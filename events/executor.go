package events

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/util/ring"
	"golang.org/x/time/rate"
	"io"
	"time"
)

// ExecutorEvents is a type that generates Events for an Executor
// of a Task
type ExecutorEvents struct {
	taskID  string
	attempt uint32
	index   uint32
	sys     *SystemLogEvents
}

// NewExecutorEvents returns a ExecutorGenerator instance.
func NewExecutorEvents(taskID string, attempt uint32, index uint32) *ExecutorGenerator {
	return &ExecutorEvents{taskID, attempt, index, &SystemLogGenerator{taskID, attempt, index}}
}

// StartTime updates an executor's start time log.
func (eg *ExecutorEvents) StartTime(t time.Time) *Event {
	return NewExecutorStartTime(eg.taskID, eg.attempt, eg.index, t)
}

// EndTime updates an executor's end time log.
func (eg *ExecutorEvents) EndTime(t time.Time) *Event {
	return NewExecutorEndTime(eg.taskID, eg.attempt, eg.index, t)
}

// ExitCode updates an executor's exit code log.
func (eg *ExecutorEvents) ExitCode(x int) *Event {
	return NewExitCode(eg.taskID, eg.attempt, eg.index, int32(x))
}

// Stdout appends to an executor's stdout log.
func (eg *ExecutorEvents) Stdout(s string) *Event {
	return NewStdout(eg.taskID, eg.attempt, eg.index, s)
}

// Stderr appends to an executor's stderr log.
func (eg *ExecutorEvents) Stderr(s string) *Event {
	return NewStderr(eg.taskID, eg.attempt, eg.index, s)
}

// Info creates an info level system log message.
func (eg *ExecutorEvents) Info(msg string, args ...interface{}) *Event {
	return eg.sys.Info(msg, args...)
}

// Debug creates a debug level system log message.
func (eg *ExecutorEvents) Debug(msg string, args ...interface{}) *Event {
	return eg.sys.Debug(msg, args...)
}

// Error creates an error level system log message.
func (eg *ExecutorEvents) Error(msg string, args ...interface{}) *Event {
	return eg.sys.Error(msg, args...)
}

// TailLogs returns stdout/err io.Writers which will track the
// tail of the content (up to "size") and emit events. Events
// are rate limited by "interval", e.g. a max of one event every
// 5 seconds.
func (ew *ExecutorWriter) TailLogs(ctx context.Context, size int64, interval time.Duration) (stdout, stderr io.Writer) {
	return TailLogs(ctx, ew.gen.taskID, ew.gen.attempt, ew.gen.index, size, interval, ew.out)
}

// TailLogs returns stdout/err io.Writers which will track the
// tail of the content (up to "size") and emit events. Events
// are rate limited by "interval", e.g. a max of one event every
// 5 seconds.
func TailLogs(ctx context.Context, taskID string, attempt, index uint32, size int64, interval time.Duration, out Writer) (stdout, stderr io.Writer) {

	// The rate limiter allows the input writers to trigger events
	// immediately, without waiting for the ticker, as long as
	// they are not exceeding the rate limit.
	limiter := rate.NewLimiter(rate.Every(interval), 1)

	stdoutbuf := ring.NewBuffer(size)
	stderrbuf := ring.NewBuffer(size)
	stdoutch := make(chan []byte)
	stderrch := make(chan []byte)
	eventch := make(chan *Event)
	// Used as an immediate timeout for flush()
	immediate := make(chan time.Time)
	close(immediate)

	flush := func(buf *ring.Buffer, t Type, timeout <-chan time.Time) {
		// Only flush if new bytes have been written to the buffer.
		if buf.TotalWritten() == 0 {
			return
		}

		// Create the event
		var e *Event
		s := buf.String()
		switch t {
		case Type_EXECUTOR_STDOUT:
			e = NewStdout(taskID, attempt, index, s)
		case Type_EXECUTOR_STDERR:
			e = NewStderr(taskID, attempt, index, s)
		}

		// Send the event to the routine which is writing out events.
		// If it's busy, don't wait because it will block the stdout/err streams
		// writing into the logs. The logs will be flushed again soon anyway.
		select {
		case eventch <- e:
			// The writer routine accepted the event, so reset the buffer byte count.
			buf.ResetTotalWritten()
		case <-timeout:
			// The writer was busy, do nothing.
		}
	}

	flushboth := func(timeout <-chan time.Time) {
		flush(stdoutbuf, Type_EXECUTOR_STDOUT, timeout)
		flush(stderrbuf, Type_EXECUTOR_STDERR, timeout)
	}

	// There are two routines below, one for accepting input, one for writing
	// out events. They are separated so that writing out events does not block
	// the input writes. If input writes are faster than output event writes,
	// flush() calls will be dropped. This is ok, because we're flushing the
	// whole buffer (log tail) every tick, so when the output event writer
	// catches up, it will write the new, complete tail.

	// output event writer routine
	go func() {
		for e := range eventch {
			out.Write(e)
		}
	}()

	// input writes and flush routine.
	go func() {
		// The ticker helps ensure content gets flushed at a regular
		// interval, so nothing is buffered for too long.
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				timeout := time.After(time.Second * 5)
				flushboth(timeout)
				close(eventch)
				return
			case <-ticker.C:
				w := stdoutbuf.TotalWritten() + stderrbuf.TotalWritten()
				// Don't use a limiter token if not content has been written.
				if w > 0 && limiter.Allow() {
					flushboth(immediate)
				}
			case b := <-stdoutch:
				stdoutbuf.Write(b)
				if limiter.Allow() {
					flushboth(immediate)
				}
			case b := <-stderrch:
				stderrbuf.Write(b)
				if limiter.Allow() {
					flushboth(immediate)
				}
			}
		}
	}()

	return &logTailWriter{stdoutch}, &logTailWriter{stderrch}
}

type logTailWriter struct {
	ch chan<- []byte
}

func (l *logTailWriter) Write(p []byte) (n int, err error) {
	l.ch <- p
	return len(p), nil
}
