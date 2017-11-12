package events

import (
	"golang.org/x/net/context"
)

// Writer provides write access to a task's events
type Writer interface {
	WriteEvent(context.Context, *Event) error
	Close() error
}

type multiwriter []Writer

// MultiWriter writes events to all the given writers.
func MultiWriter(ws ...Writer) Writer {
	return multiwriter(ws)
}

// Write writes an event to all the writers.
func (mw multiwriter) WriteEvent(ctx context.Context, ev *Event) error {
	for _, w := range mw {
		err := w.WriteEvent(ctx, ev)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mw multiwriter) Close() error {
	for _, w := range mw {
		w.Close()
	}
	return nil
}
