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
	"github.com/ohsu-comp-bio/funnel/scheduler"
)

var log = logger.New("gce")

// Plugin provides the Google Cloud Compute scheduler backend plugin.
var Plugin = &scheduler.BackendPlugin{
	Name:   "gce",
	Create: NewBackend,
}

// NewBackend returns a new Google Cloud Engine Backend instance.
func NewBackend(conf config.Config) (scheduler.Backend, error) {
	// TODO need GCE scheduler config validation. If zone is missing, nothing works.

	// Create a client for talking to the funnel scheduler
  // TODO switch to WorkerDB interface. Stick with simple in-memory for now.
	client, err := scheduler.NewClient(conf.Worker)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	// Create a client for talking to the GCE API
	gce, gerr := newClientFromConfig(conf)
	if gerr != nil {
		log.Error("Can't connect GCE client", gerr)
		return nil, gerr
	}

	s := &Backend{
		conf:   conf,
		client: client,
		gce:    gce,
	}

	return scheduler.Backend(s), nil
}

// Backend represents the GCE backend, which provides
// and interface for both scheduling and scaling.
type Backend struct {
	conf   config.Config
	client scheduler.Client
	gce    Client
}

// Schedule schedules a task on a Google Cloud VM worker instance.
func (s *Backend) Schedule(j *tes.Task) scheduler.Offer {
	log.Debug("Running GCE scheduler")

	offers := []scheduler.Offer{}
	predicates := append(scheduler.DefaultPredicates, scheduler.WorkerHasTag("gce"))

	for _, w := range s.getWorkers() {
		// Filter out workers that don't match the task request.
		// Checks CPU, RAM, disk space, ports, etc.
		if !scheduler.Match(w, j, predicates) {
			continue
		}
		offers = append(offers, makeOffer(w, j))
	}

	// No matching workers were found.
	if len(offers) == 0 {
		return nil
	}

	scheduler.SortByAverageScore(offers)
	return offers[0]
}

func (s *Backend) makeOffer(w *pbr.Worker, t *tes.Task) scheduler.Offer {

  weights := map[string]float32{}
  // TODO add useful weights:
  //      - prefer workers that are already online.

  return scheduler.Offer{
    Worker: w,
    Scores: scheduler.DefaultScores(w, t).Weighted(weights),
    OnAccept: func() error {
      // TODO check if this worker needs to be started

      // Get the template ID from the worker metadata
      template, ok := w.Metadata["gce-template"]
      if !ok || template == "" {
        return fmt.Errorf("Could not get GCE template ID from metadata")
      }

      // StartWorker calls out to GCE APIs to start a new worker instance.
      return s.gce.StartWorker(template, s.conf.RPCAddress(), w.Id)
    },
  }
}

// getWorkers returns a list of all GCE workers and appends a set of
// uninitialized workers, which the scheduler can use to create new worker VMs.
func (s *Backend) getWorkers() []*pbf.Worker {

	// Get the workers from the funnel server
	workers := []*pbf.Worker{}
	req := &pbf.ListWorkersRequest{}
	resp, err := s.client.ListWorkers(context.Background(), req)

	// If there's an error, return an empty list
	if err != nil {
		log.Error("Failed ListWorkers request. Recovering.", err)
		return workers
	}

	workers = resp.Workers

	// Include unprovisioned (template) workers.
	// This is how the scheduler can schedule tasks to workers that
	// haven't been started yet.
	for _, t := range s.gce.Templates() {
		t.Id = scheduler.GenWorkerID("funnel")
		workers = append(workers, &t)
	}

	return workers
}
