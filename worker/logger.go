package worker

import (
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
	"time"
)

// Logger provides write access to a worker's logs.
type Logger interface {
	// System logs
	Debug(msg string, fields map[string]string) error
	Info(msg string, fields map[string]string) error
	Error(msg string, fields map[string]string) error

	// Task logs
	State(tes.State) error
	StartTime(t time.Time) error
	EndTime(t time.Time) error
	Outputs(o []*tes.OutputFileLog) error
	Metadata(m map[string]string) error

	// Executor logs
	ExecutorStartTime(i int, t time.Time) error
	ExecutorEndTime(i int, t time.Time) error
	ExitCode(i int, code int) error
	Ports(i int, ports []*tes.Ports) error
	HostIP(i int, ip string) error

	Stdout(i int) io.Writer
	Stderr(i int) io.Writer
}

func NewEventLogger(id string, attempt uint32, w events.Writer) EventLogger {
	return EventLogger{events.NewAttemptWriter(id, attempt, w)}
}

type EventLogger struct {
	*events.AttemptWriter
}

func (e EventLogger) Debug(msg string, fields map[string]string) error {
	return e.AttemptWriter.SystemLog(msg, "debug", fields)
}
func (e EventLogger) Info(msg string, fields map[string]string) error {
	return e.AttemptWriter.SystemLog(msg, "info", fields)
}
func (e EventLogger) Error(msg string, fields map[string]string) error {
	return e.AttemptWriter.SystemLog(msg, "error", fields)
}
func (e EventLogger) ExecutorStartTime(i int, t time.Time) error {
	return e.AttemptWriter.ExecutorStartTime(uint32(i), t)
}
func (e EventLogger) ExecutorEndTime(i int, t time.Time) error {
	return e.AttemptWriter.ExecutorEndTime(uint32(i), t)
}
func (e EventLogger) ExitCode(i int, code int) error {
	return e.AttemptWriter.ExitCode(uint32(i), int32(code))
}
func (e EventLogger) Ports(i int, ports []*tes.Ports) error {
	return e.AttemptWriter.Ports(uint32(i), ports)
}
func (e EventLogger) HostIP(i int, ip string) error {
	return e.AttemptWriter.HostIP(uint32(i), ip)
}
func (e EventLogger) Stdout(i int) io.Writer {
	return &events.StdoutWriter{e.AttemptWriter, uint32(i)}
}
func (e EventLogger) Stderr(i int) io.Writer {
	return &events.StderrWriter{e.AttemptWriter, uint32(i)}
}
