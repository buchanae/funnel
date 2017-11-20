package datastore

import (
  "google.golang.org/appengine/datastore"
	oldctx "golang.org/x/net/context"
  "context"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

func (d *Datastore) WriteEvent(ctx context.Context, e *events.Event) error {
  tk := taskKey(ctx, e.Id)

	switch e.Type {

	case events.Type_TASK_CREATED:
		_, err := datastore.Put(ctx, tk, marshalTask(e.GetTask()))
		if err != nil {
			return err
		}

	case events.Type_EXECUTOR_STDOUT:
		_, err := datastore.Put(ctx, stdoutKey(ctx, tk, e.Attempt, e.Index), marshalEvent(e))
		return err

	case events.Type_EXECUTOR_STDERR:
		_, err := datastore.Put(ctx, stderrKey(ctx, tk, e.Attempt, e.Index), marshalEvent(e))
		return err

	default:
    err := datastore.RunInTransaction(ctx, func(ctx oldctx.Context) error {
			props := datastore.PropertyList{}
			err := datastore.Get(ctx, tk, &props)
			if err != nil {
				return err
			}

			task := &tes.Task{}
			unmarshalTask(task, props)
			tb := events.TaskBuilder{task}
			err = tb.WriteEvent(context.Background(), e)
			if err != nil {
				return err
			}

			_, err = datastore.Put(ctx, tk, marshalTask(task))
			return err
		}, nil)
		return err
	}
	return nil
}
