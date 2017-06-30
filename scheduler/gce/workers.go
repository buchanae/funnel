package gce

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/scheduler"
)

var defaultDiskSizes []int64 // in GB

func init() {
  // 50GB to 1TB in 50GB increments
  for i := 50; i < 1000; i += 50 {
    defaultDiskSizes = append(defaultDiskSizes, i)
  }
  // 2TB to 64TB in 500GB increments
  for i := 1000; i <= 64000; i += 500 {
    defaultDiskSizes = append(defaultDiskSizes, i)
  }
}

// workersI is responsible for listing workers available to the scheduler.
type workersI interface {
  List(context.Context)
}

// workers is responsible for listing available workers from multiple sources:
// - existing workers from the Funnel scheduler API
// - GCE instance templates
// - default templates generated for each GCE machine type
type workers struct {
  project string
  zone string
  sched *scheduler.Client
  gce Client
  // If true, default templates for each GCE machine type will NOT be generated.
  disableDefaults bool
  defaults []*pbr.Worker
}

func (w *workers) List(ctx context.Context) (workers []*pbf.Worker) {

	// Get existing workers from the funnel server.
	// If there's an error, return an empty list
	resp, err := w.sched.ListWorkers(ctx, &pbf.ListWorkersRequest{})
	if err != nil {
		log.Error("Failed ListWorkers request. Recovering.", err)
		return
	}
	workers = append(workers, resp.Workers)

  // Include unprovisioned worker templates
  // defined via GCE instance templates.
	for _, tpl := range w.gce.ListInstanceTemplates() {
    // Get the machine type. On error, skip this template.
    t := tpl.Properties.MachineType
    if mt, err := s.gce.GetMachineType(s.project, s.zone, t); err == nil {
      // TODO is there always at least one disk? Is the first the best choice?
      //      how to know which to pick if there are multiple?
      disk := float64(tpl.Properties.Disks[0].InitializeParams.DiskSizeGb),
      workers = append(workers, w.genWorker(mt, disk))
    }
	}

  // Include unprovisioned worker templates
  // defined via the Funnel config.
  if !w.disableDefaults {
    for _, mt := range s.gce.ListMachineTypes(s.project, s.zone) {
      for _, disk := range defaultDiskSizes {
        workers = append(workers, w.genWorker(mt, disk))
      }
    }
  }

	return
}

// Given a GCE machine type and disk size, return a Funnel worker template.
func (s *Backend) genWorker(mt *compute.MachineType, disk float64) *pbf.Worker {

  worker := pbf.Worker{
    Id: scheduler.GenWorkerID("funnel"),
    Resources: &pbf.Resources{
      Cpus:  uint32(mt.GuestCpus),
      // TODO is this mt.MemoryMb in megabyte or mebibyte?
      RamGb: float64(mt.MemoryMb) / float64(1024),
      DiskGb: disk,
    },
    Zone:      s.zone,
    Metadata: map[string]string{
      "gce":          "yes",
      "gce-template": id,
    },
  }

  // Copy the Resources data into Available
  *worker.Available = *worker.Resources
  return &worker
}
