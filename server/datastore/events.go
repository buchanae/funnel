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

"Task" holds the basic task view. It has an index
on ID and State, which allows projecting the minimal view.

"TaskFull" holds multiple types of documents making up the full view:
stdout, stderr, system logs, and input content. It does not hold the
base task document. It has an ancestor links to "Task".
*/

func stdoutKey(e *events.Event) *datastore.Key {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	k := fmt.Sprintf("stdout-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskFull", k, taskKey)
}

func stderrKey(e *events.Event) *datastore.Key {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	k := fmt.Sprintf("stderr-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskFull", k, taskKey)
}

func syslogKey(e *events.Event) *datastore.Key {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	k := fmt.Sprintf("syslog-%s", e.Timestamp)
	return datastore.NameKey("TaskFull", k, taskKey)
}

func (d *Datastore) WriteEvent(ctx context.Context, e *events.Event) error {
	taskKey := datastore.NameKey("Task", e.Id, nil)
	// TODO
	//contentKey := datastore.NameKey("TaskChunk", "0-content", taskKey)

	switch e.Type {

	case events.Type_TASK_CREATED:
		_, err := d.client.Put(ctx, taskKey, fromTask(e.GetTask()))
		if err != nil {
			return err
		}

	case events.Type_SYSTEM_LOG:
		// Just in case there's a system log that isn't task-specific.
		if e.Id == "" {
			return nil
		}
		_, err := d.client.Put(ctx, syslogKey(e), fromEvent(e))
		return err

	case events.Type_EXECUTOR_STDOUT:
		_, err := d.client.Put(ctx, stdoutKey(e), fromEvent(e))
		return err

	case events.Type_EXECUTOR_STDERR:
		_, err := d.client.Put(ctx, stderrKey(e), fromEvent(e))
		return err

	default:
		_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			res := &task{}
			err := tx.Get(taskKey, res)
			if err != nil {
				return err
			}

			task := &tes.Task{}
			toTask(res, task)
			tb := events.TaskBuilder{task}
			err = tb.WriteEvent(context.Background(), e)
			if err != nil {
				return err
			}

			_, err = tx.Put(taskKey, fromTask(task))
			return err
		})
		return err
	}
	return nil
}
