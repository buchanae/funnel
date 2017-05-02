package worker

func PollForCancel(ctx context.Context, rate time.Duration, reader TaskReader) context.Context {
  taskctx, cancel := context.WithCancel(ctx)

  // Start a goroutine that polls the server to watch for a canceled state.
  // If a cancel state is found, "taskctx" is canceled.
  go func() {
    ticker := time.NewTicker(rate)
    defer ticker.Stop()

    for {
    case <-taskctx.Done():
      return
    case <-ticker.C:
      task, err := reader.State()
      GetTask(ctx, &tes.GetTaskRequest{
        Id: r.TaskID,
      })

      if task.State == tes.State_CANCELED {
        cancel()
      }
    }
  }()
  return taskctx
}
