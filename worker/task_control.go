package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"sync"
)

type Action int
const (
  Noop Action = iota
  Stop
  Update
  Error
)

type TaskController interface {
	TaskState
  TaskLogger
}

type TaskState interface {
	Err() error
	State() tes.State
	Context() context.Context
	Complete() bool
  Canceled()
	Cancel()
}

type TaskLogger interface {
  logger.Logger
  ExitCode()
  Ports()
  IP()
  Running()
  Result()
  Stdout() io.Writer
  Stderr() io.Writer // TODO or pass as arg?
  OutputFile()
  Metadata()
  StartTime()
  EndTime()
}

func reconcile() {
  // TODO the server (or other TaskOutput backend) should handle most of the
  //      state reconciliation. An error should be returned in some cases
    resp, err := client.UpdateTask(runner.State)

    action := reconcileTaskState(resp.State, runner.State)
    switch action {
    case Update:
			task.State = runnerSt
    case Stop:
			ctrl.Cancel()
    case Error:
      log.Error("Unhandled case during worker reconciliation. Canceling.",
        "db state", db, "runner state", runner)
			ctrl.Cancel()
    }
}

// State variables for convenience
const (
	Unknown      = tes.State_UNKNOWN
	Queued       = tes.State_QUEUED
  // Active
	Initializing = tes.State_INITIALIZING
	Running      = tes.State_RUNNING
  // Terminal
	Canceled     = tes.State_CANCELED
	Complete     = tes.State_COMPLETE
	Error        = tes.State_ERROR
	SystemError  = tes.State_SYSTEM_ERROR
  // Unchecked
	Paused       = tes.State_PAUSED
)

func Active(s tes.State) bool {
	return s == Initializing || s == Running
}

func Terminal(s tes.State) bool {
	return s == Complete || s == Error || s == SystemError || s == Canceled
}

func ReconcileState(db tes.State, runner tes.State) Action {

  // This looks like a lot, but most of these are edge cases
  // that result in only a log message. Some edge cases stop
  // the task runner.
  //
  // The order of these cases is important.
  switch {

  // If the database state is terminal, stop the runner.
  // If the runner state is Unknown, Canceled, or terminal,
  // nothing will happen.
  case Terminal(db):
    return Stop

  // This is always an edge case, and could mean something
  // unusual is going on with the database, so let the task
  // runner continue in whatever state it's already in.
  case db == Unknown:
    log.Info("Unusual state during task reconciliation. Skipping.",
      "db state", db, "runner state", runner)
    return Noop

  // States match, do nothing.
  case db == runner:
    return Noop

  // Update the database with the current task state.
  case !Terminal(db) && (Active(Running) || Terminal(runner)):
    return Update

  default:
    return Error
  }
}

// NewTaskControl returns a new TaskControl instance
func NewTaskControl() TaskControl {
	ctx, cancel := context.WithCancel(context.Background())
	return &taskControl{ctx: ctx, cancelFunc: cancel}
}

type taskControl struct {
	running    bool
	complete   bool
	err        error
	mtx        sync.Mutex
  // TODO 
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (r *taskControl) Context() context.Context {
	return r.ctx
}

func (r *taskControl) SetResult(err error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// Don't set the result twice
	if !r.complete {
		r.complete = true
		r.err = err
	}
}

func (r *taskControl) SetRunning() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if !r.complete {
		r.running = true
	}
}

func (r *taskControl) SetComplete() {
  r.SetResult(nil)
}

func (r *taskControl) Err() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.err
}

func (r *taskControl) Cancel() {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.cancelFunc()
	r.err = r.ctx.Err()
	r.complete = true
}

func (r *taskControl) Complete() bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return r.complete
}

func (r *taskControl) State() tes.State {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	switch {
	case r.err == context.Canceled:
		return Canceled
	case r.err != nil:
		return Error
	case r.complete:
		return Complete
	case r.running:
		return Running
	default:
		return Initializing
	}
}
