package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"github.com/ohsu-comp-bio/funnel/storage"
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
	return &DockerWorker{c, taskID, svc.Reader, svc.Logger}, nil
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
	read   TaskReader
	log    Logger
}

// Run runs the Worker.
func (r *DockerWorker) Run(ctx context.Context) {
	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, r.read.State, r.conf.UpdateRate)

	log := logger.Sub("worker", "taskID", r.taskID)

	Start(r.log)
	defer End(r.log, nil)

	task, err := r.read.Task()
	Must(err)

	// Prepare file mapper, which maps task file URLs to host filesystem paths.
	baseDir := path.Join(r.conf.WorkDir, r.taskID)
	mapper, err := NewFileMapper(baseDir, task)
	Must(err)

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(r.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
	Must(ValidateStorageURLs(mapper.Inputs, mapper.Outputs, store))

	// Download the inputs.
	Must(Download(ctx, mapper.Inputs, store))

	// Set to running.
	r.log.State(tes.State_RUNNING)

	// TODO re-implement log tailers

	// Run task executors
	for i, exec := range task.Executors {
		// Wrap the executor to handle start/end time, context, etc.
		Must(RunExec(ctx, r.log, i, func(ctx context.Context) error {

			// Open stdin/out/err files, mapped to working directory on host
			stdio, err := mapper.NewStdio(exec.Stdin, exec.Stdout, exec.Stderr)
			Must(err)
			// Write stdout/err to the task logger event stream.
			stdio = ExecutorStdioEvents(stdio, i, r.log)

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
	r.log.Outputs(outputs)
}
