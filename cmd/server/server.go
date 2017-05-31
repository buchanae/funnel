package server

import (
	"context"

	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/logger/logutils"
	"github.com/ohsu-comp-bio/funnel/scheduler"
	"github.com/ohsu-comp-bio/funnel/scheduler/condor"
	"github.com/ohsu-comp-bio/funnel/scheduler/gce"
	"github.com/ohsu-comp-bio/funnel/scheduler/local"
	"github.com/ohsu-comp-bio/funnel/scheduler/manual"
	"github.com/ohsu-comp-bio/funnel/scheduler/openstack"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/server/badger"
	"github.com/spf13/cobra"
)

var log = logger.New("server cmd")

// Cmd represents the `funnel server` CLI command set.
var Cmd = &cobra.Command{
	Use:                "server",
	Short:              "Starts a Funnel server.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		var configFile string
		c := config.DefaultConfig()

		f := cmd.Flags()
		f.StringVarP(&configFile, "config", "c", "", "Config File")
		f.StringVar(&c.HostName, "hostname", c.HostName, "Host name or IP")
		f.StringVar(&c.RPCPort, "rpc-port", c.RPCPort, "RPC Port")
		f.StringVar(&c.WorkDir, "work-dir", c.WorkDir, "Working Directory")
		f.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Level of logging")
		f.StringVar(&c.LogPath, "log-path", c.LogPath, "File path to write logs to")
		f.StringVar(&c.HTTPPort, "http-port", c.HTTPPort, "HTTP Port")
		f.StringVar(&c.DBPath, "db-path", c.DBPath, "Database path")
		f.StringVar(&c.Scheduler, "scheduler", c.Scheduler, "Name of scheduler to enable")

		// Parse flags only to get config file.
		if err := f.Parse(args); err != nil {
			return err
		}

		// Load the config file.
		config.ParseFile(configFile, &c)

		// Parse flags again to overwrite default config values.
		if err := f.Parse(args); err != nil {
			return err
		}

		// The worker config inherits a few parts of the root config.
		c = config.WorkerInheritConfigVals(c)
		return Run(c)
	},
}

// Run runs a default Funnel server.
// This opens a database, and starts an API server and scheduler.
// This blocks indefinitely.
func Run(conf config.Config) error {
  conf.LogLevel = "debug"
	logutils.Configure(conf)

	// make sure the proper defaults are set
	conf = config.WorkerInheritConfigVals(conf)

	db, err := badger.NewTaskBadger(conf)
	if err != nil {
		log.Error("Couldn't open database", err)
		return err
	}
  /*
	db, err := server.NewTaskBolt(conf)
	if err != nil {
		log.Error("Couldn't open database", err)
		return err
	}
  */

	srv := server.DefaultServer(db, conf)

	sched, err := scheduler.NewScheduler(db, conf)
	if err != nil {
		return err
	}

	sched.AddBackend(gce.Plugin)
	sched.AddBackend(condor.Plugin)
	sched.AddBackend(openstack.Plugin)
	sched.AddBackend(local.Plugin)
	sched.AddBackend(manual.Plugin)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	var srverr error
	go func() {
		srverr = srv.Serve(ctx)
		cancel()
	}()

	// Start scheduler
	err = sched.Start(ctx)
	if err != nil {
		return err
	}

	// Block
	<-ctx.Done()
	if srverr != nil {
		log.Error("Server error", srverr)
	}
	return srverr
}
