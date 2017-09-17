package events

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// Writer defines the interface of an event writer.
type Writer interface {
	Write(*Event) error
}

// StdoutWriter provides an io.Writer interface for generating executor stdout events.
type StdoutWriter struct {
	*AttemptWriter
	Index uint32
}

// Write writes one Stdout event.
func (e *StdoutWriter) Write(p []byte) (int, error) {
	err := e.AttemptWriter.Stdout(e.Index, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// StderrWriter provides an io.Writer interface for executor stderr events.
type StderrWriter struct {
	*AttemptWriter
	Index uint32
}

// Write writes one Stderr event.
func (e *StderrWriter) Write(p []byte) (int, error) {
	err := e.AttemptWriter.Stderr(e.Index, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// AttemptWriter provides a helper for writing events for a specific task attempt.
type AttemptWriter struct {
	gen *AttemptGenerator
	w   Writer
}

// NewAttemptWriter returns a new AttemptWriter.
func NewAttemptWriter(id string, attempt uint32, w Writer) *AttemptWriter {
	return &AttemptWriter{
		&AttemptGenerator{id, attempt},
		w,
	}
}

func (a *AttemptWriter) State(s tes.State) error {
	return a.w.Write(a.gen.State(s))
}
func (a *AttemptWriter) StartTime(t time.Time) error {
	return a.w.Write(a.gen.StartTime(t))
}
func (a *AttemptWriter) EndTime(t time.Time) error {
	return a.w.Write(a.gen.EndTime(t))
}
func (a *AttemptWriter) Outputs(f []*tes.OutputFileLog) error {
	return a.w.Write(a.gen.Outputs(f))
}
func (a *AttemptWriter) Metadata(m map[string]string) error {
	return a.w.Write(a.gen.Metadata(m))
}
func (a *AttemptWriter) ExecutorStartTime(i uint32, t time.Time) error {
	return a.w.Write(a.gen.ExecutorStartTime(i, t))
}
func (a *AttemptWriter) ExecutorEndTime(i uint32, t time.Time) error {
	return a.w.Write(a.gen.ExecutorEndTime(i, t))
}
func (a *AttemptWriter) ExitCode(i uint32, x int32) error {
	return a.w.Write(a.gen.ExitCode(i, x))
}
func (a *AttemptWriter) Ports(i uint32, ports []*tes.Ports) error {
	return a.w.Write(a.gen.Ports(i, ports))
}
func (a *AttemptWriter) HostIP(i uint32, ip string) error {
	return a.w.Write(a.gen.HostIP(i, ip))
}
func (a *AttemptWriter) Stdout(i uint32, s string) error {
	return a.w.Write(a.gen.Stdout(i, s))
}
func (a *AttemptWriter) Stderr(i uint32, s string) error {
	return a.w.Write(a.gen.Stderr(i, s))
}
func (a *AttemptWriter) SystemLog(msg, lvl string, fields map[string]string) error {
	return a.w.Write(a.gen.SystemLog(msg, lvl, fields))
}

// Collector collects all events into a slice.
type Collector []*Event

func (c *Collector) Write(e *Event) error {
	*c = append(*c, e)
	return nil
}
