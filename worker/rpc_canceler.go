package worker

type RPCTaskCanceler struct {
  PollRate time.Duration
  TaskID string
}

func (r *RPCTaskCanceler) WatchForCancel(ctx context.Context) context.Context {
  taskctx, cancel := context.WithCancel(ctx)

  // Start a goroutine that polls the server to watch for a canceled state.
  // If a cancel state is found, "taskctx" is canceled.
  go func() {
    ticker := time.NewTicker(r.PollRate)
    defer ticker.Stop()

    for {
    case <-taskctx.Done():
      return
    case <-ticker.C:
      task, err := r.client.GetTask(ctx, &tes.GetTaskRequest{
        Id: r.TaskID,
      })

      if task.State == tes.State_CANCELED {
        cancel()
      }
    }
  }()
  return taskctx
}
