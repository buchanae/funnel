package worker

import (
	"context"
	"errors"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"time"
)

// NewWorker returns a new Worker instance
func NewWorker(conf config.Worker) (*Worker, error) {
	sched, err := newSchedClient(conf)
	if err != nil {
		return nil, err
	}

	log := logger.New("worker", "workerID", conf.ID)
	log.Debug("Worker Config", "config.Worker", conf)

	return &Worker{
    conf: conf,
    sched: sched,
    log: log,
    resources: detectResources(conf.Resources),
    newRunner: defaultRunnerFactory,
    runners: map[string]Runner{},
    timeout: util.NewIdleTimeout(conf.Timeout),
    stop: make(chan struct{}),
	}, nil
}

// Worker is a worker...
type Worker struct {
	conf config.Worker
	sched      *schedClient
	log        logger.Logger
	resources  *pbf.Resources
  runners    map[string]context.CancelFunc
	timeout    util.IdleTimeout
	stop       chan struct{}
}

// Run runs a worker with the given config. This is responsible for communication
// with the server and starting task runners
func (w *Worker) Run() {
	w.log.Info("Starting worker")
	w.checkConnection()

	ticker := time.NewTicker(w.conf.UpdateRate)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.Sync()
		case <-w.timeout.Done():
			// Worker timeout reached. Shutdown.
			w.Stop()
			return
		}
	}
}

// Stop stops the worker
func (w *Worker) Stop() {
	w.timeout.Stop()
	for _, runner := range w.Runners {
    runner.Stop()
	}
  w.sched.Goodbye()
	w.sched.Close()
	close(w.stop)
}

func (w *Worker) checkConnection() {
	_, err := w.sched.GetWorker(context.TODO(), &pbf.GetWorkerRequest{Id: w.conf.ID})

	if err != nil {
		log.Error("Couldn't contact server.", err)
	} else {
		log.Info("Successfully connected to server.")
	}
}

// TODO Sync should probably use a channel to sync data access.
//      Probably only a problem for test code, where Sync is called directly.
func (w *Worker) Sync() {
	r, gerr := w.sched.GetWorker(context.TODO(), &pbf.GetWorkerRequest{Id: w.conf.ID})

	if gerr != nil {
		log.Error("Couldn't get worker state during sync.", gerr)
		return
	}

  if r.Shutdown {
    w.Stop()
    return
  }

  // TODO
  // Create runners after the update has saved to the database?
  // Is there anything to update anymore though?
  // Just a worker ping?
  for _, id := range r.TaskIds {
    w.ensureRunner(id)
  }

  // Check if the worker is idle. If so, start the timeout timer.
	if len(w.Runners) == 0 {
		w.timeout.Start()
	} else {
		w.timeout.Stop()
	}

	// Worker data has been updated. Send back to server for database update.
	r.Resources = w.resources

	_, err := w.sched.UpdateWorker(r)
	if err != nil {
		log.Error("Couldn't save worker update. Recovering.", err)
	}
}

func (d *Worker) ensureRunner(id string) {
  // Ensure a runner exists for the task
  if _, ok := w.Runners[id]; !ok {
    ctx, cancel := context.WithCancel(context.TODO())
    RunTask(ctx, id)
    w.Runners[id] = cancel
  }
}
