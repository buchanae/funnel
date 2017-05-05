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
  workspace, werr := NewWorkspace(conf.WorkDir, task)
  store, serr := storage.FromConfig(conf.Storage)
  rpc, err := NewRPCTask(conf, taskID)
  docker := DockerExecutor{
    RemoveContainer: conf.RemoveContainer,
    task: task,
    logger: filetask,
    workspace: workspace,
  }

  if err := util.Check(werr, terr, serr); err != nil {
    return nil, err
  }

  return &DefaultBackend{
    Logger: log.WithFields("task", taskID),
    RPCTaskLogger: rpc,
    RPCTaskReader: rpc,
    Storage: store,
    DockerExecutor: docker,
  }, nil
}

type DefaultBackend struct {
  logger.Logger
  *RPCTaskLogger
  *RPCTaskReader
  storage.Storage
  *DockerExecutor
}

func (b *DefaultBackend) Close() {
  // TODO ?? b.RPCTaskLogger.Close()
}
