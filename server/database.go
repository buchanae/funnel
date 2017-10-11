package server

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/events"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Database represents the interface to the database used by the scheduler, scaler, etc.
// Mostly, this exists so it can be mocked during testing.
type Database interface {
	tes.TaskServiceServer
	events.EventServiceServer
	pbs.SchedulerServiceServer
	WithComputeBackend(compute.Backend)
}

func DatabaseFromConfig(conf config.Server) (db Database, err error) {
	switch strings.ToLower(conf.Server.Database) {
	case "boltdb":
		db, err = boltdb.NewBoltDB(conf)
	case "dynamodb":
		db, err = dynamodb.NewDynamoDB(conf.Server.Databases.DynamoDB)
  default:
    err = fmt.Errorf("unknown database: %s", conf.Server.Database)
	}
	if err != nil {
    err = fmt.Errorf("error occurred while connecting to or creating the database: %v", err)
	}
  return
}
