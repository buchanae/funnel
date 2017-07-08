package tes

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func WithValidation(s TaskServiceServer, l logger.Logger) TaskServiceServer {
  return &validator{s, l}
}

type validator struct {
  tes TaskServiceServer
  log logger.Logger
}
func (v *validator) CreateTask(ctx context.Context, task *Task) (*CreateTaskResponse, error) {
	if err := Validate(task); err != nil {
		v.log.Error("Invalid task message", "error", err)
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
  return v.tes.CreateTask(ctx, task)
}
func (v *validator) GetTask(ctx context.Context, req *GetTaskRequest) (*Task, error) {
  return v.tes.GetTask(ctx, req)
}
func (v *validator) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
  return v.tes.ListTasks(ctx, req)
}
func (v *validator) CancelTask(ctx context.Context, req *CancelTaskRequest) (*CancelTaskResponse, error) {
  return v.tes.CancelTask(ctx, req)
}
func (v *validator) GetServiceInfo(ctx context.Context, req *ServiceInfoRequest) (*ServiceInfo, error) {
  return v.tes.GetServiceInfo(ctx, req)
}
