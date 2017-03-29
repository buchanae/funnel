package worker

import (
	"context"
	"errors"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	pbr "tes/server/proto"
	"tes/util"
	"time"
)

// NewWorker returns a new Worker instance
func NewWorker(conf config.Worker) (*Worker, error) {
	sched, err := newSchedClient(conf)
	if err != nil {
		return nil, err
	}

	log := logger.New("worker", "workerID", conf.ID)
	res := detectResources(conf.Resources)
	// Tracks active job ctrls: job ID -> JobControl instance
	ctrls := map[string]JobControl{}
	timeout := util.NewIdleTimeout(conf.Timeout)
	stop := make(chan struct{})
	state := pbr.WorkerState_Uninitialized
  backend := TODO
	return &Worker{
		conf, sched, log, res, backend, ctrls, timeout, stop, state,
	}, nil
}

// Worker is a worker...
type Worker struct {
	conf config.Worker
	sched      *schedClient
	log        logger.Logger
	resources  *pbr.Resources
  backend    Backend
	Ctrls      map[string]JobControl
	timeout    util.IdleTimeout
	stop       chan struct{}
	state      pbr.WorkerState
}

// Run runs a worker with the given config. This is responsible for communication
// with the server and starting job runners
func (w *Worker) Run() {
	w.log.Info("Starting worker")
	w.state = pbr.WorkerState_Alive

	ticker := time.NewTicker(w.conf.UpdateRate)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.Sync()
			w.checkIdleTimer()
		case <-w.timeout.Done():
			// Worker timeout reached. Shutdown.
			w.Stop()
			return
		}
	}
}

// Stop stops the worker
// TODO need a way to shut the worker down from the server/scheduler.
func (w *Worker) Stop() {
	w.state = pbr.WorkerState_Gone
	close(w.stop)
	w.timeout.Stop()
	for _, ctrl := range w.Ctrls {
		ctrl.Cancel()
	}
	w.Sync()
	w.sched.Close()
}

// Check if the worker is idle. If so, start the timeout timer.
func (w *Worker) checkIdleTimer() {
	// The worker is idle if there are no job controllers.
	idle := len(w.Ctrls) == 0 && w.state == pbr.WorkerState_Alive
	if idle {
		w.timeout.Start()
	} else {
		w.timeout.Stop()
	}
}
