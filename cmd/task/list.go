package task

import (
  "context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
)

var DefaultList = List{
  Task: DefaultTask,
  Command: util.Command{
		Names: []string{
      "list tasks",
      "task list",
      "tasks list",
    },
		Short: "List all tasks.",
  },
  View: TaskView(tes.Basic),
}

// List tasks.
type List struct {
  Task
  util.Command

  // Minimal, basic, of full task view.
  View TaskView
  // Start at this page.
  PageToken string
  PageSize uint32
  // Return all pages.
  All bool
}

func (l *List) Run(ctx context.Context) error {
//func List(server, taskView, pageToken string, pageSize uint32, all bool, writer io.Writer) error {
	cli, err := client.NewClient(l.Server)
	if err != nil {
		return err
	}

	output := &tes.ListTasksResponse{}
  pageToken := l.PageToken

	for {
		resp, err := cli.ListTasks(ctx, &tes.ListTasksRequest{
			View:      tes.TaskView(l.View),
			PageToken: pageToken,
			PageSize:  l.PageSize,
		})
		if err != nil {
			return err
		}

		output.Tasks = append(output.Tasks, resp.Tasks...)
		output.NextPageToken = resp.NextPageToken
		pageToken = resp.NextPageToken

		if !l.All || (l.All && pageToken == "") {
			break
		}
	}

	response, err := cli.Marshaler.MarshalToString(output)
	if err != nil {
		return fmt.Errorf("marshaling error: %v", err)
	}

	fmt.Fprintf(l.Out, "%s\n", response)
	return nil
}
