package worker

import (
	"github.com/ohsu-comp-bio/funnel/util"
	"time"
)

// Service provides common worker service behavior:
// - polling/ticker loop
// - idle timeout
// - set of task runners
//
// Service does not implement a full worker service,
// it is intended to be used as a base for concrete
// worker service implementations.
type Service struct {
  TickRate  time.Duration
  Timeout   time.Duration
	Runners   runSet
  Tick      <-chan time.Time
  stop   chan bool
}

func (w *Service) Run() {
	// Don't start if already running
	if w.stop != nil {
		return
	}
  w.stop = make(chan bool)

  timeout := util.NewIdleTimeout(w.Timeout)
  defer timeout.Stop()

  ticker := time.NewTicker(w.TickRate)
  defer ticker.Stop()

  tick := make(chan time.Time)
  defer close(tick)
  w.Tick = tick

	for {
		select {
    case t := <-w.ticker.C:
      tick <- t

      // Check if the worker is idle. If so, start the timeout timer.
      if w.Runners.Count() == 0 {
        timeout.Start()
      } else {
        timeout.Stop()
      }

		case <-timeout.Done():
			// Service timeout reached.
			w.Stop()
    case <-w.stop:
      return
		}
	}
}

// Stop stops the service
func (w *Service) Stop() {
	// Don't stop if not running
  if w.stop == nil {
		return
	}
  // Wait for runners to stop.
  w.Runners.Stop()
  // Signal the for-loop in Run() to return.
  close(w.stop)
  w.stop = nil
}
