package run

import (
	"bufio"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
)

// *********************************************************************
// IMPORTANT:
// Usage/help docs are defined in usage.go.
// If you're updating flags, you probably need to update that file.
// *********************************************************************

var log = logger.New("run")
var vals flagVals

// Cmd represents the run command
var Cmd = &cobra.Command{
	Use:   "run 'CMD' [flags]",
	Short: "Run a task.",
	RunE:  runE,
}

func init() {
	f := Cmd.Flags()
	addTopLevelFlags(f, &vals)
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
	for _, xf := range vals.extraFiles {
		b, _ := ioutil.ReadFile(xf)
		vals.extra = append(vals.extra, string(b))
	}

	// Load CLI arguments from stdin, which allows bash heredoc for easily
	// spreading args over multiple lines.
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, _ := ioutil.ReadAll(os.Stdin)
		if len(b) > 0 {
			vals.extra = append(vals.extra, string(b))
		}
	}

	// Load and parse all "extra" CLI arguments.
	for _, ex := range vals.extra {
		sp, _ := shellquote.Split(ex)
		cmd.ParseFlags(sp)
	}

	// Fill in empty values with defaults.
	defaultVals(&vals)

	// Split command string based on shell syntax.
	vals.cmd, _ = shellquote.Split(args[0])

	// TES HTTP client
	tg := taskGroup{
		printTask: vals.printTask,
		client:    client.NewClient(vals.server),
	}

	// Scatter all vals into tasks
	for _, v := range append([]flagVals{vals}, scatter(vals)...) {
		// Parse inputs, outputs, environ, and tags from flagVals
		// and update the task.
		task, err := valsToTask(v)
		if err != nil {
			cmd.Usage()
			return err
		}
		tg.runTask(task, v.wait, v.waitFor)
	}

	return tg.wait()
}

// scatter reads each line from each scatter file, extending "base" flagVals
// with per-scatter vals from each line.
func scatter(base flagVals) []flagVals {
	out := []flagVals{}

	for _, sc := range base.scatterFiles {
		// Read each line of the scatter file.
		fh, _ := os.Open(sc)
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			tv := base
			// Per-scatter flags
			fl := &pflag.FlagSet{}
			addTaskFlags(fl, &tv)
			sp, _ := shellquote.Split(scanner.Text())
			fl.Parse(sp)
			// Parse scatter file flags into new flagVals
			out = append(out, tv)
		}
	}
	return out
}
