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

The "Task" kind is an empty, root Entity which allows grouping the
chunks of data which make up the various task views, and allows ListTasks
to easily query and sort the task IDs for pagination.

The "TaskChunk" kind makes up all of the bits of task data:
base task, state, basic logs, input content, stdout, and stderr.
TaskChunk entities have an ancestor ID pointing to the Task entity.
These chunks have the following design for keys:

  0-content                        input content
  0-executor-stdout-ATTEMPT-INDEX  executor stdout string
  0-executor-stderr-ATTEMPT-INDEX  executor stderr string
  0-syslog-TIMESTAMP               system logs
  1-base                           the base task, as given to CreateTask
  1-event-TIMESTAMP                task/exec log events
  2-state                          state string

The "0-", "1-", and "2-" prefixes are used to filter task data in order
to provide efficient read/write access to different task views:

  "0-" is full view
  "1-" is basic view
  "2-" is minimal view

This is because Datastore allows filtering entity keys with ">" (greater than)
only in ascending order. A GetTask or ListTasks call may filter entity keys
this way to get back all chunks in one query. For example, getting the basic view
means filtering for IDs greater than "0-" (the full view).

Datastore requires that the entire Entity be sent in each call to Put,
so stdout, stderr, and state are stored separately to allow for easy, efficient updates.
*/
func stdoutKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("0-executor-stdout-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskChunk", k, task)
}

func stderrKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("0-executor-stderr-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskChunk", k, task)
}

func logEventKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("1-event-%d", e.Timestamp)
	return datastore.NameKey("TaskChunk", k, task)
}

func (d *Datastore) WriteEvent(ctx context.Context, e *events.Event) error {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	//contentKey := datastore.NameKey("TaskChunk", "0-content", taskKey)
	baseKey := datastore.NameKey("TaskChunk", "1-base", taskKey)
	stateKey := datastore.NameKey("TaskChunk", "2-state", taskKey)

	switch e.Type {
  case events.Type_SYSTEM_LOG:
    // TODO

	case events.Type_TASK_CREATED:
		task := e.GetTask()
		_, err := d.client.Put(ctx, baseKey, fromTask(task))
		if err != nil {
			return err
		}

		_, err = d.client.Put(ctx, stateKey, fromState(tes.Queued))
		return err

	case events.Type_TASK_STATE:
		_, err := d.client.Put(ctx, stateKey, fromState(e.GetState()))
		return err

	case events.Type_TASK_START_TIME, events.Type_TASK_END_TIME,
		events.Type_TASK_OUTPUTS, events.Type_TASK_METADATA,
	  events.Type_EXECUTOR_START_TIME, events.Type_EXECUTOR_END_TIME,
		events.Type_EXECUTOR_EXIT_CODE:

    _, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
      key := logEventKey(e, taskKey)
      _, err := d.client.Put(ctx, key, fromLogEvent(e))
      return err
    })
    return err

	case events.Type_EXECUTOR_STDOUT:
		_, err := d.client.Put(ctx, stdoutKey(e, taskKey), fromStdout(e.GetStdout()))
		return err

	case events.Type_EXECUTOR_STDERR:
		_, err := d.client.Put(ctx, stderrKey(e, taskKey), fromStderr(e.GetStderr()))
		return err
	}
	return nil
}
