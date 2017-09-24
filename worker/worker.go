package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/cmd/version"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/util"
	"os"
	"path/filepath"
	"time"
)

// NewDefaultWorker returns the default task runner used by Funnel,
// which uses gRPC to read/write task details.
func NewDefaultWorker(conf config.Worker, taskID string) (Worker, error) {

	rsvc, err := NewRPCTaskReader(conf, taskID)
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate TaskReader: %v", err)
	}

	rpcWriter, err := events.NewRPCWriter(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating EventService RPC client: %v", err)
	}

	logWriter := events.NewLogger("worker")

	return &DefaultWorker{
		Conf:       conf,
		Mapper:     NewFileMapper("/"),
		Store:      storage.Storage{},
		TaskReader: rsvc,
		Event:      events.MultiWriter(rpcWriter, logWriter),
	}, nil
}

// DefaultWorker is the default task worker, which follows a basic,
// sequential process of task initialization, execution, finalization,
// and logging.
type DefaultWorker struct {
	Conf       config.Worker
	Mapper     *FileMapper
	Store      storage.Storage
	TaskReader TaskReader
	Event      events.Writer
}

// Run runs the Worker.
// TODO document behavior of slow consumer of task log updates
func (r *DefaultWorker) Run(pctx context.Context) {

	// The code here is verbose, but simple; mainly loops and simple error checking.
	//
	// The steps are:
	// - prepare the working directory
	// - map the task files to the working directory
	// - log the IP address
	// - set up the storage configuration
	// - validate input and output files
	// - download inputs
	// - run the steps (docker)
	// - upload the outputs

	var run helper

	var task *tes.Task
	var terr error
	task, terr = r.TaskReader.Task()

	if terr != nil || task == nil {
		// TODO log error
		return
	}
	event := events.NewTaskWriter(task.Id, uint32(len(task.Logs)), r.Conf.Logger.Level, r.Event)

	event.Info("Version", version.LogFields()...)

	if run.ok() {
		event.State(tes.State_INITIALIZING)
	}

	event.StartTime(time.Now())
	// Run the final logging/state steps in a deferred function
	// to ensure they always run, even if there's a missed error.
	defer func() {
		event.EndTime(time.Now())

		switch {
    case run.taskCanceled:
			event.State(tes.State_CANCELED)
    case run.syserr == context.Canceled:
      fallthrough
    case run.execerr == context.Canceled:
      event.Error("System canceled")
      event.State(tes.State_SYSTEM_ERROR)
		case run.execerr != nil:
			// One of the executors failed
			event.Error("Exec error", run.execerr)
			event.State(tes.State_ERROR)
		case run.syserr != nil:
			// Something else failed
			event.Error("System error", run.syserr)
			event.State(tes.State_SYSTEM_ERROR)
		default:
			event.State(tes.State_COMPLETE)
		}

    for _, in := range task.Inputs {
      os.RemoveAll(in.Path)
    }
	}()

	// Recover from panics
	defer handlePanic(func(e error) {
		run.syserr = e
	})

	ctx := r.pollForCancel(pctx, func() {
		run.taskCanceled = true
	})
	run.ctx = ctx

	// Create working dir
	var dir string
	if run.ok() {
		dir, run.syserr = filepath.Abs(r.Conf.WorkDir)
	}
	if run.ok() {
		run.syserr = util.EnsureDir(dir)
	}

	// Prepare file mapper, which maps task file URLs to host filesystem paths
	if run.ok() {
		run.syserr = r.Mapper.MapTask(task)
	}

	// Grab the IP address of this host. Used to send task metadata updates.
	var ip string
	if run.ok() {
		ip, run.syserr = externalIP()
	}

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	if run.ok() {
		r.Store, run.syserr = r.Store.WithConfig(r.Conf.Storage)
	}

	if run.ok() {
		run.syserr = r.validateInputs()
	}

	if run.ok() {
		run.syserr = r.validateOutputs()
	}

	// Download inputs
	for _, input := range r.Mapper.Inputs {
		if run.ok() {
			run.syserr = r.Store.Get(ctx, input.Url, input.Path, input.Type)
		}
	}

	if run.ok() {
		event.State(tes.State_RUNNING)
	}

	// Run steps
	for i, d := range task.Executors {
		s := &stepWorker{
			Conf:  r.Conf,
			Event: event.NewExecutorWriter(uint32(i)),
			IP:    ip,
			Cmd: &DockerCmd{
				ImageName:     d.ImageName,
				Cmd:           d.Cmd,
				Environ:       d.Environ,
				Volumes:       r.Mapper.Volumes,
				Workdir:       d.Workdir,
				Ports:         d.Ports,
				ContainerName: fmt.Sprintf("%s-%d", task.Id, i),
				// TODO make RemoveContainer configurable
				RemoveContainer: true,
				Event:           event.NewExecutorWriter(uint32(i)),
			},
		}

		// Opens stdin/out/err files and updates those fields on "cmd".
		if run.ok() {
			run.syserr = r.openStepLogs(s, d)
		}

		if run.ok() {
			run.execerr = s.Run(ctx)
		}
	}

	// Upload outputs
	var outputs []*tes.OutputFileLog
	for _, output := range r.Mapper.Outputs {
		if run.ok() {
			r.fixLinks(output.Path)
			var out []*tes.OutputFileLog
			out, run.syserr = r.Store.Put(ctx, output.Url, output.Path, output.Type)
			outputs = append(outputs, out...)
		}
	}

	if run.ok() {
		event.Outputs(outputs)
	}
}

