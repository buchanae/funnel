package task

import (
	"github.com/ohsu-comp-bio/funnel/client"
	"golang.org/x/net/context"
)

// Wait runs the "task wait" CLI command, which polls the server,
// calling GetTask() for each ID, and blocking until the tasks have
// reached a terminal state.
func Wait(cli *client.Client, ids []string) error {
	return cli.WaitForTask(context.Background(), ids...)
}
