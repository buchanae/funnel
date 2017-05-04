package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/util"
)

type TaskRunner struct {
  Storage
  TaskLogger
  TaskReader
  PollRate time.Duration
}

func (r *TaskRunner) RunTask(ctx context.Context, task *tes.Task) {
  l := util.CallList{}
  task, err := r.Task()
  ctx = r.PollForCancel(ctx)

  l.AddUnchecked(func() {
    r.TaskLogger.StartTime(util.Now())
  })

  // Validate the input and outputs
  for _, p := range append(task.Inputs, task.Outputs...) {
    l.Add(func() error {
      return r.Storage.Supports(p.Url, p.Path, p.Type) {
    })
  })

	// Download inputs
	for _, input := range task.Inputs {
		l.Add(func() error {
			return r.Storage.Get(ctx, input.Url, input.Path, input.Type)
		})
	}

  // Set task to running state
  l.AddUnchecked(func() {
    r.TaskLogger.Running()
  })

	// Run executors
	for i, d := range task.Executors {
		l.Add(func() error {

      // subctx ensures goroutines are cleaned up when the step exits.
      subctx, cleanup := context.WithCancel(ctx)
      defer cleanup()

      exec := r.Executor(i, d)
      defer exec.Close()

      exec.Logger.Info("Running")
      r.TaskLogger.ExecutorStartTime(i, util.Now())
      exec.Stdout(r.TaskLogger.ExecutorStdout(i))

      // Run the executor
      done := make(chan error)
      go func() {
        done <- exec.Run(subctx, d)
      }()

      // Inspect the executor for metadata
      go func() {
        meta := exec.Inspect(subctx, d)
        r.TaskLogger.ExecutorPorts(i, meta.Ports)
        r.TaskLogger.ExecutorHostIP(i, meta.HostIP)
      }()

      // Wait for executor to exit
      res := <-done
      r.TaskLogger.ExecutorEndTime(i, util.Now())
      r.TaskLogger.ExecutorExitCode(i, getExitCode(res))
      return res
		})
	}

	// Upload outputs
	for _, output := range task.Outputs {
		l.Add(func() error {
      filelist, err := r.Storage.Put(ctx, output.Url, output.Path, output.Type)
      r.TaskLogger.OutputFiles(f)
      return err
		})
	}

  l.AddUnchecked(func() {
    r.TaskLogger.EndTime(util.Now())
  })

  result := l.Run(ctx)
  r.TaskLogger.Result(result)
}




func (r *TaskRunner) PollForCancel(ctx context.Context) context.Context {
  taskctx, cancel := context.WithCancel(ctx)

  // Start a goroutine that polls the server to watch for a canceled state.
  // If a cancel state is found, "taskctx" is canceled.
  go func() {
    ticker := time.NewTicker(r.PollRate)
    defer ticker.Stop()

    for {
    case <-taskctx.Done():
      return
    case <-ticker.C:
      state, err := reader.State()
      if state == tes.State_CANCELED {
        cancel()
      }
    }
  }()
  return taskctx
}
