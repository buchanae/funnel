package worker

import (
  "fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
)

// configureWriters creates multiple event writers based on the given config.
func configureWriters(conf config.Worker) (events.Writer, error) {
	var writers []events.Writer
	for _, w := range conf.ActiveEventWriters {

		var writer events.Writer
		var err error

		switch w {
		case "dynamodb":
			writer, err = events.NewDynamoDBEventWriter(conf.EventWriters.DynamoDB)
		case "log":
			writer = events.NewLogger("worker")
		case "rpc":
      c := conf.EventWriters.RPC
			writer, err = events.NewRPCWriter(c.ServerAddress, c.ServerPassword, c.UpdateTimeout)
		default:
			err = fmt.Errorf("unknown EventWriter")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate EventWriter: %v", err)
		}

		writers = append(writers, writer)
	}
	if writers == nil {
		return events.Discard, nil
	}
	return events.MultiWriter(writers...), nil
}
