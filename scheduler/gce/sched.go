package gce

// TODO
// - resource tracking via GCP APIs
//   - check project/region quotas
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

var log = logger.Sub("gce")

// Plugin provides the Google Cloud Compute scheduler backend plugin.
var Plugin = &scheduler.BackendPlugin{
	Name:   "gce",
	Create: NewBackend,
}

// NewBackend returns a new Google Cloud Engine Backend instance.
func NewBackend(conf config.Config) (scheduler.Backend, error) {
  if conf.Zone == "" || conf.Project == "" {
    return nil, fmt.Error("invalid GCE config: missing zone or project")
  }

	// Create a client for talking to the funnel scheduler
	client, err := scheduler.NewClient(conf.Worker)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	// Create a client for talking to the GCE API
	gce, gerr := newCachingClientFromConfig(conf)
	if gerr != nil {
		log.Error("Can't connect GCE client", gerr)
		return nil, gerr
	}

  return &Backend{
		gce:    gce,
    workers: &workers{
      project: conf.Backends.GCE.Project,
      zone: conf.Backends.GCE.Zone,
      sched: sched,
      client: gce,
      disableDefaults: conf.Backends.GCE.DisableDefaultTemplates,
    },
    project: conf.Backends.GCE.Project,
    zone: conf.Backends.GCE.Zone,
    serverAddress: conf.RPCAddress(),
	}, nil
}

// Backend represents the GCE backend, which provides
// and interface for both scheduling and scaling.
type Backend struct {
	gce    Client
  workers Workers
  project string
  zone string
  serverAddress string
}

// Schedule schedules a task on a Google Cloud VM worker instance.
func (s *Backend) Schedule(task *tes.Task) *scheduler.Offer {
	log.Debug("Running GCE scheduler")
  workers := s.workers.List(task)
  weights := map[string]float32{}
  return scheduler.DefaultScheduleAlgorithm(task, workers, weights)
}

// ShouldStartWorker tells the scaler loop which workers
// belong to this scheduler backend, basically.
func (s *Backend) ShouldStartWorker(w *pbf.Worker) bool {
	// Only start works that are uninitialized and have a gce template.
	tpl, ok := w.Metadata["gce-template"]
	return ok && tpl != "" && w.State == pbf.WorkerState_UNINITIALIZED
}

// StartWorker calls out to GCE APIs to start a new worker instance.
func (s *Backend) StartWorker(worker *pbf.Worker) error {

	// Get the template ID from the worker metadata
	tplID, ok := worker.Metadata["gce-template"]
	if !ok || tplID == "" {
		return fmt.Errorf("Could not get GCE template ID from metadata")
	}

	// Get the instance template from the GCE API
  tpl, ok := s.gce.Template(tplID)
	if !ok {
		return fmt.Errorf("Instance template not found: %s", tplName)
	}

	// Add GCE instance metadata
  serverAddress := s.serverAddress
	props := tpl.Properties

	// Create the instance on GCE
  instance := &compute.Instance{
		Name:              worker.Id,
		CanIpForward:      props.CanIpForward,
		Description:       props.Description,
		Disks:             props.Disks,
		MachineType:       props.MachineType,
		NetworkInterfaces: props.NetworkInterfaces,
		Scheduling:        props.Scheduling,
		ServiceAccounts:   props.ServiceAccounts,
		Tags:              props.Tags,
		Metadata:          &compute.Metadata{
      Items: append(props.Metadata.Items, &compute.MetadataItems{
        Key:   "funnel-worker-serveraddress",
        Value: &serverAddress,
      }),
    },
	}

  // Localize values to the zone
  instance.MachineType = localize(s.zone, "machineTypes", instance.MachineType)
	for _, disk := range instance.Disks {
		dt := localize(s.zone, "diskTypes", disk.InitializeParams.DiskType)
		disk.InitializeParams.DiskType = dt
	}

	op, ierr := s.gce.InsertInstance(s.project, s.zone, instance)
	if ierr != nil {
		log.Error("Couldn't insert GCE VM instance", ierr)
		return ierr
	}

	log.Debug("GCE VM instance created", "details", op)
  return nil
}

// localize helps make a resource string zone-specific
func localize(zone, resourceType, val string) string {
	return fmt.Sprintf("zones/%s/%s/%s", zone, resourceType, val)
}
