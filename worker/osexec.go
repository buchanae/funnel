package worker

import (
  "context"
  "github.com/ohsu-comp-bio/funnel/storage"
  "github.com/ohsu-comp-bio/funnel/logger"
  "github.com/ohsu-comp-bio/funnel/proto/tes"
  osexec "os/exec"
  "time"
  "io"
)

type OSExecConfig struct {
  Storage storage.Config
  UpdateRate time.Duration
}

type OSExecWorker struct {
  taskID string
  svc TaskService
  conf OSExecConfig
}

func (o *OSExecWorker) Run(ctx context.Context) {

	log := logger.Sub("worker", "taskID", o.taskID)

	Start(o.svc)
	defer End(o.svc, log, nil)

	task, err := o.svc.Task()
	Must(err)

	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, o.svc.State, o.conf.UpdateRate)

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(o.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
	Must(ValidateStorageURLs(task.Inputs, task.Outputs, store))

	// Download the inputs.
	Must(Download(ctx, task.Inputs, store))

	// Set to running.
	o.svc.SetState(tes.State_RUNNING)

	// Run task executors
	for i, exec := range task.Executors {
		// Wrap the executor to handle start/end time, context, etc.
		Must(RunExec(ctx, o.svc, i, func(ctx context.Context) error {

			// Open stdin/out/err files
			stdio, err := OpenStdio(exec.Stdin, exec.Stdout, exec.Stderr)
			Must(err)

			// Write stdout/err to both the files and the task logger.
			stdio.Out = io.MultiWriter(stdio.Out, o.svc.ExecutorStdout(i))
			stdio.Err = io.MultiWriter(stdio.Err, o.svc.ExecutorStderr(i))

      cmd := osexec.CommandContext(ctx, exec.Cmd[0], exec.Cmd[1:]...)
      cmd.Env = formatEnv(exec.Environ)
      cmd.Dir = exec.Workdir
      cmd.Stdin = stdio.In
      cmd.Stdout = stdio.Out
      cmd.Stderr = stdio.Err

      result := cmd.Run()
		  code := GetExitCode(result)
      o.svc.ExecutorExitCode(i, code)

      // TODO does not yet log ports or IP

      if result != nil {
        return ErrExecFailed(result)
      }
      return nil
		}))
	}

	// Upload outputs
	outputs, err := Upload(ctx, task.Outputs, store)
	Must(err)

	// Log task outputs.
	o.svc.Outputs(outputs)
}

func formatEnv(in map[string]string) []string {
  var out []string

  for k, v := range in {
    out = append(out, k + "=" + v)
  }
  return out
}
