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
stdout, stderr, system logs, and input content. It has an ancestor
link to "Task". It does not hold the base task document.
*/

func stdoutKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("stdout-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskFull", k, task)
}

func stderrKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("stderr-%d-%d", e.Attempt, e.Index)
	return datastore.NameKey("TaskFull", k, task)
}

func syslogKey(e *events.Event, task *datastore.Key) *datastore.Key {
	k := fmt.Sprintf("syslog-%s", e.Timestamp)
	return datastore.NameKey("TaskFull", k, task)
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
		_, err := d.client.Put(ctx, syslogKey(e, taskKey), fromEvent(e))
		return err

	case events.Type_EXECUTOR_STDOUT:
		_, err := d.client.Put(ctx, stdoutKey(e, taskKey), fromEvent(e))
		return err

	case events.Type_EXECUTOR_STDERR:
		_, err := d.client.Put(ctx, stderrKey(e, taskKey), fromEvent(e))
		return err

	default:
		_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			res := &task{}
			err := d.client.Get(ctx, taskKey, res)
			if err != nil {
				return err
			}

			task := &tes.Task{}
			toTask(res, task)
			tb := events.TaskBuilder{task}
			tb.WriteEvent(context.Background(), e)

			_, err = d.client.Put(ctx, taskKey, fromTask(task))
			return err
		})
		return err
	}
	return nil
}
