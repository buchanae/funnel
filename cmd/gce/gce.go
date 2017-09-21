package gce

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/cmd/node"
	"github.com/ohsu-comp-bio/funnel/cmd/server"
	"github.com/ohsu-comp-bio/funnel/compute/gce"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
)

var log = logger.New("gce cmd")

// Cmd represents the 'funnel gce" CLI command set.
var Cmd = &cobra.Command{
	Use: "gce",
}

func init() {
	Cmd.AddCommand(nodeCmd)
	Cmd.AddCommand(serverCmd)
  Cmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.DefaultConfig()

		// Check that this is a GCE VM environment.
		// If not, fail.
		meta, merr := gce.LoadMetadata()
		if merr != nil {
			log.Error("Error getting GCE metadata", merr)
			return fmt.Errorf("can't find GCE metadata. This command requires a GCE environment")
		}

		var err error
		conf, err = gce.WithMetadataConfig(conf, meta)
		if err != nil {
			return err
		}
    fmt.Println(string(conf.ToYaml()))

    return nil
	},
}

var nodeCmd = &cobra.Command{
	Use: "node",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.DefaultConfig()

		// Check that this is a GCE VM environment.
		// If not, fail.
		meta, merr := gce.LoadMetadata()
		if merr != nil {
			log.Error("Error getting GCE metadata", merr)
			return fmt.Errorf("can't find GCE metadata. This command requires a GCE environment")
		}

		log.Info("Loaded GCE metadata")
		log.Debug("GCE metadata", meta)

		var err error
		conf, err = gce.WithMetadataConfig(conf, meta)
		if err != nil {
			return err
		}

    logger.Configure(conf.Scheduler.Node.Logger)
    return node.Run(conf)
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.DefaultConfig()

		// Check that this is a GCE VM environment.
		// If not, fail.
		meta, merr := gce.LoadMetadata()
		if merr != nil {
			log.Error("Error getting GCE metadata", merr)
			return fmt.Errorf("can't find GCE metadata. This command requires a GCE environment")
		}

		log.Info("Loaded GCE metadata")
		log.Debug("GCE metadata", meta)

		var err error
		conf, err = gce.WithMetadataConfig(conf, meta)
		if err != nil {
			return err
		}

		logger.Configure(conf.Server.Logger)
		return server.Run(context.Background(), conf)
	},
}
