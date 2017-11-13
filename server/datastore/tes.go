package datastore

import (
	"cloud.google.com/go/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func (d *Datastore) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

	taskKey := datastore.NameKey("Task", req.Id, nil)

  var q *datastore.Query
  switch req.View {
  case tes.Minimal:
	  q = datastore.NewQuery("Task").Project("State").Filter("__key__ =", taskKey)
  case tes.Basic:
	  q = datastore.NewQuery("Task").
      Filter("Id = ", req.Id)
  case tes.Full:
	  q = datastore.NewQuery("").Ancestor(taskKey)
  }

	pls := []datastore.PropertyList{}
	keys, err := d.client.GetAll(ctx, q, &pls)
	if err != nil {
		return nil, err
	}

  logger.Debug("VIEW", "view", req.View)

	t := &tes.Task{}
  for i, key := range keys {
    pl := pls[i]
    switch {
    case key.Kind == "Task":
      c := &task{}
      datastore.LoadStruct(c, pl)
	    toTask(c, t)
    default:
      logger.Debug("FULL KEY", key.Name)
    }
  }

	return t, nil
}

// ListTasks returns a list of taskIDs
func (d *Datastore) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	resp := &tes.ListTasksResponse{}
	q := datastore.NewQuery("Task").Project("Id", "State")

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
