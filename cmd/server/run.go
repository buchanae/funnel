package server

import (
	"context"
	"github.com/imdario/mergo"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/gcp"
	"github.com/ohsu-comp-bio/funnel/compute/gce"
	"github.com/spf13/cobra"
)

var log = logger.New("server run cmd")

// runCmd represents the `funnel server run` command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Funnel server.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		// parse config file
		conf := config.DefaultConfig()
		config.ParseFile(configFile, &conf)

		// make sure server address and password is inherited by scheduler nodes and workers
		conf = config.InheritServerProperties(conf)
		flagConf = config.InheritServerProperties(flagConf)

		// file vals <- cli val
		err := mergo.MergeWithOverwrite(&conf, flagConf)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		return Run(ctx, conf)
	},
}

// Run runs a default Funnel server.
// This opens a database, and starts an API server, scheduler and task logger.
func Run(ctx context.Context, conf config.Config) error {
	logger.Configure(conf.Server.Logger)

  /*
	db, err = server.NewTaskBolt(conf)
	if err != nil {
		log.Error("Couldn't open database", err)
		return err
	}
  */
  pubsub, err := gce.NewPubSubBackend()
  if err != nil {
    return err
  }

  db, err := gcp.NewDatastoreTES("isb-cgc-04-0029", pubsub)
  if err != nil {
    return err
  }

	srv := server.DefaultServer(conf.Server)
  srv.TaskServiceServer = db

  /*
  backend, err = gce.NewPubSubBackend()
  if err != nil {
    return err
  }
  */

	// Block

	// Start server
	errch := make(chan error)
	go func() {
		errch <- srv.Serve(ctx)
	}()

	// Block until done.
	// Server and scheduler must be stopped via the context.
	return <-errch
}
