package tes

import (
  "context"
  "time"
)

type TaskGetter func(context.Context, *GetTaskRequest) (*Task, error)

func GetFullTask(ctx context.Context, id string, get TaskGetter) (*Task, error) {
	return get(context.Background(), &GetTaskRequest{
		Id:   id,
		View: TaskView_FULL,
	})
}

func PollTaskContext(ctx context.Context, id string, get TaskGetter) context.Context {
	taskctx, cancel := context.WithCancel(ctx)

	// Start a goroutine that polls the server to watch for a canceled state.
	// If a cancel state is found, "taskctx" is canceled.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-taskctx.Done():
				return
			case <-ticker.C:
        resp, err := get(ctx, &GetTaskRequest{
          Id: id,
          View: TaskView_MINIMAL,
        })
				if TerminalState(resp.GetState()) {
					cancel()
				}
			}
		}
	}()
	return taskctx
}

