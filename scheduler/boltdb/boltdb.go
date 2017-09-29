package boltdb

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
  "github.com/golang/protobuf/proto"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

var log = logger.Sub("boltsched")

var errNotFound = errors.New("not found")

// TasksQueued defines the name of a bucket which maps
// task ID -> nil
var TasksQueued = []byte("tasks-queued")

// Nodes maps:
// node ID -> pbs.Node struct
var Nodes = []byte("nodes")


type TaskQueue struct {
}
func (tq *TaskQueue) Queue(queueID, taskID string) error {
}
func (tq *TaskQueue) Read(queueID string, chunk int) error {
}
func (tq *TaskQueue) Remove(queueID, taskID string) error {
}


type SchedulerDatabase struct {
  db *bolt.DB
  conf config.Config
  tasks tes.TaskServiceServer
}

func NewSchedulerDatabase(conf config.Config, tasks tes.TaskServiceServer) (*SchedulerDatabase, error) {

  err := util.EnsurePath(conf.Server.DBPath)
  if err != nil {
    return nil, err
  }

  db, err := bolt.Open(conf.Server.DBPath, 0600, &bolt.Options{
    Timeout: time.Second * 5,
  })
  if err != nil {
    return nil, err
  }

  // Check to make sure all the required buckets have been created
  err = db.Update(func(tx *bolt.Tx) error {
    var err error
    if tx.Bucket(TasksQueued) == nil && err == nil{
      _, err = tx.CreateBucket(TasksQueued)
    }
    if tx.Bucket(Nodes) == nil && err == nil{
      _, err = tx.CreateBucket(Nodes)
    }
    return err
  })
  if err != nil {
    return nil, err
  }

  return &SchedulerDatabase{db, conf, tasks}, nil
}

// QueueTask adds a task to the scheduler queue.
func (sd *SchedulerDatabase) QueueTask(task *tes.Task) error {
	taskID := task.Id
	idBytes := []byte(taskID)

  return sd.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(TasksQueued).Put(idBytes, []byte{})
	})
}

// ReadQueue returns a slice of queued Tasks. Up to "n" tasks are returned.
func (sd *SchedulerDatabase) ReadQueue(n int) []*tes.Task {
	tasks := make([]*tes.Task, 0)
	sd.db.View(func(tx *bolt.Tx) error {

		// Iterate over the TasksQueued bucket, reading the first `n` tasks
		c := tx.Bucket(TasksQueued).Cursor()
		for k, _ := c.First(); k != nil && len(tasks) < n; k, _ = c.Next() {
      ctx := context.Background()
			task, err := sd.tasks.GetTask(ctx, &tes.GetTaskRequest{
        Id: string(k),
        View: tes.TaskView_FULL,
      })
      if err != nil {
        return err
      }
			tasks = append(tasks, task)
		}
		return nil
	})
	return tasks
}

// AssignTask assigns a task to a node.
func (sd *SchedulerDatabase) AssignTask(t *tes.Task, w *pbs.Node) error {
  err := sd.db.Update(func(tx *bolt.Tx) error {
    return tx.Bucket(TasksQueued).Delete([]byte(t.Id))
	})
  if err != nil {
    return err
  }
  return sd.updateNode(w)
}




// UpdateNode is an RPC endpoint that is used by nodes to send heartbeats
// and status updates, such as completed tasks. The server responds with updated
// information for the node, such as canceled tasks.
func (sd *SchedulerDatabase) UpdateNode(ctx context.Context, req *pbs.Node) (*pbs.UpdateNodeResponse, error) {
  err := sd.updateNode(req)
	resp := &pbs.UpdateNodeResponse{}
	return resp, err
}

func (sd *SchedulerDatabase) updateNode(req *pbs.Node) error {
  return sd.db.Update(func(tx *bolt.Tx) error {
    // Get node
    node, err := sd.getNode(tx, req.Id)
    if err != nil {
      return err
    }

    if node.Version != 0 && req.Version != 0 && node.Version != req.Version {
      return errors.New("Version outdated")
    }

    node.LastPing = time.Now().Unix()
    node.State = req.GetState()
    node.Version = time.Now().Unix()
    return sd.putNode(tx, node)
  })
}

// GetNode gets a node
func (sd *SchedulerDatabase) GetNode(ctx context.Context, req *pbs.GetNodeRequest) (*pbs.Node, error) {
	var node *pbs.Node

	err := sd.db.View(func(tx *bolt.Tx) error {
		node, err = sd.getNode(tx, req.Id)
		return err
	})
	if err == errNotFound {
		return nil, grpc.Errorf(codes.NotFound, fmt.Sprintf("%v: nodeID: %s", err.Error(), req.Id))
	}
  if err != nil {
    return nil, err
  }
	return node, nil
}

// CheckNodes is used by the scheduler to check for dead/gone nodes.
// This is not an RPC endpoint
func (sd *SchedulerDatabase) CheckNodes() error {
  // Track dead node IDs for removal at the end.
  var dead []string

  sd.db.View(func(tx *bolt.Tx) error {

    // Loop over all the nodes, looking for dead nodes to delete.
		bucket := tx.Bucket(Nodes)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

      node, err := sd.getNode(tx, string(k))
      if err != nil {
        // If there's an error, just skip this node.
        continue
      }

			if node.LastPing == 0 {
				// This shouldn't be happening, because nodes should be
				// created with LastPing, but give it the benefit of the doubt
				// and leave it alone.
				continue
			}

			lastPing := time.Unix(node.LastPing, 0)
			d := time.Since(lastPing)

      init := node.State == pbs.NodeState_UNINITIALIZED ||
				      node.State == pbs.NodeState_INITIALIZING

      if (init && d > sd.conf.Scheduler.NodeInitTimeout) ||
			    d > sd.conf.Scheduler.NodeDeadTimeout {

        dead = append(dead, string(k))
      }
		}
		return nil
	})

  if len(dead) == 0 {
    return nil
  }

  // Delete dead nodes.
  return sd.db.Update(func(tx *bolt.Tx) error {
    for _, id := range dead {
      // Explicitly ignoring the error because there's nothing
      // we can do, and the node will be checked again in a future
      // call to CheckNodes().
      _ = tx.Bucket(Nodes).Delete([]byte(id))
    }
  })
}

// ListNodes is an API endpoint that returns a list of nodes.
func (sd *SchedulerDatabase) ListNodes(ctx context.Context, req *pbs.ListNodesRequest) (*pbs.ListNodesResponse, error) {
	resp := &pbs.ListNodesResponse{}
	resp.Nodes = []*pbs.Node{}

	err := sd.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(Nodes)
		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			node, _ := sd.getNode(tx, string(k))
			resp.Nodes = append(resp.Nodes, node)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (sd *SchedulerDatabase) getNode(tx *bolt.Tx, id string) (*pbs.Node, error) {
	node := &pbs.Node{
		Id: id,
	}

	data := tx.Bucket(Nodes).Get([]byte(id))
  // getNode returns a new, empty node when the ID is not found.
  // This allows new nodes to easily register by calling
  // UpdateNode()
	if data == nil {
    return node, errNotFound
	}

  err := proto.Unmarshal(data, node)
  if err != nil {
    return nil, err
  }

	if node.Metadata == nil {
		node.Metadata = map[string]string{}
	}

	return node, nil
}

func (sd *SchedulerDatabase) putNode(tx *bolt.Tx, node *pbs.Node) error {
  log.Info("put", node)
	data, err := proto.Marshal(node)
  if err != nil {
    return err
  }
	return tx.Bucket(Nodes).Put([]byte(node.Id), data)
}
