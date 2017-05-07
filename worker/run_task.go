package worker

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
  "time"
)

type DefaultRunner struct {
  TaskLogger
  TaskReader
  Storage storage.Storage
  PollRate time.Duration
}

func (r *DefaultRunner) Run(ctx context.Context) {

  var err error
  var task *tes.Task

  // Watch for the task to be canceled.
  ctx = r.PollForCancel(ctx)

  // If the context is canceled, set "err" so that steps below
  // will be skipped (based on "err").
  go func() {
    <-ctx.Done()
    err = ctx.Err()
  }()

  // Set the task result based on "err".
  // This deferred func should be defined first,
  // so that it will be run last.
  defer func() {
    r.TaskLogger.Result(err)
  }()

  // Recover from panics
  defer handlePanic(func(e error) {
    err = e
  })

  // Get the task
  task, err = r.Task()

  // Start and end time
  r.TaskLogger.StartTime(util.Now())
  defer r.TaskLogger.EndTime(util.Now())

  // Validate the input and outputs
  for _, p := range append(task.Inputs, task.Outputs...) {
    if err == nil {
      err = r.Storage.Supports(p.Url, p.Path, p.Type)
    }
  }

	// Download inputs
	for _, input := range task.Inputs {
    if err == nil {
      err = r.Storage.Get(ctx, input.Url, input.Path, input.Type)
		}
	}

  // Set task to running state
  if err == nil {
    r.TaskLogger.Running()
  }

	// Run executors
	for i, d := range task.Executors {
    if err == nil {
      // subctx ensures goroutines are cleaned up when the step exits.
      subctx, cleanup := context.WithCancel(ctx)
      defer cleanup()

      r.TaskLogger.ExecutorStartTime(i, util.Now())

      err = r.Execute(subctx, i)

      r.TaskLogger.ExecutorExitCode(i, getExitCode(err))
      r.TaskLogger.ExecutorEndTime(i, util.Now())
      cleanup()
		}
	}

	// Upload outputs
	for _, output := range task.Outputs {
    if err == nil {
      // var filelist []string
      err = r.Storage.Put(ctx, output.Url, output.Path, output.Type)
      // TODO r.TaskLogger.Outputs(filelist)
    }
	}
}


func (r *DefaultRunner) PollForCancel(ctx context.Context) context.Context {
  taskctx, cancel := context.WithCancel(ctx)

  // Start a goroutine that polls the server to watch for a canceled state.
  // If a cancel state is found, "taskctx" is canceled.
  go func() {
    ticker := time.NewTicker(r.PollRate)
    defer ticker.Stop()

    for {
      select {
      case <-taskctx.Done():
        return
      case <-ticker.C:
        state, err := r.TaskReader.State()
        // TODO look for any terminal state?
        if state == tes.State_CANCELED {
          cancel()
        }
      }
    }
  }()
  return taskctx
}

// recover from panic and call "cb" with an error value.
func handlePanic(cb func(error)) {
  if r := recover(); r != nil {
    if e, ok := r.(error); ok {
      cb(e)
    } else {
      cb(fmt.Errorf("Unknown task runner panic: %+v", r))
    }
  }
}
