package events

import (
	"context"
)

// Writer provides write access to a task's events
type Writer interface {
	WriteEvent(context.Context, *Event) error
}

type MultiWriter []Writer

// MultiWriter writes events to all the given writers.
func NewMultiWriter(ws ...Writer) MultiWriter {
	return MultiWriter(ws)
}

func (mw *MultiWriter) Add(ws ...Writer) {
	*mw = append(*mw, ws...)
}

// Write writes an event to all the writers.
func (mw *MultiWriter) WriteEvent(ctx context.Context, ev *Event) error {
	for _, w := range *mw {
		err := w.WriteEvent(ctx, ev)
		if err != nil {
			return err
		}
	}
	return nil
}
