package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/jsonpb"
	cmdutil "github.com/ohsu-comp-bio/funnel/cmd/util"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/database/mongodb"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
)

var (
	configFile string
	conf       config.Config
	flagConf   config.Config
)

var Cmd = &cobra.Command{
	Use:   "events",
	Short: "Access task event streams.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		conf, err = cmdutil.MergeConfigFileWithFlags(configFile, flagConf)
		if err != nil {
			return fmt.Errorf("error processing config: %v", err)
		}

		return nil
	},
}

func init() {
	Cmd.SetGlobalNormalizationFunc(cmdutil.NormalizeFlags)
	f := Cmd.PersistentFlags()
	f.StringVarP(&configFile, "config", "c", configFile, "Config file")
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read an event stream.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		w := &events.JSONWriter{Output: os.Stdout}
		ctx := context.Background()
		_, err := events.NewKafkaReader(ctx, conf.Kafka, events.KafkaOffsetOldest, w)
		if err != nil {
			return err
		}

		block := make(chan struct{})
		<-block

		return nil
	},
}
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write an event stream.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		log := logger.NewLogger("events", conf.Logger)
		ctx := context.Background()
		db, err := mongodb.NewMongoDB(conf.MongoDB)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(os.Stdin)

		for {
			ev := &events.Event{}
			err := jsonpb.UnmarshalNext(dec, ev)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			err = db.WriteEvent(ctx, ev)
			if err != nil {
				log.Error("error writing event", "error", err)
				continue
			}
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(readCmd)
	Cmd.AddCommand(writeCmd)
}
