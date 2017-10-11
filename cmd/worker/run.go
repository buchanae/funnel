package worker

import (
	"context"
  "fmt"
	"github.com/ohsu-comp-bio/funnel/config"
  "github.com/"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/worker"
)

// Run configures and runs a Worker
func Run(conf config.Worker, taskID string) error {
	logger.Configure(conf.Logger)
	w, err := worker.NewDefaultWorker(conf)
	if err != nil {
		return err
	}

  var getter tes.TaskGetter
  switch conf.TaskReader {
  case "rpc":
  case "dyanmo":
    getter = 
  case "elastic":
    getter = 
  }

	tesc, err := rpc.NewTESClient(r.ServerAddress, r.ServerPassword)
  if err != nil {
    return err
  }

  ctx := context.Background()
  task, err := tes.GetFullTask(ctx, taskID, getter)
	if err != nil {
		return fmt.Errorf("can't connect TES client %s", err)
	}

  ctx = tes.PollTaskContext(ctx, taskID, getter)
	w.Run(ctx, task)
	return nil
}
