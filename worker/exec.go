package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
  "io"
)


// ExecConfig describes a task executor's container, command, environment, etc.
type Config struct {
	ImageName       string
	Cmd             []string
	Volumes         []Volume
	Workdir         string
	Ports           []*tes.Ports
	ContainerName   string
	RemoveContainer bool
	Environ         map[string]string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
}

type ExecutorFactory(Config) Executor

type Executor interface {
  Run(ctx context.Context) (exitCode int64, err error)
  Inspect(context.Context) ([]*tes.Ports, error)
  Stop(context.Context) error
}
