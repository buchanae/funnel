package scheduler

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"golang.org/x/net/context"
	"strings"
	"time"
)

// Database represents the interface to the database used by the scheduler, scaler, etc.
// Mostly, this exists so it can be mocked during testing.
type Database interface {
	ReadQueue(n int) []*tes.Task
	AssignTask(*tes.Task, *pbf.Worker)
	CheckWorkers() error
	ListWorkers(context.Context, *pbf.ListWorkersRequest) (*pbf.ListWorkersResponse, error)
	UpdateWorker(context.Context, *pbf.Worker) (*pbf.UpdateWorkerResponse, error)
}

// NewScheduler returns a new Scheduler instance.
func NewScheduler(db Database, conf config.Config) (*Scheduler, error) {
	backends := map[string]*BackendPlugin{}

	err := util.EnsureDir(conf.WorkDir)
	if err != nil {
		return nil, err
	}

	return &Scheduler{db, conf, backends}, nil
}

// Scheduler handles scheduling tasks to workers and support many backends.
type Scheduler struct {
  // TODO switch to TaskQueue interface
	db       Database
	conf     config.Config
	backends map[string]*BackendPlugin
}

// AddBackend adds a backend plugin.
func (s *Scheduler) AddBackend(plugin *BackendPlugin) {
	s.backends[plugin.Name] = plugin
}

// Start starts the scheduling loop. This blocks.
//
// The scheduler will take a chunk of tasks from the queue,
// request the the configured backend schedule them, and
// act on offers made by the backend.
func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.conf.ScheduleRate)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			var err error
			err = s.Schedule()
			if err != nil {
				log.Error("Schedule error", err)
				return err
			}
		}
	}
}

// Schedule does a scheduling iteration. It checks the health of workers
// in the database, gets a chunk of tasks from the queue (configurable by config.ScheduleChunk),
// and calls the given scheduler. If the scheduler returns a valid offer, the
// task is assigned to the offered worker.
func (s *Scheduler) Schedule() error {
  // TODO try to get rid of this.
	backend, err := s.backend()
	if err != nil {
		return err
	}

  // TODO move this elsewhere
	s.db.CheckWorkers()

  // TODO extract this into something that streams tasks into this schedule loop
	for _, task := range s.db.ReadQueue(s.conf.ScheduleChunk) {
		offer := backend.Schedule(task)

    // In the future, a more advanced Funnel might make decisions between
    // multiple offers here. For now, we just accept the first and only offer.

		if offer != nil {
			log.Info("Assigning task to worker",
				"taskID", task.Id,
				"workerID", offer.Worker.Id,
			)

      // offer.OnAccept allows the scheduler backend to act on the offer being
      // accepted, e.g. to provision a worker in HTCondor.
      if offer.OnAccept != nil {
        if err := offer.OnAccept(); err != nil {
          // The OnAccept callback failed, so consider the task scheduling failed.
          // Break so that the task isn't marked as assigned.
          log.Error("Task was scheduled, but OnAccept failed")
          break
        }
      }
			s.db.AssignTask(task.Id, offer.Worker.Id)
		} else {
			log.Debug("No worker could be scheduled for task", "taskID", task.Id)
		}
	}
	return nil
}

// backend returns a Backend instance for the backend
// given by name in config.Scheduler.
func (s *Scheduler) backend() (Backend, error) {
	name := strings.ToLower(s.conf.Scheduler)
	plugin, ok := s.backends[name]

	if !ok {
		log.Error("Unknown scheduler backend", "name", name)
		return nil, fmt.Errorf("Unknown scheduler backend %s", name)
	}

	// Cache the scheduler instance on the plugin so that
	// we can call this backend() function repeatedly.
	if plugin.instance == nil {
		i, err := plugin.Create(s.conf)
		if err != nil {
			return nil, err
		}
		plugin.instance = i
	}
	return plugin.instance, nil
}
