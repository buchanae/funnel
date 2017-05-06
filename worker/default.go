package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
)


func NewDefaultBackend(conf config.Worker, taskID string) (*DefaultBackend, error) {
  store, serr := storage.FromConfig(conf.Storage)
  rpc, rerr := NewRPCTask(conf, taskID)
  task, terr := rpc.Task()
  mapped, merr := MapTaskFiles(conf.WorkDir, task)

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
    &DockerFactory{
      rpc,
      task,
      mapped,
      conf,
    },
  }, nil
}

type DefaultBackend struct {
  logger.Logger
  *RPCTask
  storage.Storage
  *DefaultTaskRunner
  *DockerFactory
}

func (b *DefaultBackend) Close() {
  // TODO ?? b.RPCTaskLogger.Close()
}
