package gcp

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
  "cloud.google.com/go/datastore"
)

var log = logger.Sub("gcp")
const bucketName = "buchanae-funnel"
const projectID = "funnel-165618"

type TaskServiceServer struct {}

var startupScript = `
#!/bin/sh
docker run          \
  --group-add 412   \
  --name funnel     \
  -w /var/db/funnel \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/bin/docker:/usr/bin/docker           \
  -v /var/db/funnel:/var/db/funnel             \
  docker.io/ohsucompbio/funnel:latest gce start

shutdown
`


func (gcp *TaskServiceServer) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {

  client, err := datastore.NewClient(ctx, projectID)
  if err != nil {
    return nil, err
  }

	task.Id = util.GenTaskID()
  key := datastore.NameKey("Task", task.Id, nil)
  if _, err := client.Put(ctx, key, task); err != nil {
    return nil, err
  }

  return &tes.CreateTaskResponse{Id: task.Id}, nil
}

func (gcp *TaskServiceServer) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
  return nil, nil
}

func (gcp *TaskServiceServer) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
  return nil, nil
}

// CancelTask cancels a task
func (gcp *TaskServiceServer) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {
	return &tes.CancelTaskResponse{}, nil
}

// GetServiceInfo provides an endpoint for Funnel clients to get information about this server.
// Could include:
// - resource availability
// - support storage systems
// - versions
// - etc.
func (gcp *TaskServiceServer) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{Name: "funnel-gcp-datastore"}, nil
}
