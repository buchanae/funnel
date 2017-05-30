package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger/logutils"
	"github.com/ohsu-comp-bio/funnel/scheduler"
	"github.com/ohsu-comp-bio/funnel/worker"
	"github.com/spf13/cobra"
)

// Cmd represents the worker command
var Cmd = &cobra.Command{
	Use:                "worker",
	Aliases:            []string{"workers"},
	Short:              "Starts a Funnel worker.",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		var configFile string
		c := config.DefaultConfig()

		f := cmd.Flags()
		f.StringVarP(&configFile, "config", "c", "", "Config File")
		f.StringVar(&c.Worker.ID, "id", c.Worker.ID, "Worker ID")
		f.DurationVar(&c.Worker.Timeout, "timeout", c.Worker.Timeout, "Timeout in seconds")
		f.StringVar(&c.HostName, "hostname", c.HostName, "Host name or IP")
		f.StringVar(&c.RPCPort, "rpc-port", c.RPCPort, "RPC Port")
		f.StringVar(&c.WorkDir, "work-dir", c.WorkDir, "Working Directory")
		f.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Level of logging")
		f.StringVar(&c.LogPath, "log-path", c.LogPath, "File path to write logs to")

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

func init() {
}

// Run runs a worker with the given config, blocking until the worker exits.
func Run(conf config.Config) error {

	logutils.Configure(conf)

	if conf.Worker.ID == "" {
		conf.Worker.ID = scheduler.GenWorkerID("funnel")
	}

	w, err := worker.NewWorker(conf.Worker)
	if err != nil {
		return err
	}
	w.Run()
	return nil
}
