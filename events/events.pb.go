// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events.proto

/*
Package events is a generated protocol buffer package.

It is generated from these files:
	events.proto

It has these top-level messages:
	TaskState
	TaskStartTime
	TaskEndTime
	TaskOutputs
	TaskMetadata
	ExecutorStartTime
	ExecutorEndTime
	ExecutorExitCode
	ExecutorHostIp
	ExecutorPorts
	ExecutorStdout
	ExecutorStderr
	SystemLog
	Event
	CreateEventResponse
*/
package events

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import tes "github.com/ohsu-comp-bio/funnel/proto/tes"
import google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
import google_protobuf2 "github.com/golang/protobuf/ptypes/timestamp"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TaskState struct {
	State tes.State `protobuf:"varint,1,opt,name=state,enum=tes.State" json:"state,omitempty"`
}

func (m *TaskState) Reset()                    { *m = TaskState{} }
func (m *TaskState) String() string            { return proto.CompactTextString(m) }
func (*TaskState) ProtoMessage()               {}
func (*TaskState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TaskState) GetState() tes.State {
	if m != nil {
		return m.State
	}
	return tes.State_UNKNOWN
}

type TaskStartTime struct {
	StartTime *google_protobuf2.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
}

func (m *TaskStartTime) Reset()                    { *m = TaskStartTime{} }
func (m *TaskStartTime) String() string            { return proto.CompactTextString(m) }
func (*TaskStartTime) ProtoMessage()               {}
func (*TaskStartTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TaskStartTime) GetStartTime() *google_protobuf2.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

