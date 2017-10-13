package worker

import (
	"fmt"
	cmdutil "github.com/ohsu-comp-bio/funnel/cmd/util"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/server/elastic"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/worker"
	"github.com/spf13/cobra"
	"path"
)

// NewCommand returns the worker command
func NewCommand() *cobra.Command {
	cmd, _ := newCommandHooks()
	return cmd
}

type hooks struct {
	Run func(conf config.Worker, taskID string) error
}

func newCommandHooks() (*cobra.Command, *hooks) {
	hooks := &hooks{
		Run: Run,
	}

	var (
		configFile    string
		conf          config.Config
		flagConf      config.Config
		serverAddress string
		taskID        string
	)

	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Funnel worker commands.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			flagConf, err = cmdutil.ParseServerAddressFlag(serverAddress, flagConf)
			if err != nil {
				return fmt.Errorf("error parsing the server address: %v", err)
			}

			conf, err = cmdutil.MergeConfigFileWithFlags(configFile, flagConf)
			if err != nil {
				return fmt.Errorf("error processing config: %v", err)
			}

			return nil
		},
	}
	f := cmd.PersistentFlags()
	f.StringVarP(&configFile, "config", "c", "", "Config File")
	f.StringVar(&serverAddress, "server-address", "", "RPC address of Funnel server")
	f.StringVar(&flagConf.Worker.WorkDir, "work-dir", flagConf.Worker.WorkDir, "Working Directory")
	f.StringVar(&flagConf.Worker.Logger.Level, "log-level", flagConf.Worker.Logger.Level, "Level of logging")
	f.StringVar(&flagConf.Worker.Logger.OutputFile, "log-path", flagConf.Worker.Logger.OutputFile, "File path to write logs to")

	run := &cobra.Command{
		Use:   "run",
		Short: "Run a task directly, bypassing the server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if taskID == "" {
				return fmt.Errorf("no taskID was provided")
			}

			return hooks.Run(conf.Worker, taskID)
		},
	}
	f = run.Flags()
	f.StringVar(&taskID, "task-id", "", "Task ID")

	cmd.AddCommand(run)

	return cmd, hooks
}

// NewDefaultWorker returns a new configured DefaultWorker instance.
func NewDefaultWorker(conf config.Worker, taskID string) (worker.Worker, error) {
	var err error
	var reader worker.TaskReader
	var writer events.Writer

	// Map files into this baseDir
	baseDir := path.Join(conf.WorkDir, taskID)

	err = util.EnsureDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create worker baseDir: %v", err)
	}

	switch conf.TaskReader {
	case "rpc":
		reader, err = worker.NewRPCTaskReader(conf, taskID)
	case "dynamodb":
		reader, err = worker.NewDynamoDBTaskReader(conf.TaskReaders.DynamoDB, taskID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate TaskReader: %v", err)
	}

	writers := []events.Writer{}
	for _, w := range conf.ActiveEventWriters {
		switch w {
		case "dynamodb":
			writer, err = events.NewDynamoDBEventWriter(conf.EventWriters.DynamoDB)
		case "log":
			writer = events.NewLogger("worker")
		case "rpc":
			writer, err = events.NewRPCWriter(conf)
		case "elastic":
			writer, err = elastic.NewElastic(conf.EventWriters.Elastic)
		default:
			err = fmt.Errorf("unknown EventWriter")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate EventWriter: %v", err)
		}
		writers = append(writers, writer)
	}

	return &worker.DefaultWorker{
		Conf:       conf,
		Mapper:     worker.NewFileMapper(baseDir),
		Store:      storage.Storage{},
		TaskReader: reader,
		Event:      events.NewTaskWriter(taskID, 0, conf.Logger.Level, events.MultiWriter(writers...)),
	}, nil
}
