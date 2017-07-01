package tes

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func WithValidation(s TaskService, l logger.Logger) TaskService {
  return &validator{s, l}
}

type validator struct {
  TaskService
  log logger.Logger
}
func (v *validator) CreateTask(ctx context.Context, task *Task) (*CreateTaskResponse, error) {
	if err := Validate(task); err != nil {
		v.log.Error("Invalid task message", "error", err)
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}
  return v.TaskService.CreateTask(ctx, task)
}
