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
	ExitCode
	HostIp
	Ports
	Stdout
	Stderr
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

type ExitCode struct {
	ExitCode int32  `protobuf:"varint,1,opt,name=exit_code,json=exitCode" json:"exit_code,omitempty"`
	Index    uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *ExitCode) Reset()                    { *m = ExitCode{} }
func (m *ExitCode) String() string            { return proto.CompactTextString(m) }
func (*ExitCode) ProtoMessage()               {}
func (*ExitCode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ExitCode) GetExitCode() int32 {
	if m != nil {
		return m.ExitCode
	}
	return 0
}

func (m *ExitCode) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type HostIp struct {
	HostIp string `protobuf:"bytes,1,opt,name=host_ip,json=hostIp" json:"host_ip,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *HostIp) Reset()                    { *m = HostIp{} }
func (m *HostIp) String() string            { return proto.CompactTextString(m) }
func (*HostIp) ProtoMessage()               {}
func (*HostIp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *HostIp) GetHostIp() string {
	if m != nil {
		return m.HostIp
	}
	return ""
}

func (m *HostIp) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type Ports struct {
	Ports []*tes.Ports `protobuf:"bytes,1,rep,name=ports" json:"ports,omitempty"`
	Index uint32       `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *Ports) Reset()                    { *m = Ports{} }
