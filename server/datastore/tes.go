package datastore

import (
  "google.golang.org/appengine"
  "google.golang.org/appengine/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
)

func (d *Datastore) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
  tk := taskKey(ctx, req.Id)
	res, err := d.GetTasks(ctx, []*datastore.Key{tk}, req.View)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func (d *Datastore) GetTasks(ctx context.Context, keys []*datastore.Key, view tes.TaskView) ([]*tes.Task, error) {

	proplists := make([]datastore.PropertyList, len(keys), len(keys))
	err := datastore.GetMulti(ctx, keys, proplists)
	if err != nil {
		return nil, err
	}

	var tasks []*tes.Task
	var parts []*datastore.Key
	byID := map[string]*tes.Task{}

	for _, props := range proplists {
		task := &tes.Task{}
		if err := unmarshalTask(task, props); err != nil {
			return nil, err
		}

		// Now that we have the task loaded, we know how many attempts/executors there are,
		// so we can determine the keys for the full view parts.
		switch view {
		case tes.Minimal:
			task = task.GetMinimalView()
		case tes.Full:
			tk := taskKey(ctx, task.Id)
			byID[task.Id] = task
			parts = append(parts, contentKey(ctx, tk))
			for attempt, a := range task.Logs {
				parts = append(parts, syslogKey(ctx, tk, uint32(attempt)))
				for index := range a.Logs {
					parts = append(parts, stdoutKey(ctx, tk, uint32(attempt), uint32(index)))
					parts = append(parts, stderrKey(ctx, tk, uint32(attempt), uint32(index)))
				}
			}
		}

		tasks = append(tasks, task)
	}

	// Load the full view parts
	if view == tes.Full {
		proplists := make([]datastore.PropertyList, len(parts), len(parts))
		err := datastore.GetMulti(ctx, parts, proplists)
		merr, isMerr := err.(appengine.MultiError)
		if err != nil && !isMerr {
			return nil, err
		}

		for i, props := range proplists {
			if merr[i] != nil {
				// The task doesn't have this part, e.g. it doesn't have stdout for an executor.
				// That's ok, skip this part.
				continue
			}
			id := parts[i].Parent().StringID()
			task := byID[id]
			if err := unmarshalPart(task, props); err != nil {
				return nil, err
			}
		}
	}

	return tasks, nil
}

func (d *Datastore) GetTasksMemcache(keys []*datastore.Key, view tes.TaskView) ([]*tes.Task, error) {
	/*
		var memcacheKeys []string
		for _, key := range keys {
			memcacheKeys = append(memcacheKeys, key.Name+"-"+view.String())
		}
	*/
	return nil, nil
}

// ListTasks returns a list of taskIDs
func (d *Datastore) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	page := req.PageToken
	size := tes.GetPageSize(req.GetPageSize())
	q := datastore.NewQuery("Task").KeysOnly().Limit(size)

	if page != "" {
		c, err := datastore.DecodeCursor(page)
		if err != nil {
			return nil, err
		}
		q = q.Start(c)
	}

	var keys []*datastore.Key

	it := q.Run(ctx)
	for {
		key, err := it.Next(nil)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	tasks, err := d.GetTasks(ctx, keys, req.View)
	if err != nil {
		return nil, err
	}
	resp := &tes.ListTasksResponse{Tasks: tasks}

	if len(keys) == size {
		c, err := it.Cursor()
		if err != nil {
			return nil, err
		}
		resp.NextPageToken = c.String()
	}

	return resp, nil
}
