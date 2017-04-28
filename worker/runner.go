package worker

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/util"
	"os"
	"path"
	"path/filepath"
  "time"
)

func RunTask(ctx context.Context, tw *pbr.TaskWrapper) {
  // Create a backend of the given name.
  // Backends are created per-task.
  backend, err := 
	task := tw.Task

	// The code here is verbose, but simple; mainly loops and simple error checking.
	//
	// The steps are:
	// 1. validate input and output mappings
	// 2. download inputs
	// 3. run the executors
	// 4. upload the outputs

  // Validate the input and outputs
  // TODO concat?
  params := append(task.Inputs[:], task.Outputs)
  for _, param := range params {
    r.Add(func() error {
      return backend.Supports(param.Url, param.Path, param.Type) {
    })
  })

  // TODO validate stdin/out/err?

	// Download inputs
	for _, input := range task.Inputs {
		r.Add(func() error {
			return backend.Get(ctx, input.Url, input.Path, input.Type)
		})
	}

  // Set task to running state
  r.Add(backend.Running)

	// Run executors
	for i, d := range task.Executors {
		r.Add(func() error {
      backend.Debug("Running executor", "i", i)
      backend.StartTime(i, time.Now().Format(time.RFC3339))

      // subctx ensures goroutines are cleaned up when the step exits.
      subctx, cleanup := context.WithCancel(ctx)
      defer cleanup()

      executor, err := backend.Executor(i)
      if err != nil {
        return err
      }

      // Run the executor
      done := make(chan error)
      go func() {
        done <- executor.Run(subctx)
      }()

      // Inspect the executor for metadata
      go func() {
        meta := s.Inspect(subctx)
        b.Ports(i, meta.ports)
        b.IP(i, meta.IP)
      }()

      // Wait for executor to exit
      res := <-done
      backend.EndTime(i, time.Now().Format(time.RFC3339))
      backend.ExitCode(i, getExitCode(res))
      return res
		})

	}

	// Upload outputs
	for _, output := range task.Outputs {
		r.Add(func() error {
      // TODO move to storage wrapper
			//r.fixLinks(output.Path)
			return backend.Put(ctx, output.Url, output.Path, output.Type)
		})
	}

  return r.Run()
}

type dolist []func() error
func (dl *dolist) Add(f func() error) {
  *dolist = append(dolist, f)
}
func (dl *dolist) Run() error {
  for _, f := range dl {
    err := f()
    if err != nil {
      return err
    }
  }
  return nil
}
