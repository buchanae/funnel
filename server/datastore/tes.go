package datastore

import (
	"cloud.google.com/go/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func (d *Datastore) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

	key := datastore.NameKey("Task", req.Id, nil)
	t := &tes.Task{}

	c := &task{}
	err := d.client.Get(ctx, key, c)
	if err != nil {
		return nil, err
	}
	toTask(c, t)

	return t, nil
}

// ListTasks returns a list of taskIDs
func (d *Datastore) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	resp := &tes.ListTasksResponse{}
	q := datastore.NewQuery("Task")

	for it := d.client.Run(ctx, q); ; {

		c := &task{}
		_, err := it.Next(c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		t := &tes.Task{}
		toTask(c, t)
		resp.Tasks = append(resp.Tasks, t)
	}
	return resp, nil
}
