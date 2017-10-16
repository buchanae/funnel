package server

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/compute/gce"
	"github.com/ohsu-comp-bio/funnel/compute/gridengine"
	"github.com/ohsu-comp-bio/funnel/compute/htcondor"
	"github.com/ohsu-comp-bio/funnel/compute/local"
	"github.com/ohsu-comp-bio/funnel/compute/manual"
	"github.com/ohsu-comp-bio/funnel/compute/noop"
	"github.com/ohsu-comp-bio/funnel/compute/openstack"
	"github.com/ohsu-comp-bio/funnel/compute/pbs"
	"github.com/ohsu-comp-bio/funnel/compute/scheduler"
	"github.com/ohsu-comp-bio/funnel/compute/slurm"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/server/boltdb"
	"github.com/ohsu-comp-bio/funnel/server/dynamodb"
	"github.com/ohsu-comp-bio/funnel/server/elastic"
	"strings"
)

// Run runs the "server run" command.
func Run(ctx context.Context, conf config.Config) error {
	logger.Configure(conf.Server.Logger)

	s, err := NewServer(conf)
	if err != nil {
		return err
	}
	return s.Run(ctx)
}

// Server is a Funnel server + scheduler.
type Server struct {
	*server.Server
	*scheduler.Scheduler
	DB       server.Database
	SDB      scheduler.Database
	SBackend scheduler.Backend
}

// NewServer returns a new Funnel server + scheduler based on the given config.
func NewServer(conf config.Config) (*Server, error) {
	var backend compute.Backend
	var db server.Database
	var sdb scheduler.Database
	var sched *scheduler.Scheduler
	var sbackend scheduler.Backend
	var err error

	switch strings.ToLower(conf.Server.Database) {
	case "boltdb":
		bdb, berr := boltdb.NewBoltDB(conf)
		err = berr
		db = bdb
		sdb = bdb
	case "dynamodb":
		ddb, derr := dynamodb.NewDynamoDB(conf.Server.Databases.DynamoDB)
		err = derr
		db = ddb
	case "elastic":
		es, eserr := elastic.NewTES(conf.Server.Databases.Elastic)
		err = eserr
		db = es
		sdb = es
		if err == nil {
			err = es.Init(context.Background())
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to or creating the database: %v", err)
	}

	switch strings.ToLower(conf.Backend) {
	case "gce", "manual", "openstack", "gce-mock":
		if sdb == nil {
			return nil, fmt.Errorf("database doesn't satisfy the scheduler interface")
		}

		backend = scheduler.NewComputeBackend(sdb)

		switch strings.ToLower(conf.Backend) {
		case "gce":
			sbackend, err = gce.NewBackend(conf)
		case "gce-mock":
			sbackend, err = gce.NewMockBackend(conf)
		case "manual":
			sbackend, err = manual.NewBackend(conf)
		case "openstack":
			sbackend, err = openstack.NewBackend(conf)
		}
		if err != nil {
			return nil, fmt.Errorf("error occurred while setting up backend: %v", err)
		}

		sched = scheduler.NewScheduler(sdb, sbackend, conf.Scheduler)
	case "gridengine":
		backend = gridengine.NewBackend(conf)
	case "htcondor":
		backend = htcondor.NewBackend(conf)
	case "local":
		backend = local.NewBackend(conf)
	case "noop":
		backend = noop.NewBackend(conf)
	case "pbs":
		backend = pbs.NewBackend(conf)
	case "slurm":
		backend = slurm.NewBackend(conf)
	default:
		return nil, fmt.Errorf("unknown backend: '%s'", conf.Backend)
	}

	db.WithComputeBackend(backend)
	srv := server.DefaultServer(db, conf.Server)

	return &Server{srv, sched, db, sdb, sbackend}, nil
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
