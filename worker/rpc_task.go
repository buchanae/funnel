package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
	"time"
)

// TODO document behavior of slow consumer of task log updates

func NewRPCTask(conf config.Config, taskID string) (*RPCTask, error) {
  reader := &RPCTaskReader{client, taskID}
  logger := &RPCTaskLogger{client, taskID}
}

type RPCTask struct {
  *RPCTaskReader
  *RPCTaskLogger
}

type RPCTaskReader struct {
  taskID string
}

func (r *RPCTaskReader) Task() (*tes.Task, error) {
  task, terr := r.client.GetTask(context.TODO(), &tes.GetTaskRequest{
    Id: r.taskID,
    View: tes.TaskView_FULL,
  })
}

func (r *RPCTaskReader) State() (*tes.State, error) {
  task, terr := r.client.GetTask(context.TODO(), &tes.GetTaskRequest{
    Id: r.taskID,
  })
}


type RPCTaskLogger struct {
  client
  taskID string
}

func (r *RPCTaskLogger) StartTime(t string) {
  r.client.UpdateTaskLogs({
    Id: r.taskID,
    StartTime: t,
  })
}

func (r *RPCTaskLogger) EndTime(t string) {
  r.client.UpdateTaskLogs({
    Id: r.taskID,
    EndTime: t,
  })
}

func (r *RPCTaskLogger) OutputFiles(f []string) {
  r.client.UpdateTaskLogs({
    Id: r.taskID,
    EndTime: t,
  })
}

func (r *RPCTaskLogger) Metadata(m map[string]string) {
  r.client.UpdateTaskLogs({
    Id: r.taskID,
    EndTime: t,
  })
}

func (r *RPCTaskLogger) Running() {
  r.client.UpdateTaskState({
    Id: r.taskID,
    State: tes.State_RUNNING,
  })
}

func (r *RPCTaskLogger) Result(err error) {
  r.client.UpdateTaskState({
    Id: r.taskID,
    State: tes.State_RUNNING,
  })
}

func (r *RPCTaskLogger) Close() {}



func (r *RPCTaskLogger) ExecutorStartTime(i int, t string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: i,
    StartTime: t,
  })
}

func (r *RPCTaskLogger) ExecutorEndTime(i int, t string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: i,
    EndTime: t,
  })
}

func (r *RPCTaskLogger) ExecutorExitCode(i int, x int) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: i,
    ExitCode: int32(x),
  })
}

func (r *RPCTaskLogger) ExecutorHostIP(i int, ip string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: i,
    HostIP: ip,
  })
}

func (r *RPCTaskLogger) ExecutorStdout(i int) io.Writer {
  // tailer
    // TODO
    //Stdout: io.MultiWriter(stdout, log.Stdout())
    //Stderr: io.MultiWriter(stderr, log.Stderr())
}

func (r *RPCTaskLogger) ExecutorStderr(i int) io.Writer {
  // tailer
}
