package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
)

// Run configures and runs a Worker
func Run(conf config.Worker, taskID string) error {
	logger.Configure(conf.Logger)
	w, err := NewDefaultWorker(conf, taskID)
	if err != nil {
		return err
	}
	w.Run(context.Background())
	return nil
}
