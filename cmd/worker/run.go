package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/server/elastic"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/worker"
	"github.com/ohsu-comp-bio/funnel/util"
)

// Run configures and runs a Worker
func Run(ctx context.Context, conf config.Worker, taskID string, log *logger.Logger) error {
  w := NewDefaultWorker(conf, taskID, log)
  worker.Run(ctx, w, taskID)
	return nil
}

// NewDefaultWorker returns a new configured DefaultWorker instance.
func NewDefaultWorker(conf config.Worker, taskID string, log *logger.Logger) *DefaultWorker {
	return &DefaultWorker{conf, taskID, log}
}

// DefaultWorker is the default task worker, which follows a basic,
// sequential process of task initialization, execution, finalization,
// and logging.
type DefaultWorker struct {
	conf       config.Worker
  taskID string
  log *logger.Logger
}

func (d *DefaultWorker) Config() config.Worker {
  return d.conf
}

func (d *DefaultWorker) TaskReader() (worker.TaskReader, error) {
	switch d.conf.TaskReader {
	case "rpc":
    return worker.NewRPCTaskReader(d.conf, d.taskID)
	case "dynamodb":
    return worker.NewDynamoDBTaskReader(d.conf.TaskReaders.DynamoDB, d.taskID)
  default:
    return nil, fmt.Errorf("unknown task reader: %s", d.conf.TaskReader)
	}
}

func (d *DefaultWorker) EventWriter() (events.Writer, error) {
  var errs util.MultiError
  var writers []events.Writer

	for _, w := range d.conf.ActiveEventWriters {
	  var writer events.Writer
    var err error

		switch w {
		case "dynamodb":
			writer, err = events.NewDynamoDBEventWriter(d.conf.EventWriters.DynamoDB)
		case "log":
			writer = &events.Logger{Log: d.log}
		case "rpc":
			writer, err = events.NewRPCWriter(d.conf)
		case "elastic":
			writer, err = elastic.NewElastic(d.conf.EventWriters.Elastic)
		default:
			err = fmt.Errorf("unknown EventWriter: %s", w)
		}

		if err != nil {
      errs = append(errs, err)
		} else {
		  writers = append(writers, writer)
    }
	}

  if writers == nil {
    return nil, fmt.Errorf("no event writers configured")
  }

  return events.MultiWriter(writers...), errs
}

func (d *DefaultWorker) FileMapper(task *tes.Task) (*worker.FileMapper, error) {
  return worker.MapTask(task, d.conf.WorkDir)
}

func (d *DefaultWorker) Storage(t *tes.Task) (storage.Storage, error) {
  return storage.Storage{}.WithConfig(d.conf.Storage)
}

func (d *DefaultWorker) Executor(task *tes.Task, index int) worker.Executor {
  e := task.Executors[i]
  return DockerCmd{
    ImageName:     e.ImageName,
    Cmd:           e.Cmd,
    Environ:       e.Environ,
    Volumes:       mapper.Volumes,
    Workdir:       e.Workdir,
    Ports:         e.Ports,
    ContainerName: fmt.Sprintf("%s-%d", task.Id, i),
    // TODO make RemoveContainer configurable
    RemoveContainer: true,
  }
}

type Mapper struct {
  Storage
  *worker.FileMapper
}
func (m *Mapper) Get(ctx context.Context, url string, path string, class tes.FileType) error {
}
func (m *Mapper) Put(ctx context.Context, url string, path string, class tes.FileType) ([]*tes.OutputFileLog, error) {
}
func (m *Mapper) Supports(url string, path string, class tes.FileType) bool {
}

type LinkFixer struct {
  Storage
  *worker.FileMapper
}
func (m *LinkFixer) Get(ctx context.Context, url string, path string, class tes.FileType) error {
}
// fixLinks walks the output paths, fixing cases where a symlink is
// broken because it's pointing to a path inside a container volume.
func (m *LinkFixer) Put(ctx context.Context, url string, path string, class tes.FileType) ([]*tes.OutputFileLog, error) {
	filepath.Walk(basepath, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			// There's an error, so be safe and give up on this file
			return nil
		}

		// Only bother to check symlinks
		if f.Mode()&os.ModeSymlink != 0 {
			// Test if the file can be opened because it doesn't exist
			fh, rerr := os.Open(p)
			fh.Close()

			if rerr != nil && os.IsNotExist(rerr) {

				// Get symlink source path
				src, err := os.Readlink(p)
				if err != nil {
					return nil
				}
				// Map symlink source (possible container path) to host path
				mapped, err := m(src)
				if err != nil {
					return nil
				}

				// Check whether the mapped path exists
				fh, err := os.Open(mapped)
				fh.Close()

				// If the mapped path exists, fix the symlink
				if err == nil {
					err := os.Remove(p)
					if err != nil {
						return nil
					}
					os.Symlink(mapped, p)
				}
			}
		}
		return nil
	})
}
func (m *LinkFixer) Supports(url string, path string, class tes.FileType) bool {
}