func (m *Ports) String() string            { return proto.CompactTextString(m) }
func (*Ports) ProtoMessage()               {}
func (*Ports) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Ports) GetPorts() []*tes.Ports {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *Ports) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type Stdout struct {
	Stdout string `protobuf:"bytes,1,opt,name=stdout" json:"stdout,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *Stdout) Reset()                    { *m = Stdout{} }
func (m *Stdout) String() string            { return proto.CompactTextString(m) }
func (*Stdout) ProtoMessage()               {}
func (*Stdout) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Stdout) GetStdout() string {
	if m != nil {
		return m.Stdout
	}
	return ""
}

func (m *Stdout) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type Stderr struct {
	Stderr string `protobuf:"bytes,1,opt,name=stderr" json:"stderr,omitempty"`
	Index  uint32 `protobuf:"varint,3,opt,name=index" json:"index,omitempty"`
}

func (m *Stderr) Reset()                    { *m = Stderr{} }
func (m *Stderr) String() string            { return proto.CompactTextString(m) }
func (*Stderr) ProtoMessage()               {}
func (*Stderr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Stderr) GetStderr() string {
	if m != nil {
		return m.Stderr
	}
	return ""
}

func (m *Stderr) GetIndex() uint32 {
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
	//	*Event_ExitCode
	//	*Event_HostIp
	//	*Event_Ports
	//	*Event_Stdout
	//	*Event_Stderr
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
type Event_ExitCode struct {
	ExitCode *ExitCode `protobuf:"bytes,11,opt,name=exit_code,json=exitCode,oneof"`
}
type Event_HostIp struct {
	HostIp *HostIp `protobuf:"bytes,12,opt,name=host_ip,json=hostIp,oneof"`
}
type Event_Ports struct {
	Ports *Ports `protobuf:"bytes,13,opt,name=ports,oneof"`
}
type Event_Stdout struct {
	Stdout *Stdout `protobuf:"bytes,14,opt,name=stdout,oneof"`
}
type Event_Stderr struct {
	Stderr *Stderr `protobuf:"bytes,15,opt,name=stderr,oneof"`
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
func (*Event_ExitCode) isEvent_Event()          {}
func (*Event_HostIp) isEvent_Event()            {}
func (*Event_Ports) isEvent_Event()             {}
func (*Event_Stdout) isEvent_Event()            {}
func (*Event_Stderr) isEvent_Event()            {}
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

func (m *Event) GetExitCode() *ExitCode {
	if x, ok := m.GetEvent().(*Event_ExitCode); ok {
		return x.ExitCode
	}
	return nil
}

func (m *Event) GetHostIp() *HostIp {
	if x, ok := m.GetEvent().(*Event_HostIp); ok {
		return x.HostIp
	}
	return nil
}

func (m *Event) GetPorts() *Ports {
	if x, ok := m.GetEvent().(*Event_Ports); ok {
		return x.Ports
	}
	return nil
}

func (m *Event) GetStdout() *Stdout {
	if x, ok := m.GetEvent().(*Event_Stdout); ok {
		return x.Stdout
	}
	return nil
}

func (m *Event) GetStderr() *Stderr {
	if x, ok := m.GetEvent().(*Event_Stderr); ok {
		return x.Stderr
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
		(*Event_ExitCode)(nil),
		(*Event_HostIp)(nil),
		(*Event_Ports)(nil),
		(*Event_Stdout)(nil),
		(*Event_Stderr)(nil),
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
	case *Event_ExitCode:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ExitCode); err != nil {
			return err
		}
	case *Event_HostIp:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.HostIp); err != nil {
			return err
		}
	case *Event_Ports:
		b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ports); err != nil {
			return err
		}
	case *Event_Stdout:
		b.EncodeVarint(14<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Stdout); err != nil {
			return err
		}
	case *Event_Stderr:
		b.EncodeVarint(15<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Stderr); err != nil {
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
	case 11: // event.exit_code
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ExitCode)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ExitCode{msg}
		return true, err
	case 12: // event.host_ip
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HostIp)
		err := b.DecodeMessage(msg)
		m.Event = &Event_HostIp{msg}
		return true, err
	case 13: // event.ports
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Ports)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Ports{msg}
		return true, err
	case 14: // event.stdout
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Stdout)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Stdout{msg}
		return true, err
	case 15: // event.stderr
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Stderr)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Stderr{msg}
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
	case *Event_ExitCode:
		s := proto.Size(x.ExitCode)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_HostIp:
		s := proto.Size(x.HostIp)
		n += proto.SizeVarint(12<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Ports:
		s := proto.Size(x.Ports)
		n += proto.SizeVarint(13<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Stdout:
		s := proto.Size(x.Stdout)
		n += proto.SizeVarint(14<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Stderr:
		s := proto.Size(x.Stderr)
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
	proto.RegisterType((*ExitCode)(nil), "events.ExitCode")
	proto.RegisterType((*HostIp)(nil), "events.HostIp")
	proto.RegisterType((*Ports)(nil), "events.Ports")
	proto.RegisterType((*Stdout)(nil), "events.Stdout")
	proto.RegisterType((*Stderr)(nil), "events.Stderr")
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
	// 808 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x7f, 0x6f, 0xe3, 0x44,
	0x10, 0x75, 0x5a, 0xe2, 0xc4, 0x93, 0xa4, 0x69, 0xb7, 0x77, 0x5c, 0xc8, 0xf1, 0x47, 0x65, 0x09,
	0x29, 0x48, 0x90, 0x4a, 0x41, 0x40, 0xa1, 0x82, 0x93, 0x38, 0x82, 0x0c, 0x77, 0x08, 0xb4, 0xe9,
	0xdf, 0x44, 0xbe, 0x7a, 0x2e, 0x67, 0x35, 0xf6, 0x5a, 0xde, 0x49, 0x94, 0x7e, 0x06, 0xbe, 0x11,
	0x9f, 0x0e, 0xed, 0x2f, 0xc7, 0x4d, 0x72, 0x15, 0x12, 0xf7, 0xdf, 0xce, 0xce, 0xbc, 0x37, 0xf1,
	0xe6, 0xbd, 0x07, 0x5d, 0x5c, 0x63, 0x4e, 0x72, 0x5c, 0x94, 0x82, 0x04, 0xf3, 0x4d, 0x35, 0x0c,
	0x08, 0xed, 0xd5, 0xb0, 0x2b, 0xa9, 0x5c, 0xdd, 0x92, 0xad, 0xfa, 0x94, 0x66, 0x28, 0x29, 0xce,
	0x0a, 0x7b, 0xf1, 0xe9, 0x42, 0x88, 0xc5, 0x12, 0x2f, 0xe3, 0x22, 0xbd, 0x8c, 0xf3, 0x5c, 0x50,
	0x4c, 0xa9, 0xc8, 0x2d, 0x38, 0xfc, 0x12, 0x82, 0x9b, 0x58, 0xde, 0xcd, 0x28, 0x26, 0x64, 0x17,
	0xd0, 0x94, 0xea, 0x30, 0x68, 0x5c, 0x34, 0x46, 0x27, 0x13, 0x18, 0xab, 0x25, 0xba, 0xc5, 0x4d,
	0x23, 0xfc, 0x0d, 0x7a, 0x76, 0xbc, 0xa4, 0x9b, 0x34, 0x43, 0xf6, 0x1d, 0x80, 0x54, 0xc5, 0x5c,
	0xad, 0xd5, 0xb8, 0xce, 0x64, 0x38, 0x36, 0x2b, 0xcd, 0x8a, 0x37, 0xab, 0xb7, 0xe3, 0x1b, 0xf7,
	0x9b, 0x78, 0x20, 0x1d, 0x34, 0xfc, 0x19, 0x3a, 0x8a, 0x6b, 0x9a, 0x27, 0x9a, 0xe9, 0x6b, 0x68,
	0x63, 0x9e, 0xfc, 0x57, 0x9e, 0x16, 0x1a, 0x58, 0x78, 0x6d, 0x58, 0xfe, 0x58, 0x51, 0xb1, 0x22,
	0xc9, 0xbe, 0x80, 0x96, 0x30, 0xc7, 0x41, 0xe3, 0xe2, 0x78, 0xd4, 0x99, 0x30, 0xfd, 0x11, 0xa6,
	0xfd, 0x4b, 0xba, 0xc4, 0xd7, 0x62, 0xc1, 0xdd, 0x48, 0xf8, 0x77, 0x03, 0xba, 0x0a, 0xfd, 0x3b,
	0x52, 0x9c, 0xc4, 0x14, 0xb3, 0x1f, 0xa1, 0x9d, 0xd9, 0xb3, 0xc5, 0x87, 0x63, 0xfb, 0xfe, 0xf5,
	0xb9, 0xb1, 0x3b, 0x4c, 0x73, 0x2a, 0xef, 0x79, 0x85, 0x19, 0x5e, 0x43, 0xef, 0x41, 0x8b, 0x9d,
	0xc2, 0xf1, 0x1d, 0xde, 0xeb, 0x0f, 0x0a, 0xb8, 0x3a, 0xb2, 0x27, 0xd0, 0x5c, 0xc7, 0xcb, 0x15,
	0x0e, 0x8e, 0xf4, 0x9d, 0x29, 0xbe, 0x3f, 0xba, 0x6a, 0x84, 0x09, 0x9c, 0x4d, 0x37, 0x78, 0xbb,
	0x22, 0x51, 0x7e, 0x88, 0x07, 0x56, 0x9b, 0xd2, 0x3c, 0xc1, 0xcd, 0xe0, 0xf8, 0xa2, 0x31, 0xea,
	0x71, 0x53, 0x84, 0x7f, 0x41, 0xdf, 0x6d, 0xf9, 0x7f, 0x4f, 0xff, 0x1e, 0xfe, 0x1f, 0xa0, 0x3d,
	0xdd, 0xa4, 0xf4, 0x52, 0x24, 0xc8, 0x9e, 0x43, 0x80, 0x9b, 0x94, 0xe6, 0xb7, 0x22, 0x31, 0xcc,
	0x4d, 0xde, 0x46, 0xd7, 0x3c, 0x0c, 0xff, 0x16, 0xfc, 0x48, 0x48, 0xfa, 0xb5, 0x60, 0xcf, 0xa0,
	0xf5, 0x4e, 0x48, 0x9a, 0xa7, 0x85, 0x7d, 0x3e, 0xff, 0x9d, 0x69, 0x1c, 0x06, 0xbe, 0x80, 0xe6,
	0x9f, 0xa2, 0x24, 0xa9, 0x54, 0x5c, 0xa8, 0x83, 0xfd, 0x03, 0x8d, 0x8a, 0x75, 0x8b, 0x9b, 0xc6,
	0x7b, 0x08, 0xbe, 0x01, 0x7f, 0x46, 0x89, 0x58, 0x11, 0xfb, 0x18, 0x7c, 0xa9, 0x4f, 0x6e, 0xb1,
	0xa9, 0x1e, 0xc5, 0x61, 0x59, 0x5a, 0x1c, 0x96, 0x65, 0x0d, 0xa7, 0xee, 0x0f, 0xe3, 0x12, 0x08,
	0x66, 0xf7, 0x92, 0x30, 0x7b, 0x2d, 0x16, 0x4a, 0x27, 0x99, 0x5c, 0x38, 0x9d, 0x64, 0x72, 0xa1,
	0x40, 0x4b, 0x5c, 0xe3, 0x52, 0xeb, 0xa4, 0xc7, 0x4d, 0xc1, 0x2e, 0xc1, 0x7f, 0x9b, 0xe2, 0x32,
	0x91, 0x9a, 0xab, 0x33, 0x79, 0xb6, 0xf7, 0x47, 0xcd, 0x74, 0x1a, 0x70, 0x3b, 0x16, 0xfe, 0xe3,
	0x43, 0x73, 0xaa, 0x14, 0xcc, 0x4e, 0xe0, 0x28, 0x4d, 0xec, 0x86, 0xa3, 0x34, 0x61, 0x03, 0x68,
	0xc5, 0x44, 0x98, 0x15, 0x64, 0x57, 0xb8, 0x92, 0x5d, 0x41, 0x50, 0xa5, 0x88, 0xdd, 0xf3, 0xa8,
	0xe4, 0xaa, 0x61, 0x36, 0x01, 0xa0, 0x58, 0xde, 0xcd, 0x4d, 0x8c, 0x7c, 0xa4, 0xa1, 0x67, 0x75,
	0x07, 0xe9, 0x34, 0x89, 0x3c, 0x1e, 0x50, 0x95, 0x3a, 0x2f, 0xa0, 0xef, 0x30, 0x4e, 0xe6, 0x4d,
	0x0d, 0x7c, 0xba, 0x03, 0x34, 0xb2, 0x8e, 0x3c, 0xde, 0xa3, 0x9d, 0x0c, 0xd2, 0x17, 0xf3, 0x4a,
	0xc3, 0xbe, 0x86, 0x9f, 0xd7, 0xe1, 0x56, 0xea, 0x91, 0xc7, 0x3b, 0x54, 0x0b, 0x9d, 0x2b, 0xe8,
	0x6a, 0xa8, 0xcb, 0x8c, 0xd6, 0x3e, 0xd2, 0x26, 0x8b, 0x43, 0xba, 0xa0, 0xb9, 0xb6, 0x4b, 0xab,
	0xb8, 0x68, 0x6b, 0xe8, 0x93, 0x43, 0x71, 0x11, 0x79, 0x5c, 0xaf, 0xa9, 0x62, 0xe6, 0x15, 0x9c,
	0xa3, 0xf5, 0x60, 0xfd, 0xb3, 0x03, 0x4d, 0xf1, 0x89, 0xa3, 0xd8, 0x0b, 0x83, 0xc8, 0xe3, 0x67,
	0xb8, 0x97, 0x10, 0x53, 0xa8, 0x2e, 0xb7, 0x4f, 0x00, 0x56, 0x1d, 0x3b, 0x54, 0xdb, 0x67, 0xe8,
	0xe3, 0x4e, 0x08, 0x5c, 0xd6, 0xbd, 0xda, 0xd1, 0xf0, 0xd3, 0x2d, 0xdc, 0x78, 0x36, 0xf2, 0x6a,
	0xfe, 0xfd, 0x7c, 0xeb, 0xcf, 0xae, 0x1e, 0x3f, 0x71, 0xe3, 0xc6, 0xc0, 0x91, 0x57, 0x39, 0xf6,
	0x33, 0x67, 0xc9, 0x9e, 0x1e, 0xec, 0xb9, 0x41, 0xed, 0xca, 0xc8, 0x73, 0xbe, 0x1c, 0x55, 0xbe,
	0x3b, 0x79, 0x48, 0x68, 0x7c, 0xa9, 0x08, 0xad, 0x13, 0x47, 0x95, 0xd3, 0xfa, 0x7b, 0x93, 0x58,
	0x96, 0x76, 0x52, 0x79, 0x6f, 0x02, 0x20, 0xb5, 0xcb, 0xe6, 0x4b, 0xb1, 0x18, 0x9c, 0x3e, 0x54,
	0x64, 0xe5, 0x3f, 0xa5, 0x48, 0xe9, 0x8a, 0x9f, 0x5a, 0xd0, 0xd4, 0x03, 0xe1, 0x53, 0x38, 0x7f,
	0x59, 0x62, 0x4c, 0xa8, 0x1d, 0xc4, 0x51, 0x16, 0x22, 0x97, 0x38, 0x79, 0x05, 0x5d, 0x7d, 0x31,
	0xc3, 0x72, 0x9d, 0xde, 0x22, 0xbb, 0x86, 0x4e, 0x6d, 0x8c, 0x55, 0x9f, 0xa7, 0xcb, 0xe1, 0x73,
	0x57, 0x1e, 0xa0, 0x0a, 0xbd, 0x37, 0xbe, 0x76, 0xd4, 0x57, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff,
	0x3b, 0x40, 0xb7, 0x21, 0xe8, 0x07, 0x00, 0x00,
}
