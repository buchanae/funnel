package worker

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





type RPCExecutorLogger struct {
  client
  taskID string
  executor int
  Stdout util.Tailer
  Stderr util.Tailer
}

func (r *RPCExecutorLogger) Close() {
  r.Stdout.Flush()
  r.Stderr.Flush()
}

func (r *RPCExecutorLogger) StartTime(t string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: r.executor,
    StartTime: t,
  })
}

func (r *RPCExecutorLogger) EndTime(t string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: r.executor,
    EndTime: t,
  })
}

func (r *RPCExecutorLogger) ExitCode(x int) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: r.executor,
    ExitCode: int32(x),
  })
}

func (r *RPCExecutorLogger) HostIP(ip string) {
  r.client.UpdateExecutorLogs({
    TaskId: r.taskID,
    Executor: r.executor,
    HostIP: ip,
  })
}
