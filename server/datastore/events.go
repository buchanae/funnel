package datastore

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

/*
Entity group and key structure:

Task                             initial Task message as given to CreateTask)
- state                          state string
- content-INDEX                  input content
- task-log-ATTEMPT               tes.TaskLog without ExecutorLogs
- executor-log-ATTEMPT-INDEX:    tes.ExecutorLog without stdout/err
- executor-stdout-ATTEMPT-INDEX  executor stdout string
- executor-stderr-ATTEMPT-INDEX  executor stderr string
*/

func (d *Datastore) WriteEvent(ctx context.Context, e *events.Event) error {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	stateKey := datastore.NameKey("TaskState", "state", taskKey)
	attemptName := fmt.Sprintf("task-log-%d", e.Attempt)
	attemptKey := datastore.NameKey("TaskLog", attemptName, taskKey)
	execIndexName := fmt.Sprintf("executor-log-%d-%d", e.Attempt, e.Index)
	execIndexKey := datastore.NameKey("ExecutorLog", execIndexName, taskKey)
	execStdoutName := fmt.Sprintf("executor-stdout-%d-%d", e.Attempt, e.Index)
	execStdoutKey := datastore.NameKey("ExecutorStdout", execStdoutName, taskKey)
	execStderrName := fmt.Sprintf("executor-stderr-%d-%d", e.Attempt, e.Index)
	execStderrKey := datastore.NameKey("ExecutorStderr", execStderrName, taskKey)

	switch e.Type {
	case events.Type_TASK_CREATED:
		task := e.GetTask()
		_, err := d.client.Put(ctx, taskKey, fromTask(task))
		if err != nil {
			return err
		}

		_, err = d.client.Put(ctx, stateKey, fromState(tes.Queued))
		return err

	case events.Type_TASK_STATE:
		_, err := d.client.Put(ctx, stateKey, fromState(e.GetState()))
		return err

	case events.Type_TASK_START_TIME, events.Type_TASK_END_TIME,
		events.Type_TASK_OUTPUTS, events.Type_TASK_METADATA:
		return d.updateTaskLog(ctx, attemptKey, e)

	case events.Type_EXECUTOR_START_TIME, events.Type_EXECUTOR_END_TIME,
		events.Type_EXECUTOR_EXIT_CODE:
		return d.updateExecLog(ctx, execIndexKey, e)

	case events.Type_EXECUTOR_STDOUT:
		_, err := d.client.Put(ctx, execStdoutKey, fromStdout(e.GetStdout()))
		return err

	case events.Type_EXECUTOR_STDERR:
		_, err := d.client.Put(ctx, execStderrKey, fromStderr(e.GetStderr()))
		return err
	}
	return nil
}

func (d *Datastore) updateTaskLog(ctx context.Context, key *datastore.Key, e *events.Event) error {
	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

		l := &tes.TaskLog{}
		err := d.client.Get(ctx, key, l)
		if err != nil {
			return err
		}

		switch e.Type {
		case events.Type_TASK_START_TIME:
			l.StartTime = e.GetStartTime()
		case events.Type_TASK_END_TIME:
			l.EndTime = e.GetEndTime()
		case events.Type_TASK_OUTPUTS:
			l.Outputs = e.GetOutputs().Value
		case events.Type_TASK_METADATA:
			l.Metadata = e.GetMetadata().Value
		}

		_, err = d.client.Put(ctx, key, fromTaskLog(l))
		return err
	})
	return err
}

func (d *Datastore) updateExecLog(ctx context.Context, key *datastore.Key, e *events.Event) error {
	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

		l := &tes.ExecutorLog{}
		err := d.client.Get(ctx, key, l)
		if err != nil {
			return err
		}

		switch e.Type {
		case events.Type_EXECUTOR_START_TIME:
			l.StartTime = e.GetStartTime()
		case events.Type_EXECUTOR_END_TIME:
			l.EndTime = e.GetEndTime()
		case events.Type_EXECUTOR_EXIT_CODE:
			l.ExitCode = e.GetExitCode()
		}

		_, err = d.client.Put(ctx, key, fromExecLog(l))
		return err
	})
	return err
}