type TaskEndTime struct {
	EndTime *google_protobuf2.Timestamp `protobuf:"bytes,1,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
}

func (m *TaskEndTime) Reset()                    { *m = TaskEndTime{} }
func (m *TaskEndTime) String() string            { return proto.CompactTextString(m) }
func (*TaskEndTime) ProtoMessage()               {}
func (*TaskEndTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TaskEndTime) GetEndTime() *google_protobuf2.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

type TaskOutputs struct {
	Outputs []*tes.OutputFileLog `protobuf:"bytes,1,rep,name=outputs" json:"outputs,omitempty"`
}

func (m *TaskOutputs) Reset()                    { *m = TaskOutputs{} }
func (m *TaskOutputs) String() string            { return proto.CompactTextString(m) }
func (*TaskOutputs) ProtoMessage()               {}
func (*TaskOutputs) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *TaskOutputs) GetOutputs() []*tes.OutputFileLog {
	if m != nil {
		return m.Outputs
	}
	return nil
}

type TaskMetadata struct {
	Metadata map[string]string `protobuf:"bytes,1,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *TaskMetadata) Reset()                    { *m = TaskMetadata{} }
func (m *TaskMetadata) String() string            { return proto.CompactTextString(m) }
func (*TaskMetadata) ProtoMessage()               {}
func (*TaskMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *TaskMetadata) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type ExecutorStartTime struct {
	StartTime *google_protobuf2.Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	Index     uint32                      `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorStartTime) Reset()                    { *m = ExecutorStartTime{} }
func (m *ExecutorStartTime) String() string            { return proto.CompactTextString(m) }
func (*ExecutorStartTime) ProtoMessage()               {}
func (*ExecutorStartTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ExecutorStartTime) GetStartTime() *google_protobuf2.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *ExecutorStartTime) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorEndTime struct {
	EndTime *google_protobuf2.Timestamp `protobuf:"bytes,1,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
	Index   uint32                      `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorEndTime) Reset()                    { *m = ExecutorEndTime{} }
func (m *ExecutorEndTime) String() string            { return proto.CompactTextString(m) }
func (*ExecutorEndTime) ProtoMessage()               {}
func (*ExecutorEndTime) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ExecutorEndTime) GetEndTime() *google_protobuf2.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *ExecutorEndTime) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorExitCode struct {
	ExitCode int32  `protobuf:"varint,1,opt,name=exit_code,json=exitCode" json:"exit_code,omitempty"`
	Index    uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorExitCode) Reset()                    { *m = ExecutorExitCode{} }
func (m *ExecutorExitCode) String() string            { return proto.CompactTextString(m) }
func (*ExecutorExitCode) ProtoMessage()               {}
func (*ExecutorExitCode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ExecutorExitCode) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *ExecutorExitCode) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorHostIp struct {
	HostIp string `protobuf:"bytes,1,opt,name=host_ip,json=hostIp" json:"host_ip,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorHostIp) Reset()                    { *m = ExecutorHostIp{} }
func (m *ExecutorHostIp) String() string            { return proto.CompactTextString(m) }
func (*ExecutorHostIp) ProtoMessage()               {}
func (*ExecutorHostIp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ExecutorHostIp) GetHostIp() string {
	if m != nil {
		return m.HostIp
	}
	return ""
}

func (m *ExecutorHostIp) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorPorts struct {
	Ports []*tes.Ports `protobuf:"bytes,1,rep,name=ports" json:"ports,omitempty"`
	Index uint32       `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorPorts) Reset()                    { *m = ExecutorPorts{} }
func (m *ExecutorPorts) String() string            { return proto.CompactTextString(m) }
func (*ExecutorPorts) ProtoMessage()               {}
func (*ExecutorPorts) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ExecutorPorts) GetPorts() []*tes.Ports {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *ExecutorPorts) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorStdout struct {
	Stdout string `protobuf:"bytes,1,opt,name=stdout" json:"stdout,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorStdout) Reset()                    { *m = ExecutorStdout{} }
func (m *ExecutorStdout) String() string            { return proto.CompactTextString(m) }
func (*ExecutorStdout) ProtoMessage()               {}
func (*ExecutorStdout) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ExecutorStdout) GetStdout() string {
	if m != nil {
		return m.Stdout
	}
	return ""
}

func (m *ExecutorStdout) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type ExecutorStderr struct {
	Stderr string `protobuf:"bytes,1,opt,name=stderr" json:"stderr,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExecutorStderr) Reset()                    { *m = ExecutorStderr{} }
func (m *ExecutorStderr) String() string            { return proto.CompactTextString(m) }
func (*ExecutorStderr) ProtoMessage()               {}
func (*ExecutorStderr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ExecutorStderr) GetStderr() string {
	if m != nil {
		return m.Stderr
	}
	return ""
}

func (m *ExecutorStderr) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type SystemLog struct {
	Msg    string                   `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
	Level  uint32                   `protobuf:"varint,2,opt,name=level" json:"level,omitempty"`
	Fields *google_protobuf1.Struct `protobuf:"bytes,3,opt,name=fields" json:"fields,omitempty"`
}

func (m *SystemLog) Reset()                    { *m = SystemLog{} }
func (m *SystemLog) String() string            { return proto.CompactTextString(m) }
func (*SystemLog) ProtoMessage()               {}
func (*SystemLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *SystemLog) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *SystemLog) GetLevel() uint32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *SystemLog) GetFields() *google_protobuf1.Struct {
	if m != nil {
		return m.Fields
	}
	return nil
}

