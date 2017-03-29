package worker

import (
	"context"
	"errors"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	pbr "tes/server/proto"
	"tes/util"
	"time"
)

// Sync syncs the worker's state with the server. It reports job state changes,
// handles signals from the server (new job, cancel job, etc), reports resources, etc.
//
// TODO Sync should probably use a channel to sync data access.
//      Probably only a problem for test code, where Sync is called directly.
func (w *Worker) Sync() {
	r, gerr := w.sched.GetWorker(context.TODO(), &pbr.GetWorkerRequest{Id: w.conf.ID})

	if gerr != nil {
		log.Error("Couldn't get worker state during sync.", gerr)
		return
	}

	// Reconcile server state with worker state.
	rerr := w.reconcile(r.Jobs)
	if rerr != nil {
		// TODO what's the best behavior here?
		log.Error("Couldn't reconcile worker state.", rerr)
		return
	}

	// Worker data has been updated. Send back to server for database update.
	r.Resources = w.resources
	r.State = w.state

	// Merge metadata
	if r.Metadata == nil {
		r.Metadata = map[string]string{}
	}
	for k, v := range w.conf.Metadata {
		r.Metadata[k] = v
	}

	_, err := w.sched.UpdateWorker(r)
	if err != nil {
		log.Error("Couldn't save worker update. Recovering.", err)
	}
}

// reconcile merges the server state with the worker state:
// - identifies new jobs and starts new runners for them
// - identifies canceled jobs and cancels existing runners
// - updates pbr.Job structs with current job state (running, complete, error, etc)
func (w *Worker) reconcile(jobs map[string]*pbr.JobWrapper) error {
	var (
		Unknown  = pbe.State_Unknown
		Canceled = pbe.State_Canceled
	)

	// Combine job IDs from response with job IDs from ctrls so we can reconcile
	// both sets below.
	jobIDs := map[string]bool{}
	for jobID := range w.Ctrls {
		jobIDs[jobID] = true
	}
	for jobID := range jobs {
		jobIDs[jobID] = true
	}

	for jobID := range jobIDs {
		jobSt := pbe.State_Unknown
		runSt := pbe.State_Unknown

		ctrl := w.Ctrls[jobID]
		if ctrl != nil {
			runSt = ctrl.State()
		}

		wrapper := jobs[jobID]
		var job *pbe.Job
		if wrapper != nil {
			job = wrapper.Job
			jobSt = job.GetState()
		}

		if isComplete(jobSt) {
			delete(w.Ctrls, jobID)
		}

		switch {
		case jobSt == Unknown && runSt == Unknown:
			// Edge case. Shouldn't be here, but log just in case.
			fallthrough
		case isComplete(jobSt) && runSt == Unknown:
			// Edge case. Job is complete and there's no ctrl. Do nothing.
			fallthrough
		case isComplete(jobSt) && runSt == Canceled:
			// Edge case. Job is complete and the ctrl is canceled. Do nothing.
			// This shouldn't happen but it's better to check for it anyway.
			//
			// Log so that these unexpected cases can be explored.
			log.Error("Edge case during worker reconciliation. Recovering.",
				"jobst", jobSt, "runst", runSt)

		case isActive(jobSt) && runSt == Canceled:
			// Edge case. Server says running but ctrl says canceled.
			// Possibly the worker is shutting down due to a local signal
			// and canceled its jobs.
			job.State = runSt

		case jobSt == runSt:
			// States match, do nothing.

		case isActive(jobSt) && runSt == Unknown:
			// Job needs to be started.
			ctrl := NewJobControl()
			go runJob(ctrl, wrapper)
			w.Ctrls[jobID] = ctrl
			job.State = ctrl.State()

		case isActive(jobSt) && runSt != Unknown:
			// Job is running, update state.
			job.State = runSt

		case jobSt == Canceled && runSt != Unknown:
			// Job is canceled.
			// ctrl.Cancel() is idempotent, so blindly cancel and delete.
			ctrl.Cancel()

		case jobSt == Unknown && runSt != Unknown:
			// Edge case. There's a ctrl for a non-existent job. Delete it.
			// TODO is it better to leave it? Continue in absence of explicit command principle?
			ctrl.Cancel()
			delete(w.Ctrls, jobID)

		case isComplete(jobSt) && isActive(runSt):
			// Edge case. The job is complete but the ctrl is still running.
			// This shouldn't happen but it's better to check for it anyway.
			// TODO better to update job state?
			ctrl.Cancel()

		default:
			log.Error("Unhandled case during worker reconciliation.",
				"job", job, "ctrl", ctrl)
			return errors.New("Unhandled case during worker reconciliation")
		}
	}
	return nil
}

func isActive(s pbe.State) bool {
	return s == pbe.State_Queued || s == pbe.State_Initializing || s == pbe.State_Running
}

func isComplete(s pbe.State) bool {
	return s == pbe.State_Complete || s == pbe.State_Error || s == pbe.State_SystemError
}
