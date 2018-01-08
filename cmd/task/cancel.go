package task

import (
  "context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
)

type CancelOpts struct {
  TaskOpts
  IDs []string `args`
}

func DefaultCancelOpts() CancelOpts {
  return CancelOpts{TaskOpts: DefaultTaskOpts()}
}

// Cancel tasks by ID.
func Cancel(ctx context.Context, opts CancelOpts) error {
//func Cancel(server string, ids []string, writer io.Writer) error {
  if len(opts.IDs) == 0 {
    return fmt.Errorf("zero task IDs given")
  }

	cli, err := client.NewClient(opts.Server)
	if err != nil {
		return err
	}

	res := []string{}

	for _, taskID := range opts.IDs {
		resp, err := cli.CancelTask(ctx, &tes.CancelTaskRequest{Id: taskID})
		if err != nil {
			return err
		}
		// CancelTaskResponse is an empty struct
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
