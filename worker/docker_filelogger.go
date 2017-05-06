package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
)

func NewFileBackend(conf config.Worker, taskID string) (*FileBackend, error) {
  store, serr := storage.FromConfig(conf.Storage)
  filetask, fterr := NewFileTask(conf, taskID)
  task, terr := filetask.Task()
  mapped, merr := MapTaskFiles(conf.WorkDir, task)

  if err := util.Check(serr, fterr, terr, merr); err != nil {
    return nil, err
  }

  return &FileBackend{
    log.WithFields("task", taskID),
    filetask,
    storage,
    &DefaultTaskRunner{
      storage,
      filetask,
      filetask,
      conf.PollRate,
    },
    &DockerFactory{
      filetask,
      task,
      mapped,
      conf,
    },
  }, nil
}

type FileBackend struct {
  logger.Logger
  *FileTask
  storage.Storage
  *DefaultTaskRunner
  *DockerFactory
}

func (b *FileBackend) Close() {}
