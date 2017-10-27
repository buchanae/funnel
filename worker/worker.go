package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/cmd/version"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/storage"
	"os"
	"path/filepath"
	"time"
)

// Run runs the Worker.
// TODO document behavior of slow consumer of task log updates
func Run(pctx context.Context, w Worker, taskID string) {

	// The code here is verbose, but simple; mainly loops and simple error checking.
	//
	// The steps are:
	// - prepare the working directory
	// - map the task files to the working directory
	// - log the IP address
	// - set up the storage configuration
	// - validate input and output files
	// - download inputs
	// - run the steps (docker)
	// - upload the outputs

	var run helper
  var reader TaskReader
	var task *tes.Task
  var store storage.Storage
  conf := w.Config()

  writer, werr := w.EventWriter()
  if writer == nil {
    // There's no event writer, which means no way to communicate
    // that the task failed or log the error. Just give up.
    return
  }
  if werr != nil {
    // There was an error creating a writer. This is a system error,
    // but there's still a non-nil writer which might be useful for
    // logging the error or event the failed task state.
    // Set run.syserr so that the normal failure flow happens.
    run.syserr = werr
  }
  event := events.NewTaskWriter(taskID, 0, conf.Logger.Level, writer)

	event.Info("Version", version.LogFields()...)
	event.StartTime(time.Now())

	// Run the final logging/state steps in a deferred function
	// to ensure they always run, even if there's a missed error.
	defer func() {
		event.EndTime(time.Now())

		switch {
		case run.taskCanceled:
			// The task was canceled.
			event.Info("Canceled")
			event.State(tes.State_CANCELED)
		case run.execerr != nil:
			// One of the executors failed
			event.Error("Exec error", "error", run.execerr)
			event.State(tes.State_ERROR)
		case run.syserr != nil:
			// Something else failed
			// TODO should we do something special for run.err == context.Canceled?
			event.Error("System error", "error", run.syserr)
			event.State(tes.State_SYSTEM_ERROR)
		default:
			event.State(tes.State_COMPLETE)
		}
	}()

	// Recover from panics
	defer handlePanic(func(e error) {
		run.syserr = e
	})

  if run.ok() {
    reader, run.syserr = w.TaskReader()
  }

	ctx := pollForCancel(pctx, reader, conf.UpdateRate, func() {
		run.taskCanceled = true
	})
	run.ctx = ctx

  if run.ok() {
	  task, run.syserr = reader.Task()
  }

	// Configure a task-specific storage backend.
	// This provides download/upload for inputs/outputs.
  if run.ok() {
    store, run.syserr = w.Storage(task)
  }

	if run.ok() {
		event.State(tes.State_INITIALIZING)
	}

	// Grab the IP address of this host. Used to send task metadata updates.
	var ip string
	if run.ok() {
		ip, run.syserr = externalIP()
	}

	if run.ok() {
		run.syserr = validateInputs(task.Inputs, store)
	}

	if run.ok() {
		run.syserr = validateOutputs(task.Outputs, store)
	}

	// Download inputs
	for _, input := range task.Inputs {
		if run.ok() {
			event.Info("Starting download", "url", input.Url)
			err := store.Get(ctx, input.Url, input.Path, input.Type)
			if err != nil {
				run.syserr = err
				event.Error("Download failed", "url", input.Url, "error", err)
			} else {
				event.Info("Download finished", "url", input.Url)
			}
		}
	}

	if run.ok() {
		event.State(tes.State_RUNNING)
	}

	// Run steps
	for i, d := range task.Executors {
    var stdio *Stdio
		if run.ok() {
      stdio, run.syserr = w.Stdio(d)
    }

		if run.ok() {
      s := &stepWorker{
        Conf:  conf,
        Event: event.NewExecutorWriter(uint32(i)),
        IP:    ip,
        Exec:  w.Executor(task, i),
      }
			run.execerr = s.Run(ctx)
		}
	}

	// Upload outputs
	var outputs []*tes.OutputFileLog
	for _, output := range task.Outputs {
		if run.ok() {
			event.Info("Starting upload", "url", output.Url)
			out, err := store.Put(ctx, output.Url, output.Path, output.Type)
			if err != nil {
				run.syserr = err
				event.Error("Upload failed", "url", output.Url, "error", err)
			} else {
				event.Info("Upload finished", "url", output.Url)
			}
			outputs = append(outputs, out...)
		}
	}

	if run.ok() {
		event.Outputs(outputs)
	}
}

func openStdio() {
  stdio := Stdio{}
  var err error

  // Find the path for task stdin
  if ex.Stdin != "" {
    stdio.Stdin, err = mapper.OpenFile(ex.Stdin)
    if err != nil {
      return nil, fmt.Errorf("couldn't open stdin", err)
    }
  }

  // Create file for task stdout
  if ex.Stdout != "" {
    stdio.Stdout, err = mapper.CreateFile(ex.Stdout)
    if err != nil {
      return nil, fmt.Errorf("couldn't create stdout", err)
    }
  }

  // Create file for task stderr
  if ex.Stderr != "" {
    stdio.Stderr, err = mapper.CreateFile(ex.Stderr)
    if err != nil {
      return nil, fmt.Errorf("couldn't create stderr", err)
    }
  }
  return &stdio, nil
}

// OpenHostFile opens a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Open
//
// If the path can't be mapped or the file can't be opened, an error is returned.
func (mapper *FileMapper) OpenHostFile(src string) (*os.File, error) {
	p, perr := mapper.HostPath(src)
	if perr != nil {
		return nil, perr
	}
	f, oerr := os.Open(p)
	if oerr != nil {
		return nil, oerr
	}
	return f, nil
}

// CreateHostFile creates a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Create
//
// If the path can't be mapped or the file can't be created, an error is returned.
func (mapper *FileMapper) CreateHostFile(src string) (*os.File, error) {
	p, perr := mapper.HostPath(src)
	if perr != nil {
		return nil, perr
	}
	err := util.EnsurePath(p)
	if err != nil {
		return nil, err
	}
	f, oerr := os.Create(p)
	if oerr != nil {
		return nil, oerr
	}
	return f, nil
}

// Validate the input downloads
func validateInputs(inputs []*tes.TaskParameter, s storage.Storage) error {
	for _, input := range inputs {
		if !s.Supports(input.Url, input.Path, input.Type) {
			return fmt.Errorf("Input download not supported by storage: %v", input)
		}
	}
	return nil
}

// Validate the output uploads
func validateOutputs(outputs []*tes.TaskParameter, s storage.Storage) error {
	for _, output := range outputs {
		if !s.Supports(output.Url, output.Path, output.Type) {
			return fmt.Errorf("Output upload not supported by storage: %v", output)
		}
	}
	return nil
}

func pollForCancel(ctx context.Context, r TaskReader, rate time.Duration, f func()) context.Context {
	taskctx, cancel := context.WithCancel(ctx)

	// Start a goroutine that polls the server to watch for a canceled state.
	// If a cancel state is found, "taskctx" is canceled.
	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()

		for {
			select {
			case <-taskctx.Done():
				return
			case <-ticker.C:
				state, _ := r.State()
				if tes.TerminalState(state) {
					cancel()
					f()
				}
			}
		}
	}()
	return taskctx
}
