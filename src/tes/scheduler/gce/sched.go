package gce

// TODO
// - how to re-evaluate the resource pool after a worker is created (autoscale)?
// - resource tracking via GCP APIs
// - matching requirements to existing VMs
// - provisioning limits, e.g. don't create more than 100 VMs, or
//   maybe use N VCPUs max, across all VMs
// - act on failed machines?
// - know how to shutdown machines

// TODO outside of scheduler for SMC RNA
// - client that organizes SMCRNA tasks and submits
// - dashboard for tracking jobs
// - resource tracking via TES worker stats collection

import (
	"context"
	"fmt"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	sched "tes/scheduler"
	pbr "tes/server/proto"
)

var log = logger.New("gce")

// NewScheduler returns a new Google Cloud Engine Scheduler instance.
func NewScheduler(conf config.Config) (sched.Scheduler, error) {
	client, err := sched.NewClient(conf.Worker)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	gce, gerr := newGCEClient(context.TODO(), conf)
	if gerr != nil {
		log.Error("Can't connect GCE client", gerr)
		return nil, gerr
	}

	s := &gceScheduler{
		conf:   conf,
		client: client,
		gce:    gce,
	}

	return s, nil
}

type gceScheduler struct {
	conf   config.Config
	client sched.Client
	gce    GCEClient
}

// Schedule schedules a job on a Google Cloud VM worker instance.
func (s *gceScheduler) Schedule(j *pbe.Job) *sched.Offer {
	log.Debug("Running gce scheduler")

	offers := []*sched.Offer{}

	for _, w := range s.getWorkers() {
		// Filter out workers that don't match the job request.
		// Checks CPU, RAM, disk space, ports, etc.
		if !sched.Match(w, j, sched.DefaultPredicates) {
			continue
		}

		sc := sched.DefaultScores(w, j)
		/*
			    TODO
			    if w.State == pbr.WorkerState_Alive {
					  sc["startup time"] = 1.0
			    }
		*/
		sc = sc.Weighted(s.conf.Schedulers.GCE.Weights)

		offer := sched.NewOffer(w, j, sc)
		offers = append(offers, offer)
	}

	// No matching workers were found.
	if len(offers) == 0 {
		return nil
	}

	sched.SortByAverageScore(offers)
	return offers[0]
}

// getWorkers returns a list of all GCE workers which are not dead/gone.
// Also appends extra entries for unprovisioned workers.
func (s *gceScheduler) getWorkers() []*pbr.Worker {

	req := &pbr.GetWorkersRequest{}
	resp, err := s.client.GetWorkers(context.Background(), req)
	workers := []*pbr.Worker{}

	if err != nil {
		log.Error("Failed GetWorkers request. Recovering.", err)
		return workers
	}

	// Find all workers with GCE prefix in ID, that are not Dead/Gone.
	for _, w := range resp.Workers {
		// Only include workers with a "gce" key in their metadata
		_, isGce := w.Metadata["gce"]

		if isGce && w.State != pbr.WorkerState_Dead && w.State != pbr.WorkerState_Gone {
			workers = append(workers, w)
		}
	}

	project := s.conf.Schedulers.GCE.Project
	zone := s.conf.Schedulers.GCE.Zone

	// Include unprovisioned (template) workers.
	for _, t := range s.conf.Schedulers.GCE.Templates {
		res, err := s.gce.Template(project, t)

		if err != nil {
			log.Error("Couldn't get template from GCE. Skipping.",
				"error", err,
				"template", t)
			continue
		}
		// Copy resources for available resources
		avail := *res

		workers = append(workers, &pbr.Worker{
			Id:        sched.GenWorkerID(),
			Resources: res,
			Available: &avail,
			Zone:      zone,
			Metadata: map[string]string{
				"gce":          "yes",
				"gce-template": t,
			},
		})
	}

	return workers
}

// ShouldStartWorker tells the scaler loop which workers
// belong to this scheduler backend, basically.
func (s *gceScheduler) ShouldStartWorker(w *pbr.Worker) bool {
	tpl, ok := w.Metadata["gce-template"]
	return ok && tpl != "" && w.State == pbr.WorkerState_Uninitialized
}

// StartWorker calls out to GCE APIs to start a new worker instance.
func (s *gceScheduler) StartWorker(w *pbr.Worker) error {

	// Write the funnel worker config yaml to a string
	c := s.conf.Worker
	c.ID = w.Id
	// TODO this should be set as a worker default config somewhere else
	c.ServerAddress = s.conf.ServerAddress
	c.Timeout = -1
	c.Storage = s.conf.Storage

	project := s.conf.Schedulers.GCE.Project
	zone := s.conf.Schedulers.GCE.Zone

	// Get the template ID from the worker metadata
	template, ok := w.Metadata["gce-template"]
	if !ok || template == "" {
		return fmt.Errorf("Could not get GCE template ID from metadata")
	}

	return s.gce.StartWorker(project, zone, template, c)
}