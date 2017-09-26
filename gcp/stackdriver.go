package gcp

import (
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/events"
	"golang.org/x/net/context"
  "github.com/jkawamoto/structpbconv"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
  "google.golang.org/api/iterator"
  structpb "github.com/golang/protobuf/ptypes/struct"
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



type StackdriverEventReader struct {
	client *logadmin.Client
  it     *logadmin.EntryIterator
}

func NewStackdriverEventReader(ctx context.Context, project string) (*StackdriverEventReader, error) {

	client, err := logadmin.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
  it := client.Entries(ctx, logadmin.Filter(`resource.type="global" labels."funnel_task"="yes"`))
	return &StackdriverEventReader{client, it}, nil
}

func (s *StackdriverEventReader) Close() error {
	return s.client.Close()
}

func (s *StackdriverEventReader) WriteTo(w events.Writer) error {
  for {
    entry, err := s.it.Next()
    if err == iterator.Done {
      return nil
    }
    if err != nil {
      return err
    }
    ev := &events.Event{}
    if x, ok := entry.Payload.(*structpb.Struct); ok {
      if err := structpbconv.Convert(x, ev); err == nil {
        fmt.Println(ev)
      }
    }

  }
	return nil
}
