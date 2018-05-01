package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/ohsu-comp-bio/funnel/compute/batch"
	"github.com/ohsu-comp-bio/funnel/compute/gridengine"
	"github.com/ohsu-comp-bio/funnel/compute/htcondor"
	"github.com/ohsu-comp-bio/funnel/compute/local"
	"github.com/ohsu-comp-bio/funnel/compute/pbs"
	"github.com/ohsu-comp-bio/funnel/compute/builtin"
	"github.com/ohsu-comp-bio/funnel/compute/slurm"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/database/boltdb"
	"github.com/ohsu-comp-bio/funnel/database/datastore"
	"github.com/ohsu-comp-bio/funnel/database/dynamodb"
	"github.com/ohsu-comp-bio/funnel/database/elastic"
	"github.com/ohsu-comp-bio/funnel/database/mongodb"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/tes"
)

// Run runs the "server run" command.
func Run(ctx context.Context, conf config.Config, log *logger.Logger) error {
	s, err := NewServer(ctx, conf, log)
	if err != nil {
		return err
	}

	return s.Run(ctx)
}

// Database represents the base funnel database interface
type Database interface {
	tes.ReadOnlyServer
	events.Writer
	Init() error
}

// NewServer returns a new Funnel server + scheduler based on the given config.
func NewServer(ctx context.Context, conf config.Config, log *logger.Logger) (*server.Server, error) {
	log.Debug("NewServer", "config", conf)

	var database Database
	var reader tes.ReadOnlyServer
  var schedService builtin.SchedulerServiceServer

	writers := events.MultiWriter{}

	// Database
	switch strings.ToLower(conf.Database) {
	case "boltdb":
		b, err := boltdb.NewBoltDB(conf.BoltDB)
		if err != nil {
			return nil, dberr(err)
		}
		database = b
		reader = b
		writers = append(writers, b)

	case "datastore":
		d, err := datastore.NewDatastore(conf.Datastore)
		if err != nil {
			return nil, dberr(err)
		}
		database = d
		reader = d
		writers = append(writers, d)

	case "dynamodb":
		d, err := dynamodb.NewDynamoDB(conf.DynamoDB)
		if err != nil {
			return nil, dberr(err)
		}
		database = d
		reader = d
		writers = append(writers, d)

	case "elastic":
		e, err := elastic.NewElastic(conf.Elastic)
		if err != nil {
			return nil, dberr(err)
		}
		database = e
		reader = e
		writers = append(writers, e)

	case "mongodb":
		m, err := mongodb.NewMongoDB(conf.MongoDB)
		if err != nil {
			return nil, dberr(err)
		}
		database = m
		reader = m
		writers = append(writers, m)

	default:
		return nil, fmt.Errorf("unknown database: '%s'", conf.Database)
	}

	// Initialize the Database
	if err := database.Init(); err != nil {
		return nil, fmt.Errorf("error creating database resources: %v", err)
	}

	// Event writers
	var writer events.Writer
	var err error

	eventWriterSet := make(map[string]interface{})
	for _, w := range conf.EventWriters {
		eventWriterSet[strings.ToLower(w)] = nil
	}

	for e := range eventWriterSet {
		switch e {
		case strings.ToLower(conf.Database):
			continue
		case "log":
			continue
		case "boltdb":
			writer, err = boltdb.NewBoltDB(conf.BoltDB)
		case "dynamodb":
			writer, err = dynamodb.NewDynamoDB(conf.DynamoDB)
		case "elastic":
			writer, err = elastic.NewElastic(conf.Elastic)
		case "kafka":
			writer, err = events.NewKafkaWriter(ctx, conf.Kafka)
		case "pubsub":
			writer, err = events.NewPubSubWriter(ctx, conf.PubSub)
		case "mongodb":
			writer, err = mongodb.NewMongoDB(conf.MongoDB)
		default:
			return nil, fmt.Errorf("unknown event writer: '%s'", e)
		}
		if err != nil {
			return nil, fmt.Errorf("error occurred while initializing the %s event writer: %v", e, err)
		}
		if writer != nil {
			writers = append(writers, writer)
		}
	}

	writer = &events.SystemLogFilter{Writer: &writers, Level: conf.Logger.Level}

	// Compute
	var compute server.ComputeBackend
	switch strings.ToLower(conf.Compute) {
	case "builtin":
		ev := &events.ErrLogger{Writer: writer, Log: log.Sub("scheduler")}
    var sched *builtin.Scheduler
		sched, err = builtin.NewScheduler(conf.Scheduler, log.Sub("scheduler"), ev)
		if err != nil {
			return nil, err
		}
		compute = sched
    schedService = sched

	case "aws-batch":
		compute, err = batch.NewBackend(ctx, conf.AWSBatch, reader, writer)
		if err != nil {
			return nil, err
		}

	case "local":
		compute, err = local.NewBackend(ctx, conf, log.Sub("local"))
		if err != nil {
			return nil, err
		}

	case "gridengine":
		compute = gridengine.NewBackend(conf, reader, writer)
	case "htcondor":
		compute = htcondor.NewBackend(ctx, conf, reader, writer)
	case "noop":
		compute = server.NoopCompute{}
	case "pbs":
		compute = pbs.NewBackend(ctx, conf, reader, writer)
	case "slurm":
		compute = slurm.NewBackend(ctx, conf, reader, writer)
	default:
		return nil, fmt.Errorf("unknown compute backend: '%s'", conf.Compute)
	}

	writer = &events.ErrLogger{Writer: writer, Log: log}

	return &server.Server{
    RPCAddress:       ":" + conf.Server.RPCPort,
    HTTPPort:         conf.Server.HTTPPort,
    User:             conf.Server.User,
    Password:         conf.Server.Password,
    DisableHTTPCache: conf.Server.DisableHTTPCache,
    Log:              log,
    Tasks: &server.TaskService{
      Name:    conf.Server.ServiceName,
      Event:   writer,
      Compute: compute,
      Read:    reader,
      Log:     log,
    },
    Events: &events.Service{Writer: writer},
    Nodes:  schedService,
	}, nil
}

func dberr(err error) error {
	return fmt.Errorf("error occurred while connecting to or creating the database: %v", err)
}
