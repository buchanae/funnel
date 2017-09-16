package events

import (
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

type Writer interface {
  Write(*Event) error
}

func convertTime(t time.Time) *tspb.Timestamp {
	p, _ := ptypes.TimestampProto(t)
	return p
}

// State sets the state of the task.
func NewState(id string, attempt uint32, s tes.State) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskState{
			TaskState: &TaskState{
				State: s,
			},
		},
	}
}

// StartTime updates the task's start time log.
func NewStartTime(id string, attempt uint32, t time.Time) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskStartTime{
			TaskStartTime: &TaskStartTime{
				StartTime: convertTime(t),
			},
		},
	}
}

// EndTime updates the task's end time log.
func NewEndTime(id string, attempt uint32, t time.Time) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskEndTime{
			TaskEndTime: &TaskEndTime{
				EndTime: convertTime(t),
			},
		},
	}
}

// Outputs updates the task's output file log.
func NewOutputs(id string, attempt uint32, f []*tes.OutputFileLog) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskOutputs{
			TaskOutputs: &TaskOutputs{
				Outputs: f,
			},
		},
	}
}

// Metadata updates the task's metadata log.
func NewMetadata(id string, attempt uint32, m map[string]string) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskMetadata{
			TaskMetadata: &TaskMetadata{
				Metadata: m,
			},
		},
	}
}

// ExecutorStartTime updates an executor's start time log.
func NewExecutorStartTime(id string, attempt uint32, i uint32, t time.Time) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorStartTime{
			ExecutorStartTime: &ExecutorStartTime{
				StartTime: convertTime(t),
				Index:     i,
			},
		},
	}
}

// ExecutorEndTime updates an executor's end time log.
func NewExecutorEndTime(id string, attempt uint32, i uint32, t time.Time) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorEndTime{
			ExecutorEndTime: &ExecutorEndTime{
				EndTime: convertTime(t),
				Index:   i,
			},
		},
	}
}

// ExecutorExitCode updates an executor's exit code log.
func NewExecutorExitCode(id string, attempt uint32, i uint32, x int32) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorExitCode{
			ExecutorExitCode: &ExecutorExitCode{
				ExitCode: x,
				Index:    i,
			},
		},
	}
}

// ExecutorPorts updates an executor's ports log.
func NewExecutorPorts(id string, attempt uint32, i uint32, ports []*tes.Ports) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorPorts{
			ExecutorPorts: &ExecutorPorts{
				Ports:   ports,
				Index:   i,
			},
		},
	}
}

// ExecutorHostIP updates an executor's host IP log.
func NewExecutorHostIP(id string, attempt uint32, i uint32, ip string) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorHostIp{
			ExecutorHostIp: &ExecutorHostIp{
				HostIp:  ip,
				Index:   i,
			},
		},
	}
}

// ExecutorStdout appends to an executor's stdout log.
func NewExecutorStdout(id string, attempt uint32, i uint32, s string) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorStdout{
			ExecutorStdout: &ExecutorStdout{
				Stdout:  s,
				Index:   i,
			},
		},
	}
}

// ExecutorStderr appends to an executor's stderr log.
func NewExecutorStderr(id string, attempt uint32, i uint32, s string) *Event {
	return &Event{
		Id:        id,
		Attempt: attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorStderr{
			ExecutorStderr: &ExecutorStderr{
				Stderr:  s,
				Index:   i,
			},
		},
	}
}

// AttemptGenerator is a type that emulates the TaskWriter interface
// and outputs Events.
type AttemptGenerator struct {
	TaskID  string
	Attempt uint32
}

// NewAttemptGenerator returns a new Generator instance.
func NewAttemptGenerator(taskID string, Attempt uint32) *AttemptGenerator {
	return &AttemptGenerator{taskID, Attempt}
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
func (g *AttemptGenerator) ExecutorExitCode(i uint32, x int32) *Event {
  return NewExecutorExitCode(g.TaskID, g.Attempt, i, x)
}
func (g *AttemptGenerator) ExecutorPorts(i uint32, ports []*tes.Ports) *Event {
  return NewExecutorPorts(g.TaskID, g.Attempt, i, ports)
}
func (g *AttemptGenerator) ExecutorHostIP(i uint32, ip string) *Event {
  return NewExecutorHostIP(g.TaskID, g.Attempt, i, ip)
}
func (g *AttemptGenerator) ExecutorStdout(i uint32, s string) *Event {
  return NewExecutorStdout(g.TaskID, g.Attempt, i, s)
}
func (g *AttemptGenerator) ExecutorStderr(i uint32, s string) *Event {
  return NewExecutorStderr(g.TaskID, g.Attempt, i, s)
}
