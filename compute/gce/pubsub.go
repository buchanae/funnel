package gce

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/events"
	"golang.org/x/net/context"
  "github.com/golang/protobuf/proto"
  "cloud.google.com/go/pubsub"
)

var projectID = "isb-cgc-04-0029"
var topicName = "funnel"

type PubSubBackend struct {
  client *pubsub.Client
  topic *pubsub.Topic
}

func NewPubSubBackend() (*PubSubBackend, error) {
  ctx := context.Background()

  client, err := pubsub.NewClient(ctx, projectID)
  if err != nil {
    return nil, err
  }

  topic := client.Topic(topicName)

  return &PubSubBackend{client, topic}, nil
}

func (p *PubSubBackend) Close() error {
  p.topic.Stop()
  p.client.Close()
  return nil
}

func (p *PubSubBackend) Submit(task *tes.Task) error {
  ctx := context.Background()

  b, err := proto.Marshal(task)
  if err != nil {
    return err
  }

  msg := &pubsub.Message{Data: b}
  res := p.topic.Publish(ctx, msg)

  _, gerr := res.Get(ctx)
  return gerr
}

func (p *PubSubBackend) Cancel(id string) error {
  ctx := context.Background()
  ev := events.NewState(id, 0, tes.Canceled)

  b, err := proto.Marshal(ev)
  if err != nil {
    return err
  }

  msg := &pubsub.Message{Data: b}
  res := p.topic.Publish(ctx, msg)

  _, gerr := res.Get(ctx)
  return gerr
}
