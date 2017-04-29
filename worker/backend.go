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
  stdout util.Tailer
  stderr util.Tailer
}

func (r *RPCExecutorLogger) Close() {
  r.stdout.Flush()
  r.stderr.Flush()
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

func (r *RPCExecutorLogger) Stdout() io.Writer {
  return r.stdout
}

func (r *RPCExecutorLogger) Stderr() io.Writer {
  return r.stderr
}





func NewDefaultBackend(conf config.Worker, taskID string) (*DefaultBackend, error) {
	// Map files into this baseDir
	baseDir := path.Join(conf.WorkDir, t.Task.Id)
	prepareDir(baseDir)
  NewTaskState()

  task, err := b.client.GetTask(ctx, &tes.GetTaskRequest{
    Id: taskID:
    View: tes.TaskView_FULL,
  })

  return &DefaultBackend{
  }, nil
}

type DefaultBackend struct {
  storage.Storage
  logger.Logger
  // TODO does the backend need to be rpc specific?
  *RPCTaskLogger
  task *tes.Task
  client
}

func (b *DefaultBackend) Task() *tes.Task {
  return b.task
}

func (b *DefaultBackend) Close() {
  b.client.Close()
}

func (b *DefaultBackend) WithContext(ctx context.Context) context.Context {
  taskctx, cancel := context.WithCancel(ctx)

  go func() {
    ticker := time.NewTicker(b.conf.UpdateRate)
    defer ticker.Stop()

    for {
    case <-taskctx.Done():
      return
    case <-ticker.C:
      task, err := b.client.GetTask(ctx, &tes.GetTaskRequest{
        Id: b.task.Id:
      })

      if task.State == tes.State_CANCELED {
        cancel()
      }
    }
  }()
  return taskctx
}


func (b *DefaultBackend) Executor(i int, d *tes.Executor) Executor {
  return &DefaultBackendExecutor{
    &DockerExecutor{
      ImageName:       d.ImageName,
      Cmd:             d.Cmd,
      Volumes:         r.mapper.Volumes,
      Workdir:         d.Workdir,
      Ports:           d.Ports,
      ContainerName:   fmt.Sprintf("%s-%d", task.Id, i),
      RemoveContainer: r.conf.RemoveContainer,
      Environ:         d.Environ,
    },
    &RPCExecutorLogger{
      client: b.client,
      taskID: b.task.Id,
      executor: i,
    },
  }
}


type DefaultBackendExecutor struct {
  ExecutorLogger
  *Docker
}

func (b *DefaultBackendExecutor) Run(ctx context.Context) error {
  var err error
	if d.Stdin != "" {
		exec.Stdin, err = b.mapper.OpenHostFile(d.Stdin)
	}
  exec.Stdout, err = b.teeLogFile(d.Stdout, log.Stdout(i))
  exec.Stderr, err = b.teeLogFile(d.Stderr, log.Stderr(i))
  return b.Docker.Run(ctx)
}

func (b *DefaultBackendExecutor) Close() {
  b.ExecutorLogger.Close()
  // TODO ?b.DockerExecutor.Close()
}

func (b *DefaultBackend) teeLogFile(p string, rpc io.Writer) (io.Writer, error) {
  if p == "" {
    nil, nil
  }
  f, err := b.mapper.CreateHostFile(p)
  if err != nil {
    return nil, err
  }
  return io.MultiWriter(f, rpc), nil
}




// Create working dir
func prepareDir(path string) error {
	dir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return util.EnsureDir(dir)
}

// Configure a task-specific storage backend.
// This provides download/upload for inputs/outputs.
func (r *taskRunner) prepareStorage() error {
	var err error

	for _, conf := range r.conf.Storage {
		r.store, err = r.store.WithConfig(conf)
		if err != nil {
			return err
		}
	}

	return nil
}
