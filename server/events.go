package server

import (
  "bytes"
	"fmt"
	"github.com/boltdb/bolt"
	proto "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/events"
	"golang.org/x/net/context"
)

// CreateEvent creates an event for the server to handle.
func (taskBolt *TaskBolt) CreateEvent(ctx context.Context, req *events.Event) (*events.CreateEventResponse, error) {
  log.Info("Create event", req)

	err := taskBolt.db.Update(func(tx *bolt.Tx) error {
		if req.Type == events.Type_TASK_STATE {
			err := transitionTaskState(tx, req.Id, req.GetState())
			if err != nil {
				return err
			}
		}

		// Try to load existing task log
		id := fmt.Sprintf("%s-%d-%d", req.Id, req.Attempt, req.Timestamp)
		reqbytes, err := proto.Marshal(req)
		if err != nil {
			return err
		}
		return tx.Bucket(TaskEvents).Put([]byte(id), reqbytes)
	})
	return &events.CreateEventResponse{}, err
}

func (taskBolt *TaskBolt) GetEvents(ctx context.Context, req *events.GetEventsRequest) (*events.GetEventsResponse, error) {

  resp := events.GetEventsResponse{}
  taskBolt.db.View(func(tx *bolt.Tx) error {

    prefix := []byte(req.Id)
    c := tx.Bucket(TaskEvents).Cursor()

    for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {

      ev := &events.Event{}
      err := proto.Unmarshal(v, ev)
      if err != nil {
        continue
      }
      resp.Events = append(resp.Events, ev)
    }
    return nil
  })
  return &resp, nil
}
