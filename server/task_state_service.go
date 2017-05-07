package server

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

func UpdateTask() {
	// Reconcile worker's task states with database
	for _, wrapper := range req.Tasks {
		// TODO test transition to self a noop
		task := wrapper.Task
		err := transitionTaskState(tx, task.Id, task.State)
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
			key := append([]byte(req.Id), []byte(task.Id)...)
			tx.Bucket(WorkerTasks).Delete(key)
		}
	}
}

func transitionTaskState(tx *bolt.Tx, id string, state tes.State) error {
	idBytes := []byte(id)
	current := getTaskState(tx, id)

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
		tx.Bucket(TasksQueued).Delete(idBytes)

	case Running, Initializing:
		if current != Unknown && current != Queued && current != Initializing {
			log.Error("Unexpected transition", "current", current, "requested", state)
			return errors.New("Unexpected transition to Initializing")
		}
		tx.Bucket(TasksQueued).Delete(idBytes)

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

	tx.Bucket(TaskState).Put(idBytes, []byte(state.String()))
	log.Info("Set task state", "taskID", id, "state", state.String())
	return nil
}

// UpdateExecutorLogs is an API endpoint that updates the logs of a task.
// This is used by workers to communicate task updates to the server.
func (taskBolt *TaskBolt) UpdateExecutorLogs(ctx context.Context, req *pbf.UpdateExecutorLogsRequest) (*pbf.UpdateExecutorLogsResponse, error) {

	taskBolt.db.Update(func(tx *bolt.Tx) error {
		bL := tx.Bucket(TasksLog)

		// max size (bytes) for stderr and stdout streams to keep in db
		max := taskBolt.conf.MaxExecutorLogSize
		key := []byte(fmt.Sprint(req.Id, req.Step))

		if req.Log != nil {
			// Check if there is an existing task log
			o := bL.Get(key)
			if o != nil {
				// There is an existing log in the DB, load it
				existing := &tes.ExecutorLog{}
				// max bytes to be stored in the db
				proto.Unmarshal(o, existing)

				stdout := []byte(existing.Stdout + req.Log.Stdout)
				stderr := []byte(existing.Stderr + req.Log.Stderr)

				// Trim the stdout/err logs to the max size if needed
				if len(stdout) > max {
					stdout = stdout[:max]
				}
				if len(stderr) > max {
					stderr = stderr[:max]
				}

				req.Log.Stdout = string(stdout)
				req.Log.Stderr = string(stderr)

				// Merge the updates into the existing.
				proto.Merge(existing, req.Log)
				// existing is updated, so set that to req.Log which will get saved below.
				req.Log = existing
			}

			// Save the updated log
			logbytes, _ := proto.Marshal(req.Log)
			tx.Bucket(TasksLog).Put(key, logbytes)
		}

		return nil
	})
	return &pbf.UpdateExecutorLogsResponse{}, nil
}
