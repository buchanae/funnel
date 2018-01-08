package task

import (
	"fmt"
  "context"
	"github.com/ohsu-comp-bio/funnel/cmd/util"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type Task struct {
  util.Command
  Server string
}

var DefaultTask = Task{
  Command: util.Command{
    Names: []string{"task", "tasks"},
		Short:   "Make API calls to a TES server.",
  },
  Server: "http://localhost:8000",
}


// NewCommand returns the "task" subcommands.
func NewCommand() *cobra.Command {

	get := &cobra.Command{
		Use:   "get [taskID ...]",
		Short: "Get one or more tasks by ID.",
    Run: DefaultGet(),
	}

	cancel := &cobra.Command{
		Use:   "cancel [taskID ...]",
		Short: "Cancel one or more tasks by ID.",
    Run: wrap(DefaultCancel()),
	}

	wait := &cobra.Command{
		Use:   "wait [taskID...]",
		Short: "Wait for one or more tasks to complete.",
    Run: DefaultWait(),
	}

	cmd.AddCommand(create, get, list, cancel, wait)
	return cmd
}

type TaskView tes.TaskView

func (tv *TaskView) String() string {
  return tes.TaskView(*tv).String()
}
func (tv *TaskView) Type() string {
  return "string"
}
func (tv *TaskView) Set(raw string) error {
  raw = strings.ToUpper(raw)
	view, ok := tes.TaskView_value[raw]
	if !ok {
		return fmt.Errorf("Unknown task view: %s. Valid task views: ['basic', 'minimal', 'full']", raw)
	}

  *tv = TaskView(view)
  return nil
}
