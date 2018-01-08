package task

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/client"
)

type WaitOpts struct {
  TaskOpts
  IDs []string `args`
}

func DefaultWaitOpts() WaitOpts {
  return WaitOpts{TaskOpts: DefaultTaskOpts()}
}

// Wait for a task to finish.
func Wait(ctx context.Context, opts WaitOpts) error {
//func Wait(server string, ids []string) error {
  if len(opts.IDs) == 0 {
    return fmt.Errorf("zero task IDs given")
  }

	cli, err := client.NewClient(opts.Server)
	if err != nil {
		return err
	}

	return cli.WaitForTask(ctx, opts.IDs...)
}
