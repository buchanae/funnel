package gcp

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var log = logger.Sub("gcp")

type GCP struct {
  gce GCE
}

func NewGCP(conf Config) error {
  ctx := context.Background()
}


func (gcp *GCP) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {

  // TODO use datastore ID generation? then prefix with time and hash?
	task.Id := util.GenTaskID()

  project := ""
  zone := "us-west1-a"

	// Prepare disk details by setting the specific zone
	for _, disk := range props.Disks {
		dt := localize(zone, "diskTypes", disk.InitializeParams.DiskType)
		disk.InitializeParams.DiskType = dt
	}

	// Create the instance on GCE
	instance := compute.Instance{
		Name:              genWorkerID(),
		Description:       fmt.Sprintf("Funnel worker for task ID %s", task.Id),
		Disks:             props.Disks,
		MachineType:       localize(zone, "machineTypes", props.MachineType),
		Tags:              props.Tags,
		Metadata:          &metadata,
	}

	op, ierr := s.gce.InsertInstance(project, zone, &instance)
	if ierr != nil {
		log.Error("Couldn't insert GCE VM instance", ierr)
		return ierr
	}
	log.Debug("GCE VM instance created", "details", op)

  return &tes.CreateTaskResponse{
    Id: task.Id,
  }, nil
}



func (gcp *GCP) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

  task := tes.Task{}
  var keys []*datastore.Key
  var ents []*tes.Task

  switch req.View {
  case tes.TaskView_FULL:
    fallthrough

  case tes.TaskView_BASIC:
    fallthrough

  default:
  }

	return task, err
}



func (gcp *GCP) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
  resp := tes.ListTasksResponse{}

	return &resp, nil
}



// CancelTask cancels a task
func (gcp *GCP) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {
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
