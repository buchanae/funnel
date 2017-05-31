package badger

// TODO put the boltdb implementation in a separate package
//      so that users can import pluggable backends

import (
  "bytes"
	"errors"
  "github.com/dgraph-io/badger/badger"
	proto "github.com/golang/protobuf/proto"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
	"time"
)

var log = logger.New("badger")

// State variables for convenience
const (
	Unknown      = tes.State_UNKNOWN
	Queued       = tes.State_QUEUED
	Running      = tes.State_RUNNING
	Paused       = tes.State_PAUSED
	Complete     = tes.State_COMPLETE
	Error        = tes.State_ERROR
	SystemError  = tes.State_SYSTEM_ERROR
	Canceled     = tes.State_CANCELED
	Initializing = tes.State_INITIALIZING
)

// UpdateWorker is an RPC endpoint that is used by workers to send heartbeats
// and status updates, such as completed tasks. The server responds with updated
// information for the worker, such as canceled tasks.
func (tb *TaskBadger) UpdateWorker(ctx context.Context, req *pbf.Worker) (*pbf.UpdateWorkerResponse, error) {
  err := tb.updateWorker(req)
  return &pbf.UpdateWorkerResponse{}, err
}

func (tb *TaskBadger) updateWorker(req *pbf.Worker) error {
  log.Debug("updateworker", req)
	worker, _ := tb.getWorker(req.Id)

	if worker.Version != 0 && req.Version != 0 && worker.Version != req.Version {
		return errors.New("Version outdated")
	}

	worker.LastPing = time.Now().Unix()
	worker.State = req.GetState()

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

	// Reconcile worker's task states with database
	for _, wrapper := range req.Tasks {
		// TODO test transition to self a noop
		task := wrapper.Task
		err := tb.transitionTaskState(task.Id, task.State)
		// TODO what's the proper behavior of an error?
		//      this is just ignoring the error, but it will happen again
		//      on the next update.
		//      need to resolve the conflicting states.
		//      Additionally, returning an error here will fail the db transaction,
		//      preventing all updates to this worker for all tasks.
		if err != nil {
			return err
		}

		// If the worker has acknowledged that the task is complete,
		// unlink the task from the worker.
		switch task.State {
		case Canceled, Complete, Error, SystemError:
      if _, ok := worker.Tasks[task.Id]; ok {
        delete(worker.Tasks, task.Id)
      }
		}
	}

	for k, v := range req.Metadata {
		worker.Metadata[k] = v
	}

	tb.updateAvailableResources(worker)
	worker.Version = time.Now().Unix()
	tb.putWorker(worker)
	return nil
}

// AssignTask assigns a task to a worker. This updates the task state to Initializing,
// and updates the worker (calls UpdateWorker()).
func (tb *TaskBadger) AssignTask(t *tes.Task, w *pbf.Worker) {
  // TODO this is important! write a test for this line.
  //      when a task is assigned, its state is immediately Initializing
  //      even before the worker has received it.
  tb.transitionTaskState(t.Id, tes.State_INITIALIZING)
  t.State = tes.State_INITIALIZING
  if w.Tasks == nil {
    w.Tasks = map[string]*pbf.TaskWrapper{}
  }
  w.Tasks[t.Id] = &pbf.TaskWrapper{Task: t}
  log.Debug("ASSIGN", w)
  tb.putWorker(w)
  tb.updateWorker(w)
}

