package worker

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/logger"
  "io"
)

type TaskLogger interface {
	StartTime(string)
	EndTime(string)
	Outputs(string)
	Metadata(map[string]string)
  Running()
  Result(error)

  ExecutorExitCode(int, int)
  ExecutorPorts(int, []*tes.Ports)
  ExecutorHostIP(int, string)
  ExecutorStartTime(int, string)
  ExecutorEndTime(int, string)
  ExecutorStdout(int) io.Writer
  ExecutorStderr(int) io.Writer
}

type TaskViewer interface {
  Task() (*tes.Task, error)
  State() tes.State
}

type Executor interface {
  Execute(context.Context, int) error
}

type Runner interface {
  Run(context.Context)
}

type Backend interface {
  logger.Logger
	storage.Storage
	TaskLogger
  TaskViewer
  Executor
  Runner

  Close()
}
