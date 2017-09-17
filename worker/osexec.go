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
	read   TaskReader
	log    Logger
}

func (o *OSExecWorker) Run(ctx context.Context) {
	// Poll the task service, looking for a cancel state.
	// If found, cancel the context.
	ctx = PollForCancel(ctx, o.read.State, o.conf.UpdateRate)

	// Handle start/end time, final state, panics, etc.
	finish := StartTask(o.log)
	defer finish(nil)

	task, err := o.read.Task()
	Must(err)

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
	store, err := storage.WithConfig(o.conf.Storage)
	Must(err)

	// Validate that the storage supports the input/output URLs.
  Must(ValidateStorage(store, task.Inputs, task.Outputs))

	// Download the inputs.
	Must(Download(ctx, task.Inputs, store))

	// Set to running.
	o.log.State(tes.State_RUNNING)

	// Run task executors
	for i, exec := range task.Executors {
		Must(o.runExec(ctx, i, exec))
	}

	// Upload outputs and log the outputs.
	Must(LogUpload(ctx, task.Outputs, store, o.log))
}

func (o *OSExecWorker) runExec(ctx context.Context, i int, exec *tes.Executor) error {
	ctx, finish := StartExec(ctx, o.log, i)
	defer finish()

	// Open stdin/out/err files
	stdio, err := OpenStdio(exec.Stdin, exec.Stdout, exec.Stderr)
	defer stdio.Close()
	Must(err)
	stdio = LogStdio(stdio, i, o.log)

	// Build os/exec.Cmd
	cmd := osexec.CommandContext(ctx, exec.Cmd[0], exec.Cmd[1:]...)
	cmd.Env = formatEnv(exec.Environ)
	cmd.Dir = exec.Workdir
	cmd.Stdin = stdio.In
	cmd.Stdout = stdio.Out
	cmd.Stderr = stdio.Err

	// Run and log exit code
	result := cmd.Run()
	code := GetExitCode(result)
	o.log.ExitCode(i, code)

	// TODO does not yet log ports or IP

	if result != nil {
		return ExecError{result}
	}
	return nil
}

func formatEnv(in map[string]string) []string {
	var out []string

	for k, v := range in {
		out = append(out, k+"="+v)
	}
	return out
}
