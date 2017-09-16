package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"github.com/ohsu-comp-bio/funnel/storage"
	"io"
	"path"
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
	return &DockerWorker{c, taskID, svc}, nil
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
	conf   DockerConfig
	taskID string
	svc    TaskService
}

// Run runs the Worker.
func (r *DockerWorker) Run(ctx context.Context) {

	log := logger.Sub("worker", "taskID", r.taskID)

	Start(r.svc)
	defer End(r.svc, log, nil)

	task, err := r.svc.Task()
	Must(err)

	// Map files into this baseDir
	baseDir := path.Join(r.conf.WorkDir, r.taskID)
	mapper := NewFileMapper(baseDir)

	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, r.svc.State, r.conf.UpdateRate)

	// Prepare file mapper, which maps task file URLs to host filesystem paths.
	Must(mapper.MapTask(task))

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(r.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
	Must(ValidateStorageURLs(task.Inputs, task.Outputs, store))

	// Download the inputs.
	Must(Download(ctx, task.Inputs, store))

	// Set to running.
	r.svc.SetState(tes.State_RUNNING)

	// Run task executors
	for i, exec := range task.Executors {
		// Wrap the executor to handle start/end time, context, etc.
		Must(RunExec(ctx, r.svc, i, func(ctx context.Context) error {

			log := log.WithFields("executor_index", i)
			log.Debug("Running executor")

			// Open stdin/out/err files, mapped to working directory on host
			stdio, err := mapper.OpenStdio(exec.Stdin, exec.Stdout, exec.Stderr)
			Must(err)

			// Write stdout/err to both the files and the task logger.
			stdio.Out = io.MultiWriter(stdio.Out, r.svc.ExecutorStdout(i))
			stdio.Err = io.MultiWriter(stdio.Err, r.svc.ExecutorStderr(i))

			cmd := DockerCmd{
				TaskLogger:     r.svc,
				Stdio:          stdio,
				ExecIndex:      i,
				Exec:           exec,
				Volumes:        mapper.Volumes,
				LeaveContainer: r.conf.LeaveContainer,
				ContainerName:  fmt.Sprintf("%s-%d", task.Id, i),
			}

			// docker run --rm --name [name] -i -w [workdir] -v [bindings] [imageName] [cmd]
			log.Info("Running command", "cmd", "docker "+strings.Join(cmd.Args(), " "))

			return cmd.Run(ctx)
		}))
	}

	// Fix symlinks broken by container filesystem
	for _, o := range mapper.Outputs {
		FixLinks(o.Path, mapper.HostPath)
	}

	// Upload outputs
	outputs, err := Upload(ctx, mapper.Outputs, store)
	Must(err)

	// Log task outputs.
	r.svc.Outputs(outputs)
}
