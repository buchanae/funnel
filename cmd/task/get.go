package task

import (
  "context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
)

type GetOpts struct {
  Server string
  View TaskView
  IDs []string `args`
  Out io.Writer
}

func DefaultGetOpts() GetOpts {
  return GetOpts{
    Server: "http://localhost:8000",
    View: TaskView(tes.Full),
    Out: os.Stdout,
  }
}

// Get tasks by ID.
func Get(ctx context.Context, opts GetOpts) error {
//func Get(server string, ids []string, taskView string, w io.Writer) error {
  if len(opts.IDs) == 0 {
    return fmt.Errorf("zero task IDs given")
  }

	cli, err := client.NewClient(opts.Server)
	if err != nil {
		return err
	}

	res := []string{}

	for _, taskID := range opts.IDs {
		resp, err := cli.GetTask(ctx, &tes.GetTaskRequest{
			Id:   taskID,
			View: tes.TaskView(opts.View),
		})
		if err != nil {
			return err
		}

		out, err := cli.Marshaler.MarshalToString(resp)
		if err != nil {
			return err
		}
		res = append(res, out)
	}

	for _, x := range res {
		fmt.Fprintln(opts.Out, x)
	}
	return nil
}
