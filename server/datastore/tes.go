package datastore

import (
	"cloud.google.com/go/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func (d *Datastore) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

	key := datastore.NameKey("Task", req.Id, nil)

	t := &task{}
	err := d.client.Get(ctx, key, t)
	if err != nil {
		return nil, err
	}

	return toTask(t), nil
}

// ListTasks returns a list of taskIDs
func (d *Datastore) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	resp := &tes.ListTasksResponse{}
	q := datastore.NewQuery("Task")

	for t := d.client.Run(ctx, q); ; {
		task := &tes.Task{}

		_, err := t.Next(task)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		resp.Tasks = append(resp.Tasks, task)
	}
	return resp, nil
}