type Event struct {
	Id        string                      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Attempt   uint32                      `protobuf:"varint,2,opt,name=attempt" json:"attempt,omitempty"`
	Timestamp *google_protobuf2.Timestamp `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp,omitempty"`
	// Types that are valid to be assigned to Event:
	//	*Event_TaskState
	//	*Event_TaskStartTime
	//	*Event_TaskEndTime
	//	*Event_TaskOutputs
	//	*Event_TaskMetadata
	//	*Event_ExecutorStartTime
	//	*Event_ExecutorEndTime
	//	*Event_ExecutorExitCode
	//	*Event_ExecutorHostIp
	//	*Event_ExecutorPorts
	//	*Event_ExecutorStdout
	//	*Event_ExecutorStderr
	//	*Event_SystemLog
	Event isEvent_Event `protobuf_oneof:"event"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type isEvent_Event interface {
	isEvent_Event()
}

type Event_TaskState struct {
	TaskState *TaskState `protobuf:"bytes,4,opt,name=task_state,json=taskState,oneof"`
}
type Event_TaskStartTime struct {
	TaskStartTime *TaskStartTime `protobuf:"bytes,5,opt,name=task_start_time,json=taskStartTime,oneof"`
}
type Event_TaskEndTime struct {
	TaskEndTime *TaskEndTime `protobuf:"bytes,6,opt,name=task_end_time,json=taskEndTime,oneof"`
}
type Event_TaskOutputs struct {
	TaskOutputs *TaskOutputs `protobuf:"bytes,7,opt,name=task_outputs,json=taskOutputs,oneof"`
}
type Event_TaskMetadata struct {
	TaskMetadata *TaskMetadata `protobuf:"bytes,8,opt,name=task_metadata,json=taskMetadata,oneof"`
}
type Event_ExecutorStartTime struct {
	ExecutorStartTime *ExecutorStartTime `protobuf:"bytes,9,opt,name=executor_start_time,json=executorStartTime,oneof"`
}
type Event_ExecutorEndTime struct {
	ExecutorEndTime *ExecutorEndTime `protobuf:"bytes,10,opt,name=executor_end_time,json=executorEndTime,oneof"`
}
type Event_ExecutorExitCode struct {
	ExecutorExitCode *ExecutorExitCode `protobuf:"bytes,11,opt,name=executor_exit_code,json=executorExitCode,oneof"`
}
type Event_ExecutorHostIp struct {
	ExecutorHostIp *ExecutorHostIp `protobuf:"bytes,12,opt,name=executor_host_ip,json=executorHostIp,oneof"`
}
type Event_ExecutorPorts struct {
	ExecutorPorts *ExecutorPorts `protobuf:"bytes,13,opt,name=executor_ports,json=executorPorts,oneof"`
}
type Event_ExecutorStdout struct {
	ExecutorStdout *ExecutorStdout `protobuf:"bytes,14,opt,name=executor_stdout,json=executorStdout,oneof"`
}
type Event_ExecutorStderr struct {
	ExecutorStderr *ExecutorStderr `protobuf:"bytes,15,opt,name=executor_stderr,json=executorStderr,oneof"`
}
type Event_SystemLog struct {
	SystemLog *SystemLog `protobuf:"bytes,16,opt,name=system_log,json=systemLog,oneof"`
}

func (*Event_TaskState) isEvent_Event()         {}
func (*Event_TaskStartTime) isEvent_Event()     {}
func (*Event_TaskEndTime) isEvent_Event()       {}
func (*Event_TaskOutputs) isEvent_Event()       {}
func (*Event_TaskMetadata) isEvent_Event()      {}
func (*Event_ExecutorStartTime) isEvent_Event() {}
func (*Event_ExecutorEndTime) isEvent_Event()   {}
func (*Event_ExecutorExitCode) isEvent_Event()  {}
func (*Event_ExecutorHostIp) isEvent_Event()    {}
func (*Event_ExecutorPorts) isEvent_Event()     {}
func (*Event_ExecutorStdout) isEvent_Event()    {}
func (*Event_ExecutorStderr) isEvent_Event()    {}
func (*Event_SystemLog) isEvent_Event()         {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetAttempt() uint32 {
	if m != nil {
		return m.Attempt
	}
	return 0
}

func (m *Event) GetTimestamp() *google_protobuf2.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetTaskState() *TaskState {
	if x, ok := m.GetEvent().(*Event_TaskState); ok {
		return x.TaskState
	}
	return nil
}

func (m *Event) GetTaskStartTime() *TaskStartTime {
	if x, ok := m.GetEvent().(*Event_TaskStartTime); ok {
		return x.TaskStartTime
	}
	return nil
}

func (m *Event) GetTaskEndTime() *TaskEndTime {
	if x, ok := m.GetEvent().(*Event_TaskEndTime); ok {
		return x.TaskEndTime
	}
	return nil
}

func (m *Event) GetTaskOutputs() *TaskOutputs {
	if x, ok := m.GetEvent().(*Event_TaskOutputs); ok {
		return x.TaskOutputs
	}
	return nil
}

func (m *Event) GetTaskMetadata() *TaskMetadata {
	if x, ok := m.GetEvent().(*Event_TaskMetadata); ok {
		return x.TaskMetadata
	}
	return nil
}

func (m *Event) GetExecutorStartTime() *ExecutorStartTime {
	if x, ok := m.GetEvent().(*Event_ExecutorStartTime); ok {
		return x.ExecutorStartTime
	}
	return nil
}

func (m *Event) GetExecutorEndTime() *ExecutorEndTime {
	if x, ok := m.GetEvent().(*Event_ExecutorEndTime); ok {
		return x.ExecutorEndTime
	}
	return nil
}

func (m *Event) GetExecutorExitCode() *ExecutorExitCode {
	if x, ok := m.GetEvent().(*Event_ExecutorExitCode); ok {
		return x.ExecutorExitCode
	}
	return nil
}

func (m *Event) GetExecutorHostIp() *ExecutorHostIp {
	if x, ok := m.GetEvent().(*Event_ExecutorHostIp); ok {
		return x.ExecutorHostIp
	}
	return nil
}

func (m *Event) GetExecutorPorts() *ExecutorPorts {
	if x, ok := m.GetEvent().(*Event_ExecutorPorts); ok {
		return x.ExecutorPorts
	}
	return nil
}

func (m *Event) GetExecutorStdout() *ExecutorStdout {
	if x, ok := m.GetEvent().(*Event_ExecutorStdout); ok {
		return x.ExecutorStdout
	}
	return nil
}

func (m *Event) GetExecutorStderr() *ExecutorStderr {
	if x, ok := m.GetEvent().(*Event_ExecutorStderr); ok {
		return x.ExecutorStderr
	}
	return nil
}

func (m *Event) GetSystemLog() *SystemLog {
	if x, ok := m.GetEvent().(*Event_SystemLog); ok {
		return x.SystemLog
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_TaskState)(nil),
		(*Event_TaskStartTime)(nil),
		(*Event_TaskEndTime)(nil),
		(*Event_TaskOutputs)(nil),
		(*Event_TaskMetadata)(nil),
		(*Event_ExecutorStartTime)(nil),
		(*Event_ExecutorEndTime)(nil),
		(*Event_ExecutorExitCode)(nil),
		(*Event_ExecutorHostIp)(nil),
		(*Event_ExecutorPorts)(nil),
		(*Event_ExecutorStdout)(nil),
		(*Event_ExecutorStderr)(nil),
		(*Event_SystemLog)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// event
	switch x := m.Event.(type) {
	case *Event_TaskState:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TaskState); err != nil {
			return err
		}
	case *Event_TaskStartTime:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TaskStartTime); err != nil {
			return err
		}
	case *Event_TaskEndTime:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TaskEndTime); err != nil {
			return err
		}
	case *Event_TaskOutputs:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TaskOutputs); err != nil {
			return err
		}
	case *Event_TaskMetadata:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TaskMetadata); err != nil {
			return err
		}
	case *Event_ExecutorStartTime:
		b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorStartTime); err != nil {
			return err
		}
	case *Event_ExecutorEndTime:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorEndTime); err != nil {
			return err
		}
	case *Event_ExecutorExitCode:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorExitCode); err != nil {
			return err
		}
	case *Event_ExecutorHostIp:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorHostIp); err != nil {
			return err
		}
	case *Event_ExecutorPorts:
		b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorPorts); err != nil {
			return err
		}
	case *Event_ExecutorStdout:
		b.EncodeVarint(14<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorStdout); err != nil {
			return err
		}
	case *Event_ExecutorStderr:
		b.EncodeVarint(15<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExecutorStderr); err != nil {
			return err
		}
	case *Event_SystemLog:
		b.EncodeVarint(16<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SystemLog); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 4: // event.task_state
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TaskState)
		err := b.DecodeMessage(msg)
		m.Event = &Event_TaskState{msg}
		return true, err
	case 5: // event.task_start_time
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TaskStartTime)
		err := b.DecodeMessage(msg)
		m.Event = &Event_TaskStartTime{msg}
		return true, err
	case 6: // event.task_end_time
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TaskEndTime)
		err := b.DecodeMessage(msg)
		m.Event = &Event_TaskEndTime{msg}
		return true, err
	case 7: // event.task_outputs
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TaskOutputs)
		err := b.DecodeMessage(msg)
		m.Event = &Event_TaskOutputs{msg}
		return true, err
	case 8: // event.task_metadata
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TaskMetadata)
		err := b.DecodeMessage(msg)
		m.Event = &Event_TaskMetadata{msg}
		return true, err
	case 9: // event.executor_start_time
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorStartTime)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorStartTime{msg}
		return true, err
	case 10: // event.executor_end_time
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorEndTime)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorEndTime{msg}
		return true, err
	case 11: // event.executor_exit_code
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorExitCode)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorExitCode{msg}
		return true, err
	case 12: // event.executor_host_ip
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorHostIp)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorHostIp{msg}
		return true, err
	case 13: // event.executor_ports
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorPorts)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorPorts{msg}
		return true, err
	case 14: // event.executor_stdout
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorStdout)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorStdout{msg}
		return true, err
	case 15: // event.executor_stderr
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExecutorStderr)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExecutorStderr{msg}
		return true, err
	case 16: // event.system_log
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SystemLog)
		err := b.DecodeMessage(msg)
		m.Event = &Event_SystemLog{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// event
	switch x := m.Event.(type) {
	case *Event_TaskState:
		s := proto.Size(x.TaskState)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_TaskStartTime:
		s := proto.Size(x.TaskStartTime)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_TaskEndTime:
		s := proto.Size(x.TaskEndTime)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_TaskOutputs:
		s := proto.Size(x.TaskOutputs)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_TaskMetadata:
		s := proto.Size(x.TaskMetadata)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorStartTime:
		s := proto.Size(x.ExecutorStartTime)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorEndTime:
		s := proto.Size(x.ExecutorEndTime)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorExitCode:
		s := proto.Size(x.ExecutorExitCode)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorHostIp:
		s := proto.Size(x.ExecutorHostIp)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorPorts:
		s := proto.Size(x.ExecutorPorts)
		n += proto.SizeVarint(13<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorStdout:
		s := proto.Size(x.ExecutorStdout)
		n += proto.SizeVarint(14<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ExecutorStderr:
		s := proto.Size(x.ExecutorStderr)
		n += proto.SizeVarint(15<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_SystemLog:
		s := proto.Size(x.SystemLog)
		n += proto.SizeVarint(16<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type CreateEventResponse struct {
}

func (m *CreateEventResponse) Reset()                    { *m = CreateEventResponse{} }
func (m *CreateEventResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateEventResponse) ProtoMessage()               {}
func (*CreateEventResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func init() {
	proto.RegisterType((*TaskState)(nil), "events.TaskState")
	proto.RegisterType((*TaskStartTime)(nil), "events.TaskStartTime")
	proto.RegisterType((*TaskEndTime)(nil), "events.TaskEndTime")
	proto.RegisterType((*TaskOutputs)(nil), "events.TaskOutputs")
	proto.RegisterType((*TaskMetadata)(nil), "events.TaskMetadata")
	proto.RegisterType((*ExecutorStartTime)(nil), "events.ExecutorStartTime")
	proto.RegisterType((*ExecutorEndTime)(nil), "events.ExecutorEndTime")
	proto.RegisterType((*ExecutorExitCode)(nil), "events.ExecutorExitCode")
	proto.RegisterType((*ExecutorHostIp)(nil), "events.ExecutorHostIp")
	proto.RegisterType((*ExecutorPorts)(nil), "events.ExecutorPorts")
	proto.RegisterType((*ExecutorStdout)(nil), "events.ExecutorStdout")
	proto.RegisterType((*ExecutorStderr)(nil), "events.ExecutorStderr")
	proto.RegisterType((*SystemLog)(nil), "events.SystemLog")
	proto.RegisterType((*Event)(nil), "events.Event")
	proto.RegisterType((*CreateEventResponse)(nil), "events.CreateEventResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EventService service

type EventServiceClient interface {
	CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventResponse, error)
}

type eventServiceClient struct {
	cc *grpc.ClientConn
}

func NewEventServiceClient(cc *grpc.ClientConn) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := grpc.Invoke(ctx, "/events.EventService/CreateEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EventService service

type EventServiceServer interface {
	CreateEvent(context.Context, *Event) (*CreateEventResponse, error)
}

func RegisterEventServiceServer(s *grpc.Server, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/events.EventService/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "events.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "events.proto",
}

func init() { proto.RegisterFile("events.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 835 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x5d, 0x6f, 0xe3, 0x54,
	0x10, 0x4d, 0x5a, 0x92, 0xd4, 0x93, 0x38, 0x69, 0x6f, 0xf7, 0x23, 0x64, 0x79, 0xa8, 0xfc, 0xd4,
	0x07, 0x48, 0xa5, 0x20, 0xa4, 0x42, 0xa5, 0xae, 0xd8, 0x25, 0x60, 0xd8, 0x45, 0xa0, 0x9b, 0x3e,
	0x63, 0x79, 0xeb, 0xd9, 0xac, 0xd5, 0xd8, 0xd7, 0xf2, 0x9d, 0x54, 0xed, 0x6f, 0xe0, 0x07, 0xf1,
	0xf7, 0xd0, 0xfd, 0x72, 0x1d, 0x27, 0x45, 0x95, 0xd8, 0xa7, 0xdc, 0x99, 0xb9, 0xe7, 0xcc, 0xf5,
	0x64, 0xce, 0x81, 0x01, 0xde, 0x62, 0x4e, 0x72, 0x5a, 0x94, 0x82, 0x04, 0xeb, 0x9a, 0x68, 0xe2,
	0x11, 0xda, 0xd4, 0x64, 0x20, 0xa9, 0x5c, 0x5f, 0x93, 0x8d, 0x46, 0x94, 0x66, 0x28, 0x29, 0xce,
	0x0a, 0x9b, 0xf8, 0x6a, 0x29, 0xc4, 0x72, 0x85, 0x67, 0x71, 0x91, 0x9e, 0xc5, 0x79, 0x2e, 0x28,
	0xa6, 0x54, 0xe4, 0x16, 0x1c, 0x7c, 0x03, 0xde, 0x55, 0x2c, 0x6f, 0x16, 0x14, 0x13, 0xb2, 0x13,
	0xe8, 0x48, 0x75, 0x18, 0xb7, 0x4f, 0xda, 0xa7, 0xc3, 0x19, 0x4c, 0x55, 0x13, 0x5d, 0xe2, 0xa6,
	0x10, 0xfc, 0x06, 0xbe, 0xbd, 0x5e, 0xd2, 0x55, 0x9a, 0x21, 0xfb, 0x1e, 0x40, 0xaa, 0x20, 0x52,
	0x6d, 0x35, 0xae, 0x3f, 0x9b, 0x4c, 0x4d, 0x4b, 0xd3, 0xe2, 0xc3, 0xfa, 0xe3, 0xf4, 0xca, 0xbd,
	0x89, 0x7b, 0xd2, 0x41, 0x83, 0x9f, 0xa0, 0xaf, 0xb8, 0xe6, 0x79, 0xa2, 0x99, 0xbe, 0x83, 0x03,
	0xcc, 0x93, 0xa7, 0xf2, 0xf4, 0xd0, 0xc0, 0x82, 0x0b, 0xc3, 0xf2, 0xc7, 0x9a, 0x8a, 0x35, 0x49,
	0xf6, 0x35, 0xf4, 0x84, 0x39, 0x8e, 0xdb, 0x27, 0xfb, 0xa7, 0xfd, 0x19, 0xd3, 0x1f, 0x61, 0xca,
	0x3f, 0xa7, 0x2b, 0x7c, 0x2f, 0x96, 0xdc, 0x5d, 0x09, 0xfe, 0x6e, 0xc3, 0x40, 0xa1, 0x7f, 0x47,
	0x8a, 0x93, 0x98, 0x62, 0x76, 0x09, 0x07, 0x99, 0x3d, 0x5b, 0x7c, 0x30, 0xb5, 0xf3, 0xaf, 0xdf,
	0x9b, 0xba, 0xc3, 0x3c, 0xa7, 0xf2, 0x9e, 0x57, 0x98, 0xc9, 0x05, 0xf8, 0x1b, 0x25, 0x76, 0x08,
	0xfb, 0x37, 0x78, 0xaf, 0x3f, 0xc8, 0xe3, 0xea, 0xc8, 0x9e, 0x41, 0xe7, 0x36, 0x5e, 0xad, 0x71,
	0xbc, 0xa7, 0x73, 0x26, 0xf8, 0x61, 0xef, 0xbc, 0x1d, 0x24, 0x70, 0x34, 0xbf, 0xc3, 0xeb, 0x35,
	0x89, 0xf2, 0x73, 0x0c, 0x58, 0x75, 0x4a, 0xf3, 0x04, 0xef, 0xc6, 0xfb, 0x27, 0xed, 0x53, 0x9f,
	0x9b, 0x20, 0xf8, 0x0b, 0x46, 0xae, 0xcb, 0xff, 0x1b, 0xfd, 0x23, 0xfc, 0x73, 0x38, 0xac, 0xf8,
	0xef, 0x52, 0x7a, 0x2b, 0x12, 0x64, 0xaf, 0xc0, 0xc3, 0xbb, 0x94, 0xa2, 0x6b, 0x91, 0x98, 0x0e,
	0x1d, 0x7e, 0x80, 0xae, 0xb8, 0x9b, 0xe6, 0x35, 0x0c, 0x1d, 0x4d, 0x28, 0x24, 0xfd, 0x5a, 0xb0,
	0x97, 0xd0, 0xfb, 0x24, 0x24, 0x45, 0x69, 0x61, 0xc7, 0xd9, 0xfd, 0x64, 0x0a, 0xbb, 0x09, 0x7e,
	0x01, 0xdf, 0x11, 0xfc, 0x29, 0x4a, 0x92, 0x6a, 0xbb, 0x0b, 0x75, 0xb0, 0x7f, 0xac, 0xd9, 0x6e,
	0x5d, 0xe2, 0xa6, 0xf0, 0x08, 0xd1, 0xe5, 0xc3, 0x4b, 0x16, 0x94, 0x88, 0x35, 0xb1, 0x17, 0xd0,
	0x95, 0xfa, 0xe4, 0x1e, 0x62, 0xa2, 0x27, 0xe1, 0xb1, 0x2c, 0x2d, 0x1e, 0xcb, 0xb2, 0x86, 0x57,
	0xf9, 0xdd, 0xf8, 0x04, 0xbc, 0xc5, 0xbd, 0x24, 0xcc, 0xde, 0x8b, 0xa5, 0xda, 0xa7, 0x4c, 0x2e,
	0xdd, 0x3e, 0x65, 0x72, 0xa9, 0x40, 0x2b, 0xbc, 0xc5, 0x95, 0xde, 0x27, 0x9f, 0x9b, 0x80, 0x9d,
	0x41, 0xf7, 0x63, 0x8a, 0xab, 0x44, 0x6a, 0xae, 0xfe, 0xec, 0xe5, 0xd6, 0x1f, 0xba, 0xd0, 0xae,
	0xc1, 0xed, 0xb5, 0xe0, 0x9f, 0x1e, 0x74, 0xe6, 0x6a, 0xd3, 0xd9, 0x10, 0xf6, 0xd2, 0xc4, 0x76,
	0xd8, 0x4b, 0x13, 0x36, 0x86, 0x5e, 0x4c, 0x84, 0x59, 0x41, 0xb6, 0x85, 0x0b, 0xd9, 0x39, 0x78,
	0x95, 0xdb, 0xd8, 0x3e, 0xff, 0xb9, 0x9a, 0xd5, 0x65, 0x36, 0x03, 0xa0, 0x58, 0xde, 0x44, 0xc6,
	0x6e, 0xbe, 0xd0, 0xd0, 0xa3, 0xba, 0xd2, 0xb4, 0xeb, 0x84, 0x2d, 0xee, 0x51, 0xe5, 0x4e, 0xaf,
	0x61, 0xe4, 0x30, 0x4e, 0x0e, 0x1d, 0x0d, 0x7c, 0xde, 0x00, 0x9a, 0xf5, 0x0f, 0x5b, 0xdc, 0xa7,
	0x86, 0x57, 0xe9, 0x44, 0x54, 0xed, 0x7a, 0x57, 0xc3, 0x8f, 0xeb, 0x70, 0x2b, 0x89, 0xb0, 0xc5,
	0xfb, 0x54, 0x33, 0xa7, 0x73, 0x18, 0x68, 0xa8, 0xf3, 0x96, 0xde, 0x36, 0xd2, 0x3a, 0x90, 0x43,
	0x3a, 0x43, 0xba, 0xb0, 0x4d, 0x2b, 0x5b, 0x39, 0xd0, 0xd0, 0x67, 0xbb, 0x6c, 0x25, 0x6c, 0x71,
	0xdd, 0xa6, 0xb2, 0xa3, 0x77, 0x70, 0x8c, 0x76, 0x75, 0xea, 0x9f, 0xed, 0x69, 0x8a, 0x2f, 0x1d,
	0xc5, 0x96, 0x69, 0x84, 0x2d, 0x7e, 0x84, 0x5b, 0x4e, 0x32, 0x87, 0x2a, 0xf9, 0x30, 0x02, 0xb0,
	0xdb, 0xd1, 0xa0, 0x7a, 0x18, 0xc3, 0x08, 0x1b, 0x66, 0x11, 0x02, 0x7b, 0xa0, 0xa9, 0x44, 0xdd,
	0xd7, 0x3c, 0xe3, 0x2d, 0x1e, 0x2b, 0xf2, 0xb0, 0xc5, 0x0f, 0xb1, 0xe9, 0x0a, 0x6f, 0xa0, 0xca,
	0x45, 0x4e, 0xd9, 0x03, 0xcd, 0xf3, 0xa2, 0xc9, 0x63, 0x2c, 0x20, 0x6c, 0xf1, 0x21, 0x6e, 0x9a,
	0xc2, 0x25, 0x54, 0x99, 0xc8, 0xa8, 0xdb, 0xdf, 0xdc, 0x89, 0x0d, 0x0f, 0x50, 0x3b, 0x81, 0x1b,
	0xa6, 0xf0, 0x23, 0x8c, 0x6a, 0x13, 0xd6, 0x9a, 0x1e, 0xee, 0x7e, 0x82, 0xd1, 0x7e, 0xfd, 0x09,
	0xd6, 0x0d, 0x1a, 0x14, 0x4a, 0xd6, 0xa3, 0x47, 0x29, 0xb0, 0x2c, 0x1b, 0x14, 0x4a, 0xf8, 0x33,
	0x00, 0xa9, 0x25, 0x1e, 0xad, 0xc4, 0x72, 0x7c, 0xb8, 0x29, 0x87, 0x4a, 0xfc, 0x4a, 0x0e, 0xd2,
	0x05, 0x6f, 0x7a, 0xd0, 0xd1, 0x17, 0x82, 0xe7, 0x70, 0xfc, 0xb6, 0xc4, 0x98, 0x50, 0xcb, 0x97,
	0xa3, 0x2c, 0x44, 0x2e, 0x71, 0xf6, 0x0e, 0x06, 0x3a, 0xb1, 0xc0, 0xf2, 0x36, 0xbd, 0x46, 0x76,
	0x01, 0xfd, 0xda, 0x35, 0xe6, 0x57, 0x8f, 0x53, 0x3f, 0x93, 0x57, 0x2e, 0xdc, 0x41, 0x15, 0xb4,
	0x3e, 0x74, 0xb5, 0x9c, 0xbf, 0xfd, 0x37, 0x00, 0x00, 0xff, 0xff, 0x92, 0xcd, 0xe4, 0xb5, 0x8d,
	0x08, 0x00, 0x00,
}