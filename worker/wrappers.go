package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

type ErrExecFailed error

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Start(svc TaskService) {
	svc.StartTime(time.Now())
}

func SetFinalState(svc TaskService, log logger.Logger, err error) {
	if x, ok := err.(ErrExecFailed); ok {
		// One of the executors failed
		log.Error("Exec error", x)
		svc.SetState(tes.State_ERROR)

		// If something else failed (system error)
	} else if err != nil {
		log.Error("System error", err)
		svc.SetState(tes.State_SYSTEM_ERROR)

		// Otherwise, success
	} else {
		svc.SetState(tes.State_COMPLETE)
	}
}

func End(svc TaskService, log logger.Logger, err error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			err = e
		} else {
			err = fmt.Errorf("Unknown worker panic: %+v", r)
		}
	}

	svc.EndTime(time.Now())
	SetFinalState(svc, log, err)
}

func RunExec(ctx context.Context, tl TaskLogger, i int, f func(context.Context) error) error {
	tl.ExecutorStartTime(i, time.Now())
	defer tl.ExecutorEndTime(i, time.Now())

	// subctx helps ensure that goroutines started while running the executor
	// are cleaned up when the executor function exits.
	subctx, cleanup := context.WithCancel(ctx)
	defer cleanup()

	return f(subctx)
}