// TODO include active ports. maybe move Available out of the protobuf message
//      and expect this helper to be used?
func (tb *TaskBadger) updateAvailableResources(worker *pbf.Worker) {
	// Calculate available resources
	a := pbf.Resources{
		Cpus:   worker.GetResources().GetCpus(),
		RamGb:  worker.GetResources().GetRamGb(),
		DiskGb: worker.GetResources().GetDiskGb(),
	}
	for taskID := range worker.Tasks {
    t, _ := tb.getTask(taskID)
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
func (tb *TaskBadger) GetWorker(ctx context.Context, req *pbf.GetWorkerRequest) (*pbf.Worker, error) {
  return tb.getWorker(req.Id)
}

// CheckWorkers is used by the scheduler to check for dead/gone workers.
// This is not an RPC endpoint
func (tb *TaskBadger) CheckWorkers() error {
  return nil
}

// ListWorkers is an API endpoint that returns a list of workers.
func (tb *TaskBadger) ListWorkers(ctx context.Context, req *pbf.ListWorkersRequest) (*pbf.ListWorkersResponse, error) {
  var workers []*pbf.Worker

  opt := badger.DefaultIteratorOptions
  itr := tb.db.NewIterator(opt)
  defer itr.Close()

  for itr.Seek(workerKey("")); itr.Valid(); itr.Next() {
    item := itr.Item()
    key := item.Key()
    if !bytes.HasPrefix(key, workerKey("")) {
      break
    }
    val := item.Value()
    var worker pbf.Worker
    proto.Unmarshal(val, &worker)
    // TODO hack hack hack
    if worker.Tasks == nil {
      worker.Tasks = map[string]*pbf.TaskWrapper{}
    }
    workers = append(workers, &worker)
  }

  return &pbf.ListWorkersResponse{
    Workers: workers,
	}, nil
}

func (tb *TaskBadger) transitionTaskState(id string, state tes.State) error {
  task, geterr := tb.getTask(id)
  if geterr != nil {
    return geterr
  }
  current := task.State

	switch current {
	case state:
		// Current state matches target state. Do nothing.
		return nil

	case Complete, Error, SystemError, Canceled:
		// Current state is a terminal state, can't do that.
		err := errors.New("Invalid state change")
		log.Error("Cannot change state of a task already in a terminal state",
			"error", err,
			"current", current,
			"requested", state)
		return err
	}

	switch state {
	case Canceled, Complete, Error, SystemError:
		// Remove from queue
    tb.db.Delete(queueKey(id))

	case Running, Initializing:
		if current != Unknown && current != Queued && current != Initializing {
			log.Error("Unexpected transition", "current", current, "requested", state)
			return errors.New("Unexpected transition to Initializing")
		}
    tb.db.Delete(queueKey(id))

	case Unknown, Paused:
		log.Error("Unimplemented task state", "state", state)
		return errors.New("Unimplemented task state")

	case Queued:
		log.Error("Can't transition to Queued state")
		return errors.New("Can't transition to Queued state")
	default:
		log.Error("Unknown task state", "state", state)
		return errors.New("Unknown task state")
	}

  task.State = state
  tb.putTask(task)
	log.Info("Set task state", "taskID", id, "state", state.String())
	return nil
}

// UpdateExecutorLogs is an API endpoint that updates the logs of a task.
// This is used by workers to communicate task updates to the server.
func (tb *TaskBadger) UpdateExecutorLogs(ctx context.Context, req *pbf.UpdateExecutorLogsRequest) (*pbf.UpdateExecutorLogsResponse, error) {

  task, _ := tb.getTask(req.Id)

  // max size (bytes) for stderr and stdout streams to keep in db
  max := tb.conf.MaxExecutorLogSize

  if req.Log != nil {
    if task.Logs == nil {
      task.Logs = []*tes.TaskLog{&tes.TaskLog{}}
    }

    for i := len(task.Logs[0].Logs); i < int(req.Step) + 1; i++ {
      task.Logs[0].Logs = append(task.Logs[0].Logs, &tes.ExecutorLog{})
    }

    existing := task.Logs[0].Logs[req.Step]
    stdout := []byte(existing.Stdout + req.Log.Stdout)
    stderr := []byte(existing.Stderr + req.Log.Stderr)

    // Trim the stdout/err logs to the max size if needed
    if len(stdout) > max {
      stdout = stdout[:max]
    }
    if len(stderr) > max {
      stderr = stderr[:max]
    }

    existing.Stdout = string(stdout)
    existing.Stderr = string(stderr)
  }
  tb.putTask(task)

	return &pbf.UpdateExecutorLogsResponse{}, nil
}

func (tb *TaskBadger) getWorker(id string) (*pbf.Worker, error) {
  key := workerKey(id)
  var item badger.KVItem
  err := tb.db.Get(key, &item)
  if err != nil {
    return nil, err
  }

  data := item.Value()
  var worker pbf.Worker
	proto.Unmarshal(data, &worker)

  // TODO hackity hack
  worker.Id = id
  if worker.Tasks == nil {
    worker.Tasks = map[string]*pbf.TaskWrapper{}
  }

	return &worker, nil
}

func (tb *TaskBadger) putWorker(worker *pbf.Worker) {
  log.Debug("putWorker", worker)
	data, _ := proto.Marshal(worker)
  key := workerKey(worker.Id)
  tb.db.Set(key, data)
}
