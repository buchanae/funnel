package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

type ExecError struct {
  error
}

type mustError struct {
  orig error
}
func (m mustError) Error() string {
  return m.orig.Error()
}

func Must(err error) {
	if err != nil {
		panic(mustError{err})
	}
}

func StartTask(log Logger) func(error) {
	log.StartTime(time.Now())
  log.State(tes.State_INITIALIZING)

	return func(err error) {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("Unknown worker panic: %+v", r)
			}
		}

    /*
    // Unwrap a panic error raised by Must
    if m, ok := err.(mustError); ok {
      fmt.Println("MUST UNWRAP", m.orig)
      err = m.orig
    }
    */

		log.EndTime(time.Now())
		LogFinalState(log, err)
	}
}

func LogFinalState(log Logger, err error) {
	if x, ok := err.(ExecError); ok {
		// One of the executors failed
    fmt.Println("mathc", x)
		log.Error("Exec error", map[string]string{
      "error": x.Error(),
    })
		log.State(tes.State_ERROR)

  } else if err == context.Canceled {
    // context.Canceled is a special case, because it can happen from multiple sources:
    //   - if the task is canceled by the user
    //   - if the worker is shutdown by the host (e.g. SIGKILL)
		log.State(tes.State_CANCELED)

	} else if err != nil {
		// If something else failed (system error)
		log.Error("System error", map[string]string{
      "error": err.Error(),
    })
		log.State(tes.State_SYSTEM_ERROR)

		// Otherwise, success
	} else {
		log.State(tes.State_COMPLETE)
	}
}

func StartExec(ctx context.Context, log Logger, index int) (context.Context, func()) {
	// subctx helps ensure that goroutines started while running the executor
	// are cleaned up when the executor function exits.
	subctx, cleanup := context.WithCancel(ctx)

	log.ExecutorStartTime(index, time.Now())

	return subctx, func() {
		log.ExecutorEndTime(index, time.Now())
		cleanup()
	}
}
