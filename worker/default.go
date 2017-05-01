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

func NewDefaultBackend(conf config.Worker, taskID string) (*DefaultBackend, error) {
  workspace, werr := NewWorkspace(conf.WorkDir, taskID)
  NewTaskState()

  task, terr := client.GetTask(context.TODO(), &tes.GetTaskRequest{
    Id: taskID:
    View: tes.TaskView_FULL,
  })

  store, serr := storage.FromConfig(conf.Storage)

  if err := util.Check(werr, terr, serr); err != nil {
    return nil, err
  }

  return &DefaultBackend{
    Logger: log.WithFields("task", taskID),
    RPCTaskLogger: &RPCTaskLogger{client, taskID},
    Storage: store,
    task: task,
    client: client
    workspace: workspace,
  }, nil
}

type DefaultBackend struct {
  logger.Logger
  *RPCTaskLogger
  storage.Storage
  task *tes.Task
  client
  workspace *Workspace
}

func (b *DefaultBackend) Task() *tes.Task {
  return b.task
}

func (b *DefaultBackend) Close() {
  b.client.Close()
}

func (b *DefaultBackend) Executor(i int, d *tes.Executor) Executor {
  log := &RPCExecutorLogger{
    client: b.client,
    taskID: b.task.Id,
    executor: i,
  }

  stdin, ierr := b.workspace.Reader(d.Stdin)
  stdout, oerr := b.workspace.Writer(d.Stdout)
  stderr, eerr := b.workspace.Writer(d.Stderr)

  if err := util.Check(ierr, oerr, eerr); err != nil {
    return nil, err
  }

  return &Docker{
    log,
    ImageName:       d.ImageName,
    Cmd:             d.Cmd,
    Volumes:         r.mapper.Volumes,
    Workdir:         d.Workdir,
    Ports:           d.Ports,
    ContainerName:   fmt.Sprintf("%s-%d", task.Id, i),
    RemoveContainer: r.conf.RemoveContainer,
    Environ:         d.Environ,
    Stdin: stdin,
    Stdout: io.MultiWriter(stdout, log.Stdout())
    Stderr: io.MultiWriter(stderr, log.Stderr())
  }, nil
}
