package memory

import (
  "context"
	oldCtx "golang.org/x/net/context"
  "github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/events"
)

type Memory struct {
  // task ID => task
  Tasks map[string]*tes.Task
}

func NewMemory() *Memory {
  return &Memory{Tasks: map[string]*tes.Task{}}
}

func (m *Memory) GetTask(ctx oldCtx.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
  if task, ok := m.Tasks[req.Id]; ok {
    return task, nil
  }
  return nil, tes.ErrNotFound
}

func (m *Memory) ListTasks(ctx oldCtx.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
  resp := &tes.ListTasksResponse{}
  for _, task := range m.Tasks {
    resp.Tasks = append(resp.Tasks, task)
  }
  return resp, nil
}

func (m *Memory) WriteEvent(ctx context.Context, req *events.Event) error {
  switch req.Type {
  case events.Type_TASK_CREATED:
    task := req.GetTask()
    m.Tasks[task.Id] = task
  case events.Type_TASK_STATE:
    if task, ok := m.Tasks[req.Id]; ok {
      task.State = req.GetState()
    }
    return tes.ErrNotFound
  }
  return nil
}
