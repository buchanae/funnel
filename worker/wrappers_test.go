package worker

import (
	"context"
	"errors"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"testing"
)

func TestStartTask(t *testing.T) {
	col := events.Collector{}
	log := NewEventLogger("task-id", 0, &col)

	func() {
		finish := StartTask(log)
		defer finish(nil)
		Must(nil)
	}()

	if col[0].Type != events.Type_START_TIME {
		t.Error("expected start time")
	}
	if col[1].State != tes.State_INITIALIZING {
		t.Error("expected initializing")
	}
	if col[2].Type != events.Type_END_TIME {
		t.Error("expected end time")
	}
	if col[3].State != tes.State_COMPLETE {
		t.Error("expected complete")
	}

	func() {
		finish := StartTask(log)
		defer finish(nil)
		Must(errors.New("sys err"))
	}()

	if col[7].State != tes.State_SYSTEM_ERROR {
		t.Error("expected sys err")
	}

	func() {
		finish := StartTask(log)
		defer finish(nil)
		Must(ExecError{errors.New("exec err")})
	}()

	if col[11].State != tes.State_ERROR {
		t.Error("expected error")
	}

	func() {
		finish := StartTask(log)
		defer finish(nil)
	}()

	if col[15].State != tes.State_COMPLETE {
		t.Error("expected complete")
	}

	func() {
		finish := StartTask(log)
		defer finish(errors.New("sys err"))
	}()

	if col[19].State != tes.State_SYSTEM_ERROR {
		t.Error("expected sys err")
	}

	func() {
		finish := StartTask(log)
		defer finish(ExecError{errors.New("exec err")})
	}()

	if col[23].State != tes.State_ERROR {
		t.Error("expected exec err")
	}

	func() {
		finish := StartTask(log)
		defer finish(context.Canceled)
	}()

	if col[27].State != tes.State_CANCELED {
		t.Error("expected canceled")
	}
}

func TestLogFinalState(t *testing.T) {
	col := events.Collector{}
	log := NewEventLogger("task-id", 0, &col)

	LogFinalState(log, nil)
	LogFinalState(log, errors.New("sys err foo"))
	LogFinalState(log, ExecError{errors.New("exec err foo")})

	if col[0].State != tes.State_COMPLETE {
		t.Error("expected complete but got", col[0].State)
	}
	if col[1].State != tes.State_SYSTEM_ERROR {
		t.Error("expected sys err but got", col[1].State)
	}
	if col[2].State != tes.State_ERROR {
		t.Error("expected exec err but got", col[2].State)
	}
	// TODO need system log events
}

func TestStartExec(t *testing.T) {
	col := events.Collector{}
	log := NewEventLogger("task-id", 0, &col)

	func() {
		_, finish := StartExec(context.Background(), log, 1)
		defer finish()
	}()

	if col[0].Type != events.Type_EXECUTOR_START_TIME {
		t.Error("expected exec start time")
	}
	if col[1].Type != events.Type_EXECUTOR_END_TIME {
		t.Error("expected exec start time")
	}
	if col[0].Index != 1 {
		t.Error("expected exec index 1")
	}
	if col[1].Index != 1 {
		t.Error("expected exec index 1")
	}
}
