package tes

import (
  "github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
)

func WithLogging(s TaskService, l logger.Logger) TaskService {
  return logger{s, l}
}

type logger struct {
  TaskService
  log logger.Logger
}

func (l *logger) CreateTask(ctx context.Context, task *Task) (*CreateTaskResponse, error) {
  l.log.Debug("Create task", "task", task)
  return l.TaskService.CreateTask(ctx, task)
}

func (l *logger) GetTask(ctx context.Context, req *GetTaskRequest) (*Task, error) {
  l.log.Debug("Get task", "req", req)
  return l.TaskService.GetTask(ctx, req)
}

func (l *logger) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
  l.log.Debug("List tasks", "req", req)
  return l.TaskService.ListTasks(ctx, req)
}

func (l *logger) CancelTask(ctx context.Context, req *CancelTaskRequest) (*CancelTaskResponse, error) {
  l.log.Debug("Cancel task", "req", req)
  return l.TaskService.CancelTask(ctx, req)
}

func (l *logger) GetServiceInfo(ctx context.Context, req *ServiceInfoRequest) (*ServiceInfo, error) {
  l.log.Debug("Get service info", "req", req)
  return l.TaskService.GetServiceInfo(ctx, req)
}
