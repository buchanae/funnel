package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"github.com/ohsu-comp-bio/funnel/storage"
	"strings"
	"time"
)

// NewDefaultWorker returns the default task runner used by Funnel,
// which uses gRPC to read/write task details.
func NewDockerRPCWorker(taskID string, c DockerConfig, r rpc.Config) (*DockerWorker, error) {
	svc, err := newRPCTask(r, taskID)
	if err != nil {
		return nil, err
	}
	// TODO baseDir := path.Join(r.conf.WorkDir, r.taskID)
	return &DockerWorker{c, svc.Reader, svc.Logger}, nil
}

type DockerConfig struct {
	Storage        storage.Config
	LeaveContainer bool
	UpdateRate     time.Duration
	WorkDir        string
}

// DefaultWorker is the default task worker, which follows a basic,
// sequential process of task initialization, execution, finalization,
// and logging.
type DockerWorker struct {
	conf DockerConfig
	read TaskReader
	log  Logger
}

// Run runs the Worker.
func (r *DockerWorker) Run(ctx context.Context) {
	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, r.read.State, r.conf.UpdateRate)

	// Handle start/end time, final state, panics, etc.
	finish := StartTask(r.log)
	defer finish(nil)

	task, err := r.read.Task()
	Must(err)

	// Prepare file mapper, which maps task file URLs to host filesystem paths.
	mapper, err := NewFileMapper(r.conf.WorkDir, task)
	Must(err)

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(r.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
	Must(store.SupportsParams(mapper.Inputs))
	Must(store.SupportsParams(mapper.Outputs))

	// Download the inputs.
	Must(Download(ctx, mapper.Inputs, store))

	// Set to running.
	r.log.State(tes.State_RUNNING)

	for i, exec := range mapper.Executors {
		func() {
			ctx, finish := StartExec(ctx, r.log, i)
			defer finish()

			// Open stdin/out/err files with log events.
			stdio, err := OpenStdio(exec.Stdin, exec.Stdout, exec.Stderr)
			defer stdio.Close()
			Must(err)
			stdio = LogStdio(stdio, i, r.log)

			cmd := DockerCmd{
				Logger:         r.log,
				Stdio:          stdio,
				ExecIndex:      i,
				Exec:           exec,
				Volumes:        mapper.Volumes,
				LeaveContainer: r.conf.LeaveContainer,
				ContainerName:  fmt.Sprintf("%s-%d", task.Id, i),
			}

			// docker run --rm --name [name] -i -w [workdir] -v [bindings] [imageName] [cmd]
			r.log.Info("Running command", map[string]string{
				"cmd": "docker " + strings.Join(cmd.Args(), " "),
			})

			defer cmd.Stop()
			Must(cmd.Run(ctx))
		}()
	}

	// Fix symlinks broken by container filesystem
	for _, o := range mapper.Outputs {
		FixLinks(o.Path, mapper.HostPath)
	}

	// Upload outputs and log the outputs.
	Must(LogUpload(ctx, mapper.Outputs, store, r.log))
}
