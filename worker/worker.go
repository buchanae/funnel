package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"path"
  "io"
  "strings"
)

// NewDefaultWorker returns the default task runner used by Funnel,
// which uses gRPC to read/write task details.
func NewDefaultWorker(conf config.Worker, taskID string) Worker {
	return &DefaultWorker{conf, taskID}
}

// DefaultWorker is the default task worker, which follows a basic,
// sequential process of task initialization, execution, finalization,
// and logging.
type DefaultWorker struct {
	Conf   config.Worker
  TaskID string
}

// Run runs the Worker.
func (r *DefaultWorker) Run(ctx context.Context) {

	log := logger.Sub("worker", "taskID", r.TaskID)

	svc, err := newRPCTask(r.Conf, r.TaskID)
  if err != nil {
    // TODO how to best expose this error?
    return
  }

  Start(svc)
  defer End(svc, log, nil)

	task, err := svc.Task()
  Must(err)

	// Map files into this baseDir
	baseDir := path.Join(r.Conf.WorkDir, r.TaskID)
  mapper := NewFileMapper(baseDir)

  // Poll the task service, looking for a cancel state.
  // If found, cancel the context.
	ctx = PollForCancel(ctx, svc.State, r.Conf.UpdateRate)

	// Prepare file mapper, which maps task file URLs to host filesystem paths.
  Must(mapper.MapTask(task))

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
  store, err := storage.WithConfig(r.Conf.Storage)
  Must(err)

  // Validate that the storage supports the input/output URLs.
  Must(ValidateStorageURLs(task.Inputs, task.Outputs, store))

  // Download the inputs.
  Must(Download(ctx, task.Inputs, store))

  // Set to running.
	svc.SetState(tes.State_RUNNING)

	// Run task executors
	for i, exec := range task.Executors {
    // Wrap the executor to handle start/end time, context, etc.
    Must(RunExec(ctx, svc, i, func(ctx context.Context) error {

      log := log.WithFields("executor_index", i)
      log.Debug("Running executor")

      // Open stdin/out/err files, mapped to working directory on host
      stdio, err := mapper.OpenStdio(exec.Stdin, exec.Stdout, exec.Stderr)
      Must(err)

      // Write stdout/err to both the files and the task logger.
      stdio.Out = io.MultiWriter(stdio.Out, svc.ExecutorStdout(i))
      stdio.Err = io.MultiWriter(stdio.Err, svc.ExecutorStderr(i))

      cmd := DockerCmd{
        TaskLogger:    svc,
        Stdio: stdio,
        ExecIndex: i,
        Exec: exec,
        Volumes: mapper.Volumes,
        // TODO make RemoveContainer configurable
        RemoveContainer: true,
        ContainerName: fmt.Sprintf("%s-%d", task.Id, i),
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
	svc.Outputs(outputs)
}
