package server

import (
	"github.com/boltdb/bolt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

type transitioner struct {
	id string
	tx *bolt.Tx
  tb *TaskBolt
}

func (th *transitioner) Dequeue(to tes.State) error {
	err := th.tx.Bucket(TasksQueued).Delete([]byte(th.id))
	if err != nil {
		return err
	}
	return th.SetState(to)
}

func (th *transitioner) Queue() error {
  var err error
	err = th.tx.Bucket(TasksQueued).Put([]byte(th.id), []byte{})
	if err != nil {
		return err
	}

  task := &tes.Task{}
  err = loadBasicTaskView(th.tx, th.id, task)
  if err != nil {
    return err
  }

  err = th.tb.backend.Submit(task)
	if err != nil {
		return err
	}
	return th.SetState(tes.Queued)
}

func (th *transitioner) SetState(to tes.State) error {
	log.Info("Set task state", "taskID", th.id, "state", to.String())
	return th.tx.Bucket(TaskState).Put([]byte(th.id), []byte(to.String()))
}

func (taskBolt *TaskBolt) transitionTaskState(tx *bolt.Tx, id string, to tes.State) error {
	from := getTaskState(tx, id)
	th := &transitioner{id, tx, taskBolt}
	return tes.Transition(from, to, th)
}
