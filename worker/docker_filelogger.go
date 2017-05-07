package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/worker/mapper"
  "path/filepath"
)

func NewFileBackend(conf config.Worker, taskID string) (*FileBackend, error) {
  base := filepath.Join(conf.WorkDir, taskID)
  filetask, fterr := NewFileTask(conf, taskID)
  task, terr := filetask.Task()
  mapped, merr := mapper.MapTask(base, task)
  store, serr := storage.FromConfig(conf.Storage)
  store = mapper.MapStorage(base, store)

  if err := util.Check(serr, fterr, terr, merr); err != nil {
    return nil, err
  }

  return &FileBackend{
    log.WithFields("task", taskID),
    filetask,
    storage,
    &DefaultRunner{
      storage,
      filetask,
      filetask,
      conf.PollRate,
    },
    &DockerExecutor{
      filetask,
      task,
      mapped,
      true,
    },
  }, nil
}

type FileBackend struct {
  logger.Logger
  *FileTask
  storage.Storage
  *DefaultRunner
  *DockerExecutor
}

func (b *FileBackend) Close() {}
