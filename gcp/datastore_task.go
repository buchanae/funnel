package gcp

import (
  "cloud.google.com/go/datastore"
  "golang.org/x/net/context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
  "google.golang.org/api/iterator"
)

type Backend interface {
  Submit(*tes.Task) error
}

type DatastoreTaskEventWriter struct {
  client *datastore.Client
  builder *events.TaskBuilder
  task *tes.Task
}

func NewDatastoreTaskEventWriter(project string, task *tes.Task) (*DatastoreTaskEventWriter, error) {
  ctx := context.Background()
  client, err := datastore.NewClient(ctx, project)
  if err != nil {
    return nil, err
  }
  builder := events.NewTaskBuilder(task)
  return &DatastoreTaskEventWriter{client, builder, task}, nil
}

func (d *DatastoreTaskEventWriter) Close() error {
  return d.client.Close()
}

func (d *DatastoreTaskEventWriter) Flush() error {
  ctx := context.Background()
  key := datastore.NameKey(taskKind, d.task.Id, nil)
  _, err := d.client.Put(ctx, key, d.task)
  return err
}

func (d *DatastoreTaskEventWriter) Write(ev *events.Event) error {
	return d.builder.Write(ev)
}



const taskKind = "Task"




type DatastoreTES struct {
  client *datastore.Client
  backend Backend
}

func NewDatastoreTES(project string, b Backend) (*DatastoreTES, error) {
  ctx := context.Background()
  client, err := datastore.NewClient(ctx, project)
  if err != nil {
    return nil, err
  }
  return &DatastoreTES{client, b}, nil
}

func (d *DatastoreTES) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {

	if err := tes.Validate(task); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

  task.Id = util.GenTaskID()
  task.State = tes.Queued

  key := datastore.NameKey(taskKind, task.Id, nil)
  _, err := d.client.Put(ctx, key, task)
  if err != nil {
    return nil, err
  }

  err = d.backend.Submit(task)
  if err != nil {
    // TODO handle error better
    return nil, err
  }

  return &tes.CreateTaskResponse{Id: task.Id}, nil
}

func (d *DatastoreTES) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

  task := &tes.Task{}
  key := datastore.NameKey(taskKind, req.Id, nil)

  err := d.client.Get(ctx, key, task)
  if err != nil {
    return nil, err
  }
  if req.View == tes.TaskView_MINIMAL {
    task = minView(task)
  }

  return task, nil
}

func minView(t *tes.Task) *tes.Task {
  o := &tes.Task{}
  o.Id = t.Id
  o.State = t.State
  return o
}

// ListTasks returns a list of taskIDs
func (d *DatastoreTES) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

  resp := &tes.ListTasksResponse{}
  q := datastore.NewQuery(taskKind)

  for t := d.client.Run(ctx, q); ; {
    task := &tes.Task{}

    _, err := t.Next(task)
    if err == iterator.Done {
      break
    }
    if err != nil {
      return nil, err
    }

    if req.View == tes.TaskView_MINIMAL {
      task = minView(task)
    }
    resp.Tasks = append(resp.Tasks, task)
  }
  return resp, nil
}

// CancelTask cancels a task
func (d *DatastoreTES) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {

  key := datastore.NameKey(taskKind, req.Id, nil)
  _, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

    task := &tes.Task{}
    err := d.client.Get(ctx, key, task)
    if err != nil {
      return err
    }
    task.State = tes.Canceled

    _, err = d.client.Put(ctx, key, task)
    return err
  })
  if err != nil {
    return nil, err
  }
  return &tes.CancelTaskResponse{}, nil
}

func (d *DatastoreTES) RestartTask(ctx context.Context, req *tes.RestartTaskRequest) (*tes.RestartTaskResponse, error) {

  _, err := d.CancelTask(ctx, &tes.CancelTaskRequest{Id: req.Id})
  if err != nil {
    return nil, err
  }

  task := &tes.Task{}
  key := datastore.NameKey(taskKind, req.Id, nil)
  err = d.client.Get(ctx, key, task)
  if err != nil {
    return nil, err
  }
  task.Logs = nil
  task.State = tes.Queued
  task.Id = ""
  _, err = d.CreateTask(ctx, task)

  return nil, err
}

// GetServiceInfo provides an endpoint for Funnel clients to get information about this server.
func (d *DatastoreTES) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{Name: "funnel-gcp-datastore"}, nil
}
