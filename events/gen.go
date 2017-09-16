package events

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// Writer defines the interface of an event writer.
type Writer interface {
	Write(*Event) error
}

// StdoutWriter
type StdoutWriter struct {
	Writer
	TaskID  string
	Attempt uint32
	Index   uint32
}

func (e *StdoutWriter) Write(p []byte) (int, error) {
	err := e.Writer.Write(NewStdout(e.TaskID, e.Attempt, e.Index, string(p)))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type StderrWriter struct {
	Writer
	TaskID  string
	Attempt uint32
	Index   uint32
}

func (e *StderrWriter) Write(p []byte) (int, error) {
	err := e.Writer.Write(NewStderr(e.TaskID, e.Attempt, e.Index, string(p)))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// AttemptGenerator helps create events for a single task attempt.
type AttemptGenerator struct {
	TaskID  string
	Attempt uint32
}

func (g *AttemptGenerator) State(s tes.State) *Event {
	return NewState(g.TaskID, g.Attempt, s)
}
func (g *AttemptGenerator) StartTime(t time.Time) *Event {
	return NewStartTime(g.TaskID, g.Attempt, t)
}
func (g *AttemptGenerator) EndTime(t time.Time) *Event {
	return NewEndTime(g.TaskID, g.Attempt, t)
}
func (g *AttemptGenerator) Outputs(f []*tes.OutputFileLog) *Event {
	return NewOutputs(g.TaskID, g.Attempt, f)
}
func (g *AttemptGenerator) Metadata(m map[string]string) *Event {
	return NewMetadata(g.TaskID, g.Attempt, m)
}
func (g *AttemptGenerator) ExecutorStartTime(i uint32, t time.Time) *Event {
	return NewExecutorStartTime(g.TaskID, g.Attempt, i, t)
}
func (g *AttemptGenerator) ExecutorEndTime(i uint32, t time.Time) *Event {
	return NewExecutorEndTime(g.TaskID, g.Attempt, i, t)
}
func (g *AttemptGenerator) ExitCode(i uint32, x int32) *Event {
	return NewExitCode(g.TaskID, g.Attempt, i, x)
}
func (g *AttemptGenerator) Ports(i uint32, ports []*tes.Ports) *Event {
	return NewPorts(g.TaskID, g.Attempt, i, ports)
}
func (g *AttemptGenerator) HostIP(i uint32, ip string) *Event {
	return NewHostIP(g.TaskID, g.Attempt, i, ip)
}
func (g *AttemptGenerator) Stdout(i uint32, s string) *Event {
	return NewStdout(g.TaskID, g.Attempt, i, s)
}
func (g *AttemptGenerator) Stderr(i uint32, s string) *Event {
	return NewStderr(g.TaskID, g.Attempt, i, s)
}
