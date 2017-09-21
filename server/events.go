package server

import (
	"fmt"
	"github.com/boltdb/bolt"
	proto "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/events"
	"golang.org/x/net/context"
)

// CreateEvent creates an event for the server to handle.
func (taskBolt *TaskBolt) CreateEvent(ctx context.Context, req *events.Event) (*events.CreateEventResponse, error) {

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
