package tes

import (
  "github.com/ohsu-comp-bio/funnel/logger"
	"golang.org/x/net/context"
)

func WithLogging(s TaskServiceServer, l logger.Logger) TaskServiceServer {
  return &teslogger{s, l}
}

type teslogger struct {
  tes TaskServiceServer
  log logger.Logger
}

func (l *teslogger) CreateTask(ctx context.Context, task *Task) (*CreateTaskResponse, error) {
  l.log.Debug("Create task", "task", task)
  return l.tes.CreateTask(ctx, task)
}

func (l *teslogger) GetTask(ctx context.Context, req *GetTaskRequest) (*Task, error) {
  l.log.Debug("Get task", "req", req)
  return l.tes.GetTask(ctx, req)
}

func (l *teslogger) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
  l.log.Debug("List tasks", "req", req)
  return l.tes.ListTasks(ctx, req)
}

func (l *teslogger) CancelTask(ctx context.Context, req *CancelTaskRequest) (*CancelTaskResponse, error) {
  l.log.Debug("Cancel task", "req", req)
  return l.tes.CancelTask(ctx, req)
}

func (l *teslogger) GetServiceInfo(ctx context.Context, req *ServiceInfoRequest) (*ServiceInfo, error) {
  l.log.Debug("Get service info", "req", req)
  return l.tes.GetServiceInfo(ctx, req)
}
