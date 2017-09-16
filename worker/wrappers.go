package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

type ErrExecFailed error

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Start(log Logger) {
	log.StartTime(time.Now())
}

func LogFinalState(log Logger, err error) {
	if x, ok := err.(ErrExecFailed); ok {
		// One of the executors failed
		log.Error("Exec error", x)
		log.State(tes.State_ERROR)

		// If something else failed (system error)
	} else if err != nil {
		log.Error("System error", err)
		log.State(tes.State_SYSTEM_ERROR)

		// Otherwise, success
	} else {
		log.State(tes.State_COMPLETE)
	}
}

func End(log Logger, err error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			err = e
		} else {
			err = fmt.Errorf("Unknown worker panic: %+v", r)
		}
	}

	log.EndTime(time.Now())
	LogFinalState(log, err)
}

func RunExec(ctx context.Context, log Logger, i int, f func(context.Context) error) error {
	log.ExecutorStartTime(i, time.Now())
	defer log.ExecutorEndTime(i, time.Now())

	// subctx helps ensure that goroutines started while running the executor
	// are cleaned up when the executor function exits.
	subctx, cleanup := context.WithCancel(ctx)
	defer cleanup()

	return f(subctx)
}
