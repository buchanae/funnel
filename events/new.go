package events

import (
	"github.com/golang/protobuf/ptypes"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// NewState creates a state change event.
func NewState(taskID string, attempt uint32, s tes.State) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskState{
			TaskState: &TaskState{
				State: s,
			},
		},
	}
}

// NewStartTime creates a task start time event.
func NewStartTime(taskID string, attempt uint32, t time.Time) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskStartTime{
			TaskStartTime: &TaskStartTime{
				StartTime: convertTime(t),
			},
		},
	}
}

// NewEndTime creates a task end time event.
func NewEndTime(taskID string, attempt uint32, t time.Time) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskEndTime{
			TaskEndTime: &TaskEndTime{
				EndTime: convertTime(t),
			},
		},
	}
}

// NewOutputs creates a task output file log event.
func NewOutputs(taskID string, attempt uint32, f []*tes.OutputFileLog) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskOutputs{
			TaskOutputs: &TaskOutputs{
				Outputs: f,
			},
		},
	}
}

// NewMetadata creates a task metadata log event.
func NewMetadata(taskID string, attempt uint32, m map[string]string) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_TaskMetadata{
			TaskMetadata: &TaskMetadata{
				Metadata: m,
			},
		},
	}
}

// NewExecutorStartTime creates an executor start time event
// for the executor at the given index.
func NewExecutorStartTime(taskID string, attempt uint32, index uint32, t time.Time) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorStartTime{
			ExecutorStartTime: &ExecutorStartTime{
				StartTime: convertTime(t),
				Index:     index,
			},
		},
	}
}

// NewExecutorEndTime creates an executor end time event.
// for the executor at the given index.
func NewExecutorEndTime(taskID string, attempt uint32, index uint32, t time.Time) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExecutorEndTime{
			ExecutorEndTime: &ExecutorEndTime{
				EndTime: convertTime(t),
				Index:   index,
			},
		},
	}
}

// NewExitCode creates an executor exit code event
// for the executor at the given index.
func NewExitCode(taskID string, attempt uint32, index uint32, x int32) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_ExitCode{
			ExitCode: &ExitCode{
				ExitCode: x,
				Index:    index,
			},
		},
	}
}

// NewPorts creates an executor port metadata event
// for the executor at the given index.
func NewPorts(taskID string, attempt uint32, index uint32, ports []*tes.Ports) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_Ports{
			Ports: &Ports{
				Ports: ports,
				Index: index,
			},
		},
	}
}

// NewHostIP creates an executor host IP metadata event
// for the executor at the given index.
func NewHostIP(taskID string, attempt uint32, index uint32, ip string) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_HostIp{
			HostIp: &HostIp{
				HostIp: ip,
				Index:  index,
			},
		},
	}
}

// NewStdout creates an executor stdout chunk event
// for the executor at the given index.
func NewStdout(taskID string, attempt uint32, index uint32, s string) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_Stdout{
			Stdout: &Stdout{
				Stdout: s,
				Index:  index,
			},
		},
	}
}

// NewStderr creates an executor stderr chunk event
// for the executor at the given index.
func NewStderr(taskID string, attempt uint32, index uint32, s string) *Event {
	return &Event{
		Id:        taskID,
		Attempt:   attempt,
		Timestamp: ptypes.TimestampNow(),
		Event: &Event_Stderr{
			Stderr: &Stderr{
				Stderr: s,
				Index:  index,
			},
		},
	}
}

func convertTime(t time.Time) *tspb.Timestamp {
	p, _ := ptypes.TimestampProto(t)
	return p
}
