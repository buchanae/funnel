package server

import (
	"context"
	"fmt"
	workerCmd "github.com/ohsu-comp-bio/funnel/cmd/worker"
	"github.com/ohsu-comp-bio/funnel/compute/batch"
	"github.com/ohsu-comp-bio/funnel/compute/gridengine"
	"github.com/ohsu-comp-bio/funnel/compute/htcondor"
	"github.com/ohsu-comp-bio/funnel/compute/local"
	"github.com/ohsu-comp-bio/funnel/compute/noop"
	"github.com/ohsu-comp-bio/funnel/compute/pbs"
	"github.com/ohsu-comp-bio/funnel/compute/scheduler"
	"github.com/ohsu-comp-bio/funnel/compute/slurm"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	schedProto "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/server/boltdb"
	"github.com/ohsu-comp-bio/funnel/server/datastore"
	"github.com/ohsu-comp-bio/funnel/server/dynamodb"
	"github.com/ohsu-comp-bio/funnel/server/elastic"
	"github.com/ohsu-comp-bio/funnel/server/mongodb"
	"strings"
)

// Run runs the "server run" command.
func Run(ctx context.Context, conf config.Config) error {
	log := logger.NewLogger("server", conf.Server.Logger)
	s, err := NewServer(conf, log)
	if err != nil {
		return err
	}
	return s.Run(ctx)
}

// Server is a Funnel server + scheduler.
type Server struct {
	*server.Server
	*scheduler.Scheduler
}

// NewServer returns a new Funnel server + scheduler based on the given config.
func NewServer(conf config.Config, log *logger.Logger) (*Server, error) {
	log.Debug("NewServer", "config", conf)

	var reader tes.ReadOnlyServer
	var nodes schedProto.SchedulerServiceServer
	var sched *scheduler.Scheduler
	var queue scheduler.TaskQueue
	writers := events.MultiWriter{}

	switch strings.ToLower(conf.Server.Database) {
	case "boltdb":
		b, err := boltdb.NewBoltDB(conf)
		if err != nil {
			return nil, dberr(err)
		}
		reader = b
		nodes = b
		queue = b
		writers.Add(b)

	case "datastore":
		d, err := datastore.NewDatastore(conf.Server.Databases.Datastore)
		if err != nil {
			return nil, dberr(err)
		}
		reader = d
		writers.Add(d)

	case "dynamodb":
		d, err := dynamodb.NewDynamoDB(conf.Server.Databases.DynamoDB)
		if err != nil {
			return nil, dberr(err)
		}
		reader = d
		writers.Add(d)

	case "elastic":
		e, err := elastic.NewElastic(conf.Server.Databases.Elastic)
		if err != nil {
			return nil, dberr(err)
		}
		reader = e
		nodes = e
		queue = e
		writers.Add(e)

	case "mongodb":
		m, err := mongodb.NewMongoDB(conf.Server.Databases.MongoDB)
		if err != nil {
			return nil, dberr(err)
		}
		reader = m
		nodes = m
		queue = m
		writers.Add(m)
	}

	switch strings.ToLower(conf.Backend) {
	case "manual":
		if nodes == nil {
			return nil, fmt.Errorf(
				"cannot enable manual compute backend, database %s does not implement "+
					"the scheduler service", conf.Server.Database)
		}
		if queue == nil {
			return nil, fmt.Errorf(
				"cannot enable manual compute backend, database %s does not implement "+
					"a task queue", conf.Server.Database)
		}
		sched = &scheduler.Scheduler{
			Conf:  conf.Scheduler,
			Log:   log.Sub("scheduler"),
			Nodes: nodes,
			Queue: queue,
			Event: &writers,
		}

	case "aws-batch":
		b, err := batch.NewBackend(conf.Backends.Batch)
		if err != nil {
			return nil, err
		}
		writers.Add(b)

	case "gridengine":
		writers.Add(gridengine.NewBackend(conf))
	case "htcondor":
		writers.Add(htcondor.NewBackend(conf))
	case "local":
		writers.Add(local.NewBackend(conf, log.Sub("local"), workerCmd.Run))
	case "noop":
		writers.Add(noop.NewBackend(conf))
	case "pbs":
		writers.Add(pbs.NewBackend(conf))
	case "slurm":
		writers.Add(slurm.NewBackend(conf))
	case "", "none":
		// Do nothing.
	default:
		return nil, fmt.Errorf("unknown backend: '%s'", conf.Backend)
	}

	return &Server{
		Server: &server.Server{
			RPCAddress:       ":" + conf.Server.RPCPort,
			HTTPPort:         conf.Server.HTTPPort,
			Password:         conf.Server.Password,
			DisableHTTPCache: conf.Server.DisableHTTPCache,
			Log:              log,
			Tasks: &server.TaskService{
				Name:  conf.Server.ServiceName,
				Event: &writers,
				Read:  reader,
			},
			Events: &events.Service{&writers},
			Nodes:  nodes,
		},
		Scheduler: sched,
	}, nil
}

// Run runs a default Funnel server.
// This opens a database, and starts an API server, scheduler and task logger.
// This blocks indefinitely.
func (s *Server) Run(ctx context.Context) error {

	// Start server
	errch := make(chan error)
	go func() {
		errch <- s.Server.Serve(ctx)
	}()

	// Start Scheduler
	if s.Scheduler != nil {
		go func() {
			errch <- s.Scheduler.Run(ctx)
		}()
	}

	// Block until done.
	// Server and scheduler must be stopped via the context.
	return <-errch
}

func dberr(err error) error {
	return fmt.Errorf("error occurred while connecting to or creating the database: %v", err)
}
