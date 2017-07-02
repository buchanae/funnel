package gcp

import (
  "bytes"
  "sync"
	"github.com/ohsu-comp-bio/funnel/logger"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
  "cloud.google.com/go/storage"
)

var log = logger.Sub("gcp")

type GCP struct {
  gce GCE
}

func NewGCP(conf Config) error {
  ctx := context.Background()
}

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

type partial struct {
  task *tes.Task
  handle *storage.ObjectHandle
}

type partials []partial
func (ps partials) read(ctx) error {
}
func (ps partials) write(ctx) error {
  ctx, cancel = context.WithCancel(ctx)
  defer cancel()

  errs := make(chan error)
  for _, p := range ps {
    go func(p partial) {
      w := p.handle.NewWriter(ctx)
      errs <- marshaler.Marshal(w, p.task)
      errs <- w.Close()
    }(p)
  }

  for i := 0; i < len(ps) * 2; i++ {
    if err := <-errs; err != nil {
      return error
    }
  }

  return nil
}



type taskHandle struct {
  bkt string
  id string
}
func (h *taskHandle) URL() string {
  return fmt.Sprintf("tasks/%s/", h.id)
}
func (h *taskHandle) Base() *storage.ObjectHandle {
  p := fmt.Sprintf("tasks/%s/task.base.json", h.id)
  return h.bkt.Object(p)
}
func (h *taskHandle) State() *storage.ObjectHandle {
  p := fmt.Sprintf("tasks/%s/task.state.json", h.id)
  return h.bkt.Object(p)
}
func (h *taskHandle) Stdout() *storage.ObjectHandle {
  p := fmt.Sprintf("tasks/%s/task.stdout.json", h.id)
  return h.bkt.Object(p)
}
func (h *taskHandle) Stderr() *storage.ObjectHandle {
  p := fmt.Sprintf("tasks/%s/task.stderr.json", h.id)
  return h.bkt.Object(p)
}
func (h *taskHandle) Contents() *storage.ObjectHandle {
  p := fmt.Sprintf("tasks/%s/task.contents.json", h.id)
  return h.bkt.Object(p)
}



func (gcp *GCP) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
  ctx, cancel = context.WithCancel(ctx)
  defer cancel()

	task.Id := util.GenTaskID()
  project := "funnel-165618"
  zone := "us-west1-a"
  machineType := "n1-standard-1"
  diskSize := int64(50)

  h := taskHandle{store, "funnel", task.Id}
  // Split task input contents into a separate task.
  // This modifies "task"
  contents := splitInputContents(task)
  taskURL := h.URL()

  // Write the task message as JSON
  p := []partials{
    // Base task message
    {
      task: task,
      handle: h.Base(),
    },
    // Task state
    {
      task: &tes.Task{State: tes.State_INITIALIZING},
      handle: h.State(),
    },
    // Task input contents
    {
      task: contents,
      handle: h.Contents(),
    },
    // Task stdout
    {
      task: &tes.Task{},
      handle: h.Stdout(),
    },
    // Task stderr
    {
      task: &tes.Task{},
      handle: h.Stderr(),
    },
  }

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

  return &tes.CreateTaskResponse{
    Id: task.Id,
  }, nil
}



func (gcp *GCP) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  bucketName := "funnel"
  project := "funnel-165618"
  zone := "us-west1-a"
  base := "/tasks/" + task.Id + "/"
  task := &tes.Task{}

  bkt := store.Bucket(bucketName)

  switch req.View {
  case tes.TaskView_FULL:
    ps = append(ps, partial{
      task: task,
      handle: h.Stdout(),
    })
    ps = append(ps, partial{
      task: task,
      handle: h.Stderr(),
    })
    ps = append(ps, partial{
      task: task,
      handle: h.Contents(),
    })
    fallthrough

  case tes.TaskView_BASIC:
    ps = append(ps, partial{
      task: task,
      handle: h.Base(),
    })
    fallthrough

  default:
    ps = append(ps, partial{
      task: task,
      handle: h.State(),
    })
  }

  if err := ps.read(ctx); err != nil {
    return nil, err
  }

	return task, nil
}



func (gcp *GCP) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

  bucket := "funnel"
  project := "funnel-165618"
  zone := "us-west1-a"
  resp := tes.ListTasksResponse{}

  objects, _ := objs.List(bucket).Prefix("tasks").Do()
  for _, obj := range objects.Items {
    call := objs.Get(url.bucket, obj.Name)

    // TODO get tasks in parallel
    task, err := gcp.GetTask(ctx, &tes.GetTaskRequest{id, req.View})

    if err != nil {
      return err
    }
    resp.Tasks = append(resp.Tasks, task)
  }

	return &resp, nil
}



// CancelTask cancels a task
func (gcp *GCP) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {

  bucket := "funnel"
  project := "funnel-165618"
  zone := "us-west1-a"

  st := strings.NewReader(tes.State_CANCELED.String())

	_, err := objs.Insert(bucket, &storage.Object{
		Name: "tasks/" + req.Id + "/task.state",
  }).Media(st).Do()

	return &tes.CancelTaskResponse{}, nil
}

// GetServiceInfo provides an endpoint for Funnel clients to get information about this server.
// Could include:
// - resource availability
// - support storage systems
// - versions
// - etc.
func (gcp *GCP) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{Name: "funnel-gcp"}, nil
}

func genWorkerID() string {
	u, _ := uuid.NewV4()
	return fmt.Sprintf("funnel-worker-%s", u.String())
}

// localize helps make a resource string zone-specific
func localize(zone, resourceType, val string) string {
	return fmt.Sprintf("zones/%s/%s/%s", zone, resourceType, val)
}
