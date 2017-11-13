package datastore

import (
	"cloud.google.com/go/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func (d *Datastore) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

	key := datastore.NameKey("Task", req.Id, nil)
	q := datastore.NewQuery("TaskChunk").Ancestor(key)
  task := &tes.Task{}

	for it := d.client.Run(ctx, q); ; {

		c := &chunk{}
		_, err := it.Next(c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
    toTask(c, task)
  }

	return task, nil
}

// ListTasks returns a list of taskIDs
func (d *Datastore) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	resp := &tes.ListTasksResponse{}
  /*
	q := datastore.NewQuery("TaskChunk")

	for it := d.client.Run(ctx, q); ; {

		c := &chunk{}
		_, err := it.Next(c)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		resp.Tasks = append(resp.Tasks, toTask(c))
	}
  */
	return resp, nil
}
