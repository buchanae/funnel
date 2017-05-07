package scheduler

import (
	"github.com/ohsu-comp-bio/funnel/config"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Backend is responsible for scheduling a task. It has a single method which
// is responsible for taking a Task and returning an Offer, or nil if there is
// no worker matching the task request. An Offer includes the ID of the offered
// worker.
//
// Offers include scores which describe how well the task fits the worker.
// Scores may describe a wide variety of metrics: resource usage, packing,
// startup time, cost, etc. Scores and weights are used to control the behavior
// of schedulers, and to combine offers from multiple schedulers.
type Backend interface {
	Schedule(*tes.Task) Offer
}

// Offer describes a worker offered by a scheduler for a task.
// The Scores describe how well the task fits this worker,
// which could be used by other a scheduler to pick the best offer.
type Offer struct {
	Worker *pbf.Worker
	Scores Scores
  OnAccept func() error
}

// BackendPlugin is provided by backends when they register with Scheduler,
// which allows to the scheduler to create a backend instance by name.
type BackendPlugin struct {
	Name     string
	Create   func(config.Config) (Backend, error)
	instance Backend
}
