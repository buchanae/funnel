package gcp

import (
  "context"
	"github.com/spf13/cobra"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/gcp"
)

var Cmd = &cobra.Command{
  Use: "gcp",
}

func init() {
  Cmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
  Use: "server",
  RunE: func(cmd *cobra.Command, args []string) error {
    conf := config.DefaultConfig()
    srv := server.NewServer(conf)
    srv.TaskServiceServer = &gcp.TaskServiceServer{}
    ctx := context.Background()
    return srv.Serve(ctx)
  },
}
