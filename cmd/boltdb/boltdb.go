package boltdb

import (
  "context"
  "fmt"
	"github.com/spf13/cobra"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/ohsu-comp-bio/funnel/config"
)

// Cmd represents the 'funnel gce" CLI command set.
var Cmd = &cobra.Command{
	Use: "boltdb",
}

func init() {
  Cmd.AddCommand(getEventsCmd)
}

var getEventsCmd = &cobra.Command{
	Use: "get-events",
	RunE: func(cmd *cobra.Command, args []string) error {
    if len(args) != 1 {
      cmd.Usage()
      return nil
    }

    conf := config.Config{
      Server: config.Server{
        DBPath: args[0],
      },
    }
    db, err := server.NewTaskBolt(conf)
    if err != nil {
      return err
    }

    ctx := context.Background()
    resp, err := db.GetEvents(ctx, &events.GetEventsRequest{})
    if err != nil {
      return err
    }

    for _, ev := range resp.Events {
      fmt.Println(ev)
    }

		return nil
	},
}
