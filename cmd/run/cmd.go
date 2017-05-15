package run

import (
	"bufio"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strings"
)


// *********************************************************************
// IMPORTANT:
// Usage/help docs are defined in usage.go.
// If you're updating flags, you probably need to update that file.
// *********************************************************************


var log = logger.New("run")
var vars cmdvars

// Cmd represents the run command
var Cmd = &cobra.Command{
	Use:   "run 'CMD' [flags]",
	Short: "Run a task.",
	RunE:  runE,
}

func init() {
	f := Cmd.Flags()
  addCmdFlags(f, &vars)
	Cmd.SetUsageTemplate(usage)
}

// runE is the main cobra CLI command handlers.
func runE(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("you must specify a command to run")
	}

	if len(args) > 1 {
		cmd.Usage()
		return fmt.Errorf("--in, --out and --env args should have the form 'KEY=VALUE' not 'KEY VALUE'. Extra args: %s", args[1:])
	}

  // Load CLI arguments from files, which allows reusing common CLI args.
	for _, xf := range extraFiles {
		b, _ := ioutil.ReadFile(xf)
		extra = append(extra, string(b))
	}

  // Load CLI arguments from stdin, which allows bash heredoc for easily
  // spreading args over multiple lines.
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, _ := ioutil.ReadAll(os.Stdin)
		if len(b) > 0 {
			extra = append(extra, string(b))
		}
	}

  // Load and parse all "extra" CLI arguments.
	for _, ex := range extra {
		sp, _ := shellquote.Split(ex)
		cmd.ParseFlags(sp)
	}

	task, err := valsToTask(executorCmd, vals)
	if err != nil {
		cmd.Usage()
		return err
	}

  // Split command string based on shell syntax.
	executorCmd, _ := shellquote.Split(args[0])

  // TES HTTP client
	cli := client.NewClient(server)

	if vals.container == "" {
		return nil, fmt.Errorf("you must specify a container")
	}

	if len(scatterFiles) > 0 {
		tg := taskGroup{}
		for _, sc := range scatterFiles {
			fh, _ := os.Open(sc)
			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				// Copy task vals to act as base task
				tv := vals
				// Per-scatter flags
				fl := &pflag.FlagSet{}
				addTaskFlags(fl, &tv)
				sp, _ := shellquote.Split(scanner.Text())
				fl.Parse(sp)
				t, _ := valsToTask(executorCmd, tv)
				tg.runTask(t, cli, tv.wait, tv.waitFor)
			}
		}
		return tg.wait()
	}

	return runTask(task, cli, vals.wait, vals.waitFor)
}


func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
