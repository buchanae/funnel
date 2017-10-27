package worker

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/config"
)

type Executor interface {
  Run() error
  String() string
  Inspect(context.Context) ([]*tes.Ports, error)
  Stop() error
}

type Factory interface {
  Storage(*tes.Task) (storage.Storage, error)
  EventWriter() (events.Writer, error)
  TaskReader() (TaskReader, error)
  Config() config.Worker
  FileMapper(*tes.Task) (*FileMapper, error)
  Executor() Executor
}

// TaskReader is a type which reads task information
// during task execution.
type TaskReader interface {
	Task() (*tes.Task, error)
	State() (tes.State, error)
}
