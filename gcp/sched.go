package gce

// TODO
// - resource tracking via GCP APIs
// - provisioning limits, e.g. don't create more than 100 VMs, or
//   maybe use N VCPUs max, across all VMs
// - act on failed machines?
// - know how to shutdown machines

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

const Name = "gcp"

var log = logger.Sub(Name)

// NewBackend returns a new Google Cloud Engine Backend instance.
func NewBackend(conf config.Config) (scheduler.Backend, error) {
	// TODO need GCE scheduler config validation. If zone is missing, nothing works.

	// Create a client for talking to the GCE API
	gce, gerr := newClientFromConfig(conf)
	if gerr != nil {
		log.Error("Can't connect GCE client", gerr)
		return nil, gerr
	}

	s := &Backend{
		conf:   conf,
		gce:    gce,
	}
}

// Backend represents the GCE backend, which provides
// and interface for both scheduling and scaling.
type Backend struct {
	conf   config.Config
	client scheduler.Client
	gce    Client
}

// Schedule schedules a task on a Google Cloud VM worker instance.
func (s *Backend) Schedule(j *tes.Task) *scheduler.Offer {
  w := pbf.Worker{
  }

  return scheduler.NewOffer(w, j, sc)
}

// getWorkers returns a list of all GCE workers and appends a set of
// uninitialized workers, which the scheduler can use to create new worker VMs.
func (s *Backend) getWorkers() []*pbf.Worker {

  // Include instance templates
	for _, t := range s.gce.Templates() {
		t.Id = scheduler.GenWorkerID("funnel")
		workers = append(workers, &t)
	}

	return workers
}
