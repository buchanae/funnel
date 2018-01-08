package task

import (
  "context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"os"
)

var DefaultCreate = Create{
  Task: DefaultTask,
  Command: util.Command{
    Names: []string{
      "task create",
      "create tasks",
    },
		Short: "Create one or more tasks to run on the server.",
  },
}

// Create tasks.
type Create struct {
  Task
  util.Command
}

// Run the command.
func (c *Create) Run(ctx context.Context) error {
  files := c.Args

  if len(files) == 0 {
    return fmt.Errorf("zero task files given")
  }

	cli, err := client.NewClient(c.Server)
	if err != nil {
		return err
	}

	res := []string{}

	for _, taskFile := range files {
		var err error
		var task tes.Task

		f, err := os.Open(taskFile)
		defer f.Close()
		if err != nil {
			return err
		}

		err = jsonpb.Unmarshal(f, &task)
		if err != nil {
			return fmt.Errorf("can't load task: %s", err)
		}

		r, err := cli.CreateTask(ctx, &task)
		if err != nil {
			return err
		}
		res = append(res, r.Id)
	}

	for _, x := range res {
		fmt.Fprintln(c.Out, x)
	}

	return nil
}
