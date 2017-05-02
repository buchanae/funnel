package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/util"
)

  // Create a backend of the given name.
  // Backends are created per-task.
  //backend, err := loadBackend(conf, taskID)
  // defer backend.Close()

func RunTask(ctx context.Context, backend Backend) {
  l := util.CallList{}
  task := backend.Task()
  ctx = PollForCancel(ctx, conf.PollRate, backend)

  l.AddUnchecked(func() {
    backend.TaskLogger.StartTime(util.Now())
  })

  // Validate the input and outputs
  for _, p := range append(task.Inputs, task.Outputs...) {
    l.Add(func() error {
      return backend.Storage.Supports(p.Url, p.Path, p.Type) {
    })
  })

	// Download inputs
	for _, input := range task.Inputs {
		l.Add(func() error {
			return backend.Storage.Get(ctx, input.Url, input.Path, input.Type)
		})
	}

  // Set task to running state
  l.AddUnchecked(func() {
    backend.TaskLogger.Running()
  })

	// Run executors
	for i, d := range task.Executors {
		l.Add(func() error {

      // subctx ensures goroutines are cleaned up when the step exits.
      subctx, cleanup := context.WithCancel(ctx)
      defer cleanup()

      exec := backend.Executor(i, d)
      defer exec.Close()

      exec.Logger.Info("Running")
      backend.TaskLogger.ExecutorStartTime(i, util.Now())
      exec.Stdout(backend.TaskLogger.ExecutorStdout(i))

      // Run the executor
      done := make(chan error)
      go func() {
        done <- exec.Run(subctx, d)
      }()

      // Inspect the executor for metadata
      go func() {
        meta := exec.Inspect(subctx, d)
        backend.TaskLogger.ExecutorPorts(i, meta.Ports)
        backend.TaskLogger.ExecutorHostIP(i, meta.HostIP)
      }()

      // Wait for executor to exit
      res := <-done
      backend.TaskLogger.ExecutorEndTime(i, util.Now())
      backend.TaskLogger.ExecutorExitCode(i, getExitCode(res))
      return res
		})
	}

	// Upload outputs
	for _, output := range task.Outputs {
		l.Add(func() error {
      filelist, err := backend.Storage.Put(ctx, output.Url, output.Path, output.Type)
      backend.TaskLogger.OutputFiles(f)
      return err
		})
	}

  l.AddUnchecked(func() {
    backend.TaskLogger.EndTime(util.Now())
  })

  result := l.Run(ctx)
  backend.TaskLogger.Result(result)
}
