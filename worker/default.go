package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/worker/mapper"
  "path/filepath"
)

func NewDefaultBackend(conf config.Worker, taskID string) (*DefaultBackend, error) {
  base := filepath.Join(conf.WorkDir, taskID)
  rpc, rerr := NewRPCTask(conf, taskID)
  task, terr := rpc.Task()
  mapped, merr := mapper.MapTask(base, task)
  store, serr := storage.FromConfig(conf.Storage)
  store = mapper.MapStorage(base, store)

  if err := util.Check(serr, rerr, terr, merr); err != nil {
    return nil, err
  }

  return &DefaultBackend{
    log.WithFields("task", taskID),
    rpc,
    storage,
    &DefaultTaskRunner{
      storage,
      rpc,
      rpc,
      conf.PollRate,
    },
    &DockerExecutor{
      rpc,
      task,
      mapped,
      true,
    },
  }, nil
}

type DefaultBackend struct {
  logger.Logger
  *RPCTask
  storage.Storage
  *DefaultRunner
  *DockerExecutor
}

func (b *DefaultBackend) Close() {
  // TODO ?? b.RPCTaskLogger.Close()
}
