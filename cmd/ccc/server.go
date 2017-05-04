package ccc

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/logger/logutils"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/ccc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var log = logger.New("ccc")
var configFile string

func init() {
	flags := Cmd.Flags()
	flags.StringVarP(&configFile, "config", "c", "", "Config File")
}

// Cmd represents the `funnel server` CLI command set.
var Cmd = &cobra.Command{
	Use:   "ccc",
	Short: "Starts a Funnel CCC Proxy.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var conf = config.DefaultConfig()
		config.ParseFile(configFile, &conf)
		return Run(conf)
	},
}

func Run(conf config.Config) error {
	logutils.Configure(conf)

  proxy, err := ccc.NewTaskProxy(conf)
  if err != nil {
    return err
  }

	srv := server.Server{
    RPCAddress: conf.RPCAddress(),
    HTTPPort: conf.HTTPPort,
    TaskServiceServer: proxy,
		DisableHTTPCache:       conf.DisableHTTPCache,
		DialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
  }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv.Serve(ctx)
	return nil
}
