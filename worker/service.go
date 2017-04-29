package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/util"
	"golang.org/x/sync/syncmap"
	"time"
)

func NewService(conf config.Worker) (*Service, error) {
	sched, err := newSchedClient(conf)
	if err != nil {
		return nil, err
	}

	return &Service{
		conf:      conf,
		sched:     sched,
		log:       logger.New("worker", "workerID", conf.ID),
		resources: detectResources(conf.Resources),
		runners:   runnerSet{},
		timeout:   util.NewIdleTimeout(conf.Timeout),
		stopped:   make(chan struct{}),
		clear:     make(chan string),
	}, nil
}

type Service struct {
	conf      config.Worker
	sched     *schedClient
	log       logger.Logger
	resources *pbf.Resources
	runners   runnerSet
	timeout   util.IdleTimeout
	stopped   chan struct{}
	stop      context.CancelFunc
}

func (w *Service) Run() {
	// Don't start if already running
	if w.stop != nil {
		return
	}

	w.log.Info("Starting worker service")
	w.checkConnection()
	ctx, cancel := context.WithCancel(context.Background())
	w.stop = cancel

	ticker := time.NewTicker(w.conf.UpdateRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.sync(ctx)
		case <-w.timeout.Done():
			// Service timeout reached. Shutdown.
			w.Stop()
		case <-ctx.Done():
			// Clean up
			w.timeout.Stop()
			w.sched.Goodbye()
			w.sched.Close()
			w.runners.Wait()
			close(w.stopped)
			return
		}
	}
}

// Stop stops the service
func (w *Service) Stop() {
	// Don't stop if not running
	if w.stop == nil {
		return
	}
	w.stop()
	<-w.stopped
}

func (w *Service) checkConnection(ctx context.Context) {
	_, err := w.sched.GetWorker(ctx, &pbf.GetWorkerRequest{Id: w.conf.ID})

	if err != nil {
		log.Error("Couldn't contact server.", err)
		// TODO what should happen if the service can't connect?
	} else {
		log.Info("Successfully connected to server.")
	}
}

// TODO Sync should probably use a channel to sync data access.
//      Probably only a problem for test code, where Sync is called directly.
func (w *Service) sync(ctx context.Context) {
	r, gerr := w.sched.GetWorker(ctx, &pbf.GetWorkerRequest{Id: w.conf.ID})

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
		w.runners.Run(ctx, w.conf, id)
	}

	// Check if the worker is idle. If so, start the timeout timer.
	if w.runners.Count() == 0 {
		w.timeout.Start()
	} else {
		w.timeout.Stop()
	}

	// Worker data has been updated. Send back to server for database update.
	r.resources = w.resources

	_, err := w.sched.UpdateWorker(r)
	if err != nil {
		log.Error("Couldn't save worker update. Recovering.", err)
	}
}

// runnerSet manages concurrent access to a set of runners.
// runnerSet is safe for concurrent access.
type runnerSet struct {
	wg  sync.WaitGroup
	mtx sync.Mutex
	ids map[string]bool
}

// Run will call RunTask in a gouroutine and increment the waitgroup count.
// It also tracks task IDs, to ensure there's only one runner per task ID.
func (r *runnerSet) Run(ctx context.Context, c config.Worker, taskID string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	// If there's already a runner for the given task ID,
	// do nothing.
	if ok := r.ids[taskID]; ok {
		return
	}

	// Initialize map if needed
	if r.ids == nil {
		r.ids = make(map[string]bool)
	}

	r.ids[taskID] = true

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		RunTask(ctx, c, taskID)
		// Remove the task ID from the set
		r.mtx.Lock()
		defer r.mtx.Unlock()
		delete(r.ids, taskID)
	}()
}

// Wait for all runners to exit.
func (r *runnerSet) Wait() {
	r.wg.Wait()
}

// Count returns the number of runners currently running.
func (r *runnerSet) Count() int {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	return len(r.ids)
}
