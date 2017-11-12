package events

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// TaskEvents is a type that generates Events for a given Task execution
// attempt.
type TaskEvents struct {
	taskID  string
	attempt uint32
	sys     *SystemLogGenerator
}

// NewTaskEvents creates a TaskGenerator instance.
func NewTaskEvents(taskID string, attempt uint32) *TaskGenerator {
	return &TaskEvents{taskID, attempt, &SystemLogGenerator{taskID, attempt, 0}}
}

// State sets the state of the task.
func (eg *TaskEvents) State(s tes.State) *Event {
	return NewState(eg.taskID, eg.attempt, s)
}

// StartTime updates the task's start time log.
func (eg *TaskEvents) StartTime(t time.Time) *Event {
	return NewStartTime(eg.taskID, eg.attempt, t)
}

// EndTime updates the task's end time log.
func (eg *TaskEvents) EndTime(t time.Time) *Event {
	return NewEndTime(eg.taskID, eg.attempt, t)
}

// Outputs updates the task's output file log.
func (eg *TaskEvents) Outputs(f []*tes.OutputFileLog) *Event {
	return NewOutputs(eg.taskID, eg.attempt, f)
}

// Metadata updates the task's metadata log.
func (eg *TaskEvents) Metadata(m map[string]string) *Event {
	return NewMetadata(eg.taskID, eg.attempt, m)
}

// Info creates an info level system log message.
func (eg *TaskEvents) Info(msg string, args ...interface{}) *Event {
	return eg.sys.Info(msg, args...)
}

// Debug creates a debug level system log message.
func (eg *TaskEvents) Debug(msg string, args ...interface{}) *Event {
	return eg.sys.Debug(msg, args...)
}

// Error creates an error level system log message.
func (eg *TaskEvents) Error(msg string, args ...interface{}) *Event {
	return eg.sys.Error(msg, args...)
}
