package gcp

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/events"
	"golang.org/x/net/context"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

type StackdriverEventWriter struct {
	client *logging.Client
	logger *logging.Logger
}

func NewStackdriverEventWriter(ctx context.Context, name, project string) (*StackdriverEventWriter, error) {

	client, err := logging.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}

	logger := client.Logger(name, logging.CommonLabels(map[string]string{
    "funnel_task": "yes",
  }))

	return &StackdriverEventWriter{client, logger}, nil
}

func (s *StackdriverEventWriter) Close() error {
	return s.client.Close()
}

func (s *StackdriverEventWriter) Write(e *events.Event) error {
	s.logger.Log(logging.Entry{
		Payload: e,
    Operation: &logpb.LogEntryOperation{
      Id: fmt.Sprintf("%s-%d", e.Id, e.Attempt),
      Producer: "worker.funnel.ohsu.edu",
    },
	})
	return nil
}
