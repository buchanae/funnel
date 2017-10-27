package events

import (
	"time"
)

// Retrier will call Write on the underlying Writer. If an error
// is returned, Retrier will try again after a short delay.
// If MaxAttempts is not set, it defaults to 3.
// If Delay is not set, it defaults to 1 second.
type Retrier struct {
	Writer
	MaxAttempts int
	Delay       time.Duration
}

func (r Retrier) Write(ev *Event) error {
	max := 3
	if r.MaxAttempts != 0 {
		max = r.MaxAttempts
	}
	delay := time.Second
	if r.Delay != 0 {
		delay = r.Delay
	}
	var err error
	for i := 0; i < max; i++ {
		err = r.Writer.Write(ev)
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}
