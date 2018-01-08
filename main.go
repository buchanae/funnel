package main

import (
	//"github.com/ohsu-comp-bio/funnel/cmd/aws"
	"github.com/ohsu-comp-bio/funnel/cmd/examples"
	"github.com/ohsu-comp-bio/funnel/cmd/gce"
	"github.com/ohsu-comp-bio/funnel/cmd/node"
	"github.com/ohsu-comp-bio/funnel/cmd/run"
	"github.com/ohsu-comp-bio/funnel/cmd/server"
	"github.com/ohsu-comp-bio/funnel/cmd/task"
	"github.com/ohsu-comp-bio/funnel/cmd/termdash"
	"github.com/ohsu-comp-bio/funnel/cmd/version"
	"github.com/ohsu-comp-bio/funnel/cmd/worker"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var rootcmd = &cobra.Command{
	Use:           "funnel",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var genMarkdownCmd = &cobra.Command{
	Use:    "genmarkdown",
	Short:  "generate markdown formatted documentation for the funnel commands",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return doc.GenMarkdownTree(RootCmd, "./funnel-cmd-docs")
	},
}

var genBashCompletionCmd = &cobra.Command{
	Use:    "genbash",
	Short:  "generate bash completions for the funnel commands",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

var cmds = []util.Interface{
  task.Task,
  task.Create,
  task.List,
  task.Get,
  task.Cancel,
  task.Wait,
  examples.Examples,
  gce.Run,
  termdash.Termdash,
  node.Run,
  server.Run,
  worker.Run,
  version.Version,
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logger.PrintSimpleError(err)
		os.Exit(-1)
	}
}
