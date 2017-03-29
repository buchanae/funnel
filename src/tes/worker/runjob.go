package worker

import (
	"context"
	"io"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	pbr "tes/server/proto"
	"time"
  "tes/util/ticker"
)


func (w *Worker) runJob(ctrl JobControl, wpr *pbr.JobWrapper) {
	job := wpr.Job
  log := w.log.WithFields("jobID", job.JobID)
	log.Debug("Running job")
  var store storage.Storage

  // Pipeline helps connect multiple sequential steps
  // while stopping on the first error. All steps here
  // should be wrapped in p.Run() so that errors and
  // cancelation are handeled correctly.
  p := pipeline.WithContext(ctrl.Context())

  p.Run(func() error {
    var err error
    store, err = w.backend.Storage(job)
    return err
  })

  // Validate job
  p.Run(func() error {
    return w.backend.Validate(job)
  })

  // Validate inputs against storage
  // TODO would be better if this logged a reason
	for _, input := range job.Task.Inputs {
    p.Run(func() error {
      if !store.Supports(input.Location, input.Path, input.Class) {
        return fmt.Errorf("Input download not supported by storage: %v", input)
      }
    }
  }

  // Validate outputs against storage
	for _, output := range job.Task.Outputs {
    p.Run(func() error {
      if !store.Supports(output.Location, output.Path, output.Class) {
        return fmt.Errorf("Output upload not supported by storage: %v", output)
      }
    })
  }

	// Download inputs
	for _, input := range job.Task.Inputs {
    p.Run(func() error {
      return store.Get(ctx, input.Location, input.Path, input.Class, true)
    })
  }

  p.Run(func() error {
    ctrl.SetRunning()
  })

	// Run executors
	for i, _ := range job.Task.Docker {
    p.Run(func() error {
      // Use context to manage executor goroutines
      subctx, cleanup := context.WithCancel(ctrl.Context())
      defer cleanup()
      // Start goroutines to watch executor for metadata and logs
      go w.inspectExecutor(subctx, job, i)
      go w.watchExecutorLogs(subctx, job, i)
      // Start executor
      res := w.backend.Execute(subctx, job, i)
      // Send update for exit code
      w.update(jobID, i, &pbe.JobLog{
        ExitCode: getExitCode(res),
      })
      return res
    })
	}

  // Upload outputs
	for _, output := range job.Task.Outputs {
    p.Run(func() error {
      return store.Get(ctx, output.Location, output.Path, output.Class)
    })
  }

  ctrl.SetResult(p.Err)
}

// Inspect container metadata when it's available,
// e.g. get the ports mapped by docker.
func (w *Worker) inspectExecutor(ctx, job, i) {
  meta := w.backend.Inspect(ctx, job, i)
  w.update(jobID, i, &pbe.JobLog{
    HostIP: meta.IP,
    Ports: meta.Ports,
  })
}

// Tail the stdout/err logs and send updates back to the server
func (w *Worker) watchExecutorLogs(ctx context.Context) {
  // Send updates to server on every tick.
  // Stop when context is canceled.
  t := ticker.WithContext(ctx, conf.LogUpdateRate)
  defer t.Stop()
  for _ := range t.C {
    w.flushLogs(logs)
  }
  // Ensure one last flush on exit
  w.flushLogs(logs)
}

// update sends an update of the JobLog for the given job executor.
// Used to update stdout/err logs, port mapping, etc.
func (w *Worker) update(jobID string, i int, log *pbe.JobLog) {
  up := &pbr.UpdateJobLogsRequest{
		Id:   j.job.JobID,
		Step: int64(stepID),
		Log:  log,
	}
	// UpdateJobLogs() is more lightweight than UpdateWorker(),
	// which is why it happens separately and at a different rate.
	err := w.sched.UpdateJobLogs(up)
	if err != nil {
		// TODO if the request failed, the job update is lost and the logs
		//      are corrupted. Cache logs to prevent this?
    //      At least return error from update?
		w.log.Error("Job log update failed", err)
	}
}

func (w *Worker) flushLogs(l stepLogs) {
  // TODO don't flush until the update has been acknowledged
  //      in order to protect against dropped logs
  //      i.e. when the server can't be contacted, logs
  //      should stay buffered
  w.update(jobID, stepID, *pbe.JobLog{
    Stdout: l.Stdout.Flush(),
    Stderr: l.Stderr.Flush(),
  })
}

func logTails(size int64) (*tailer, *tailer) {
	stdout, _ := newTailer(size)
	stderr, _ := newTailer(size)

	if s.Cmd.Stdout != nil {
		s.Cmd.Stdout = io.MultiWriter(s.Cmd.Stdout, stdout)
	}
	if s.Cmd.Stderr != nil {
		s.Cmd.Stderr = io.MultiWriter(s.Cmd.Stderr, stderr)
	}
	return stdout, stderr
}
