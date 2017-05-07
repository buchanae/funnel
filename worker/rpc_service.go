package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"time"
)

func NewRPCService(conf config.Worker) (*RPCService, error) {
	client, err := newSchedClient(conf)
	if err != nil {
		return nil, err
	}

	return &RPCService{
    service: &Service{
      TickRate: conf.UpdateRate,
      Timeout: conf.Timeout,
    },
		conf:      conf,
		client:     client,
	}, nil
}

type RPCService struct {
  service *Service
  conf config.Config
  client *schedClient
}

func (rs *RPCService) Run() {
  go rs.service.Run()

  log := logger.New("rpc-worker", "workerID", conf.ID)
  resources := detectResources(rs.conf.Resources)
  // TODO
  hostname := ""

  for range rs.service.Tick {

    r, gerr := w.client.SyncWorker(context.TODO(), &pbf.Worker{
      Id: w.conf.ID
      Hostname: w.hostname,
      Resources: w.resources,
    })

    if gerr != nil {
      log.Error("Couldn't get worker state during sync.", gerr)
      return
    }

    if r.Shutdown {
      w.Stop()
      return
    }

    // Start task runners. runSet will track task IDs
    // to ensure there's only one runner per ID, so it's ok
    // to call this multiple times with the same task ID.
    for _, id := range r.TaskIds {
      go w.Runners.Add(id, func(ctx context.Context, id string) {
        // TODO
        b := w.createBackend(id)
        b.Run(ctx)
      })
    }
  }
}

func (rs *RPCService) Stop() {
  rs.client.Goodbye()
  rs.client.Close()
  rs.service.Stop()
}
