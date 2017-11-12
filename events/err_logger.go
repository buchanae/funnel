package events

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/logger"
)

// ErrLogger writes an error message to the given logger when an event write fails.
type ErrLogger struct {
	Writer
	Log *logger.Logger
}

func (e *ErrLogger) WriteEvent(ctx context.Context, ev *Event) error {
	err := e.Writer.WriteEvent(ctx, ev)
	if err != nil {
		e.Log.Error("error writing event", err)
	}
	return err
}
