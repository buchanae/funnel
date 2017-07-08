package gcp

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
  "cloud.google.com/go/storage"
  "google.golang.org/api/iterator"
)

var log = logger.Sub("gcp")
const bucketName = "buchanae-funnel"

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

  /*
func (gcp *TaskServiceServer) StartWorker(ctx context.Context) {
  //project := "funnel-165618"
  //zone := "us-west1-a"
  //machineType := "n1-standard-1"
  //diskSize := int64(50)
  // Create a worker
	instance := compute.Instance{
		Name:              genWorkerID(),
		Description:       fmt.Sprintf("Funnel worker for task ID %s", task.Id),
		Disks: []*compute.AttachedDisk{
      {
        AutoDelete: true,
        Boot: true,
        InitializeParams: &compute.AttachedDiskInitializeParams{
          DiskSizeGb: diskSize,
        },
      },
    },
		MachineType:       localize(zone, "machineTypes", machineType,
    Metadata: &compute.Metadata{
      Items: []*compute.MetadataItems{
        {
          Key: "taskID",
          Value: &task.Id,
        },
        {
          Key: "taskURL",
          Value: &taskURL,
        },
        {
          Key: "startup-script",
          Value: &startupScript,
        },
      },
    },
	}

	op, ierr := comp.InsertInstance(project, zone, &instance)
	if ierr != nil {
		log.Error("Couldn't insert GCE VM instance", ierr)
		return nil, ierr
	}
	log.Debug("GCE VM instance created", "details", op)
}
  */

func (gcp *TaskServiceServer) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {

  client, err := storage.NewClient(ctx)
  if err != nil {
    return nil, err
  }

	task.Id = util.GenTaskID()
  bkt := taskBucket{client.Bucket(bucketName), task.Id}

  // Write MINIMAL task view to object storage
  /*
  min := viewTask(task, tes.TaskView_MINIMAL)
  err = writeTask(min, bkt.minimal().NewWriter(ctx))
  if err != nil {
    return nil, err
  }
  */

  // Write BASIC task view to object storage
  basic := viewTask(task, tes.TaskView_BASIC)
  err = writeTask(basic, bkt.basic().NewWriter(ctx))
  if err != nil {
    return nil, err
  }

  // Write FULL task view to object storage
  /*
  full := viewTask(task, tes.TaskView_FULL)
  err = writeTask(full, bkt.full().NewWriter(ctx))
  if err != nil {
    return nil, err
  }
  */

  return &tes.CreateTaskResponse{Id: task.Id}, nil
}

func (gcp *TaskServiceServer) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

  client, err := storage.NewClient(ctx)
  if err != nil {
    return nil, err
  }

  bkt := taskBucket{client.Bucket(bucketName), req.Id}
  var obj *storage.ObjectHandle

  switch req.View {
  case tes.TaskView_FULL:
    obj = bkt.full()

  case tes.TaskView_BASIC:
    obj = bkt.basic()

  default:
    obj = bkt.minimal()
  }

  r, err := obj.NewReader(ctx)
  if err == storage.ErrObjectNotExist {
    // TODO return proper 404 error
    return nil, err
  }
  if err != nil {
    return nil, err
  }

  return readTask(r)
}

func (gcp *TaskServiceServer) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

  client, err := storage.NewClient(ctx)
  if err != nil {
    return nil, err
  }

  resp := tes.ListTasksResponse{}
  bkt := client.Bucket(bucketName)

  it := bkt.Objects(ctx, &storage.Query{
    Prefix: "tasks/",
    //Delimiter: "/",
  })

  for {
    attrs, err := it.Next()
    if err == iterator.Done {
      break
    }
    if err != nil {
      return nil, err
    }
    id := attrs.Name

    // Get task view, append to response
    // TODO get tasks in parallel
    //task, err := gcp.GetTask(ctx, &tes.GetTaskRequest{id, req.View})
    //if err != nil {
      //return nil, err
    //}
    resp.Tasks = append(resp.Tasks, &tes.Task{
      Id: id,
    })
  }

	return &resp, nil
}

// CancelTask cancels a task
func (gcp *TaskServiceServer) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {

  client, err := storage.NewClient(ctx)
  if err != nil {
    return nil, err
  }

  bkt := taskBucket{client.Bucket(bucketName), req.Id}
  w := bkt.cancel().NewWriter(ctx)
  if err := w.Close(); err != nil {
    return nil, err
  }

	return &tes.CancelTaskResponse{}, nil
}

// GetServiceInfo provides an endpoint for Funnel clients to get information about this server.
// Could include:
// - resource availability
// - support storage systems
// - versions
// - etc.
func (gcp *TaskServiceServer) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{Name: "funnel-gcp"}, nil
}
