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

type Executor interface {
  Run(context.Context) error
  Inspect(context.Context) []*tes.Ports
}

func (r *taskRunner) Stdin(i int) io.Reader {
  d := task.Executors[i]
	if d.Stdin != "" {
    // Ignoring the error because it is expected to have been checked
    // during initialization.
    // TODO maybe save file handles in an array
    f, _ := r.mapper.OpenHostFile(d.Stdin)
    return r
	}
  return nil
}

func (r *taskRunner) Stdout(i int) io.Writer {
  d := task.Executors[i]
	// Create file for task stdout
	if d.Stdout != "" {
    // Ignoring the error because it is expected to have been checked
    // during initialization.
    // TODO maybe save file handles in an array
    f, _ := r.mapper.CreateHostFile(d.Stdout)
    f = io.MultiWriter(f, r.TaskLogger.Stdout(i))
    return f
	}
  return nil
}

func (r *taskRunner) Stderr(i int) io.Writer {
  d := task.Executors[i]
	// Create file for task stderr
	if d.Stderr != "" {
    // Ignoring the error because it is expected to have been checked
    // during initialization.
    // TODO maybe save file handles in an array
    f, _ := r.mapper.CreateHostFile(d.Stderr)
    f = io.MultiWriter(f, r.TaskLogger.Stderr(i))
    return f
	}
  return nil
}
