package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	osexec "os/exec"
	"time"
)

type OSExecConfig struct {
	Storage    storage.Config
	UpdateRate time.Duration
}

type OSExecWorker struct {
	conf   OSExecConfig
	taskID string
	read   TaskReader
	log    Logger
}

func (o *OSExecWorker) Run(ctx context.Context) {
	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, o.read.State, o.conf.UpdateRate)

	Start(o.log)
	defer End(o.log, nil)

	task, err := o.read.Task()
	Must(err)

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(o.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
	Must(ValidateStorageURLs(task.Inputs, task.Outputs, store))

	// Download the inputs.
	Must(Download(ctx, task.Inputs, store))

	// Set to running.
	o.log.State(tes.State_RUNNING)

	// Run task executors
	for i, exec := range task.Executors {
		// Wrap the executor to handle start/end time, context, etc.
		Must(RunExec(ctx, o.log, i, func(ctx context.Context) error {

			// Open stdin/out/err files
			stdio, err := NewStdio(exec.Stdin, exec.Stdout, exec.Stderr)
			Must(err)
			// Write stdout/err to the task logger event stream.
			stdio = ExecutorStdioEvents(stdio, i, o.log)

			cmd := osexec.CommandContext(ctx, exec.Cmd[0], exec.Cmd[1:]...)
			cmd.Env = formatEnv(exec.Environ)
			cmd.Dir = exec.Workdir
			cmd.Stdin = stdio.In
			cmd.Stdout = stdio.Out
			cmd.Stderr = stdio.Err

			result := cmd.Run()
			code := GetExitCode(result)
			o.log.ExecutorExitCode(i, code)

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
	o.log.Outputs(outputs)
}

func formatEnv(in map[string]string) []string {
	var out []string

	for k, v := range in {
		out = append(out, k+"="+v)
	}
	return out
}
