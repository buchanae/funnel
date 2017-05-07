package server

// TODO put the boltdb implementation in a separate package
//      so that users can import pluggable backends

import (
	"errors"
	"github.com/boltdb/bolt"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"time"
)

// UpdateWorker is an RPC endpoint that is used by workers to send heartbeats
// and status updates, such as completed tasks. The server responds with updated
// information for the worker, such as canceled tasks.
func (tb *TaskBolt) SyncWorker(ctx context.Context, req *pbf.Worker) (*pbf.Worker, error) {

	err := tb.db.Update(func(tx *bolt.Tx) error {
    // Get worker
    worker := getWorker(tx, req.Id)

    if worker.Version != 0 && req.Version != 0 && worker.Version != req.Version {
      return errors.New("Version outdated")
    }

    // TODO WorkerPings bucket
    worker.LastPing = time.Now().Unix()

    if req.Resources != nil {
      if worker.Resources == nil {
        worker.Resources = &pbf.Resources{}
      }
      // Merge resources
      if req.Resources.Cpus > 0 {
        worker.Resources.Cpus = req.Resources.Cpus
      }
      if req.Resources.RamGb > 0 {
        worker.Resources.RamGb = req.Resources.RamGb
      }
      if req.Resources.DiskGb > 0 {
        worker.Resources.DiskGb = req.Resources.DiskGb
      }
    }

    for k, v := range req.Metadata {
      worker.Metadata[k] = v
    }

    worker.Version = time.Now().Unix()
    putWorker(tx, worker)
	})
	return resp, err
}

// AssignTask assigns a task to a worker. This updates the task state to Initializing,
// and updates the worker (calls UpdateWorker()).
// TODO would be awesome if worker could call ListTasks(tag="funnel-worker-id=1234")
//      to get assigned tasks. Or an internal SearchTaskQueue(tag="funnel-worker-id=1234")
func (tb *TaskBolt) AssignTask(t *tes.Task, w *pbf.Worker) error {
	return tb.db.Update(func(tx *bolt.Tx) error {

		// TODO this is important! write a test for this line.
		//      when a task is assigned, its state is immediately Initializing
		//      even before the worker has received it.
    // TODO maybe don't change state. Maybe workers need to be able to show
    //      update and run tasks without going through the scheduler?
		transitionTaskState(tx, t.Id, tes.State_INITIALIZING)
		taskIDBytes := []byte(t.Id)
		workerIDBytes := []byte(w.Id)

		// TODO the database needs tests for this stuff. Getting errors during dev
		//      because it's easy to forget to link everything.
    // TODO nested bucket instead of concat key?
		key := append(workerIDBytes, taskIDBytes...)
		tx.Bucket(WorkerTasks).Put(key, taskIDBytes)
		tx.Bucket(TaskWorker).Put(taskIDBytes, workerIDBytes)

		if err != nil {
			return err
		}
		return nil
	})
}

// TODO include active ports. maybe move Available out of the protobuf message
//      and expect this helper to be used?
func updateAvailableResources(tx *bolt.Tx, worker *pbf.Worker) {
	// Calculate available resources
	a := pbf.Resources{
		Cpus:   worker.GetResources().GetCpus(),
		RamGb:  worker.GetResources().GetRamGb(),
		DiskGb: worker.GetResources().GetDiskGb(),
	}
	for taskID := range worker.Tasks {
		t := getTask(tx, taskID)
		res := t.GetResources()

		// Cpus are represented by an unsigned int, and if we blindly
		// subtract it will rollover to a very large number. So check first.
		rcpus := res.GetCpuCores()
		if rcpus >= a.Cpus {
			a.Cpus = 0
		} else {
			a.Cpus -= rcpus
		}

		a.RamGb -= res.GetRamGb()

		if a.Cpus < 0 {
			a.Cpus = 0
		}
		if a.RamGb < 0.0 {
			a.RamGb = 0.0
		}
	}
	worker.Available = &a
}

// GetWorker gets a worker
func (tb *TaskBolt) GetWorker(ctx context.Context, req *pbf.GetWorkerRequest) (*pbf.Worker, error) {
	var worker *pbf.Worker
	err := tb.db.View(func(tx *bolt.Tx) error {
		worker = getWorker(tx, req.Id)
		return nil
	})
	return worker, err
}

// CheckWorkers is used by the scheduler to check for and delete dead workers.
// This is not an RPC endpoint.
func (tb *TaskBolt) CheckWorkers() error {
  return tb.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(WorkerPings)
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
      workerID := k
      lastPing := time.Unix(v, 0)

			if lastPing == 0 {
				// This shouldn't be happening, because workers should be
				// created with LastPing, but give it the benefit of the doubt
				// and leave it alone.
				continue
			}

      if time.Since(lastPing) > tb.conf.WorkerPingTimeout {
        // TODO delete the worker record
			}
		}
		return nil
	})
}

// ListWorkers is an API endpoint that returns a list of workers.
func (tb *TaskBolt) ListWorkers(ctx context.Context, req *pbf.ListWorkersRequest) (*pbf.ListWorkersResponse, error) {
	resp := &pbf.ListWorkersResponse{}
	resp.Workers = []*pbf.Worker{}

	err := tb.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(Workers)
		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			worker := getWorker(tx, string(k))
			resp.Workers = append(resp.Workers, worker)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Look for an auth token related to the given task ID.
func getTaskAuth(tx *bolt.Tx, taskID string) string {
	idBytes := []byte(taskID)
	var auth string
	data := tx.Bucket(TaskAuthBucket).Get(idBytes)
	if data != nil {
		auth = string(data)
	}
	return auth
}
