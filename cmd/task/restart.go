package task

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart [taskID ...]",
	Short: "Restart one or more tasks by ID.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		res, err := doRestart(tesServer, args)
		if err != nil {
			return err
		}

		for _, x := range res {
			fmt.Println(x)
		}
		return nil
	},
}

func doRestart(server string, ids []string) ([]string, error) {
	cli := client.NewClient(server)
	res := []string{}

	for _, taskID := range ids {
		resp, err := cli.RestartTask(taskID)
		if err != nil {
			return nil, err
		}
		// RestartTaskResponse is an empty struct
		out, err := cli.Marshaler.MarshalToString(resp)
		if err != nil {
			return nil, err
		}
		res = append(res, out)
	}
	return res, nil
}