// fixLinks walks the output paths, fixing cases where a symlink is
// broken because it's pointing to a path inside a container volume.
func (r *DefaultWorker) fixLinks(basepath string) {
	filepath.Walk(basepath, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			// There's an error, so be safe and give up on this file
			return nil
		}

		// Only bother to check symlinks
		if f.Mode()&os.ModeSymlink != 0 {
			// Test if the file can be opened because it doesn't exist
			fh, rerr := os.Open(p)
			fh.Close()

			if rerr != nil && os.IsNotExist(rerr) {

				// Get symlink source path
				src, err := os.Readlink(p)
				if err != nil {
					return nil
				}
				// Map symlink source (possible container path) to host path
				mapped, err := r.Mapper.HostPath(src)
				if err != nil {
					return nil
				}

				// Check whether the mapped path exists
				fh, err := os.Open(mapped)
				fh.Close()

				// If the mapped path exists, fix the symlink
				if err == nil {
					err := os.Remove(p)
					if err != nil {
						return nil
					}
					os.Symlink(mapped, p)
				}
			}
		}
		return nil
	})
}

// openLogs opens/creates the logs files for a step and updates those fields.
func (r *DefaultWorker) openStepLogs(s *stepWorker, d *tes.Executor) error {

	// Find the path for task stdin
	var err error
	if d.Stdin != "" {
		s.Cmd.Stdin, err = r.Mapper.OpenHostFile(d.Stdin)
		if err != nil {
			s.Event.Error("Couldn't prepare log files", err)
			return err
		}
	}

	// Create file for task stdout
	if d.Stdout != "" {
		s.Cmd.Stdout, err = r.Mapper.CreateHostFile(d.Stdout)
		if err != nil {
			s.Event.Error("Couldn't prepare log files", err)
			return err
		}
	}

	// Create file for task stderr
	if d.Stderr != "" {
		s.Cmd.Stderr, err = r.Mapper.CreateHostFile(d.Stderr)
		if err != nil {
			s.Event.Error("Couldn't prepare log files", err)
			return err
		}
	}
	return nil
}

// Validate the input downloads
func (r *DefaultWorker) validateInputs() error {
	for _, input := range r.Mapper.Inputs {
		if !r.Store.Supports(input.Url, input.Path, input.Type) {
			return fmt.Errorf("Input download not supported by storage: %v", input)
		}
	}
	return nil
}

// Validate the output uploads
func (r *DefaultWorker) validateOutputs() error {
	for _, output := range r.Mapper.Outputs {
		if !r.Store.Supports(output.Url, output.Path, output.Type) {
			return fmt.Errorf("Output upload not supported by storage: %v", output)
		}
	}
	return nil
}

func (r *DefaultWorker) pollForCancel(ctx context.Context, f func()) context.Context {
	taskctx, cancel := context.WithCancel(ctx)

	// Start a goroutine that polls the server to watch for a canceled state.
	// If a cancel state is found, "taskctx" is canceled.
	go func() {
		ticker := time.NewTicker(r.Conf.UpdateRate)
		defer ticker.Stop()

		for {
			select {
			case <-taskctx.Done():
				return
			case <-ticker.C:
				state, _ := r.TaskReader.State()
				if tes.TerminalState(state) {
					cancel()
					f()
				}
			}
		}
	}()
	return taskctx
}
