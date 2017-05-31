package badger

import (
	"bytes"
	"github.com/dgraph-io/badger/badger"
	proto "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func taskKey(id string) []byte {
	return append([]byte("tasks"), []byte(id)...)
}

func queueKey(id string) []byte {
	return append([]byte("queue"), []byte(id)...)
}

func workerKey(id string) []byte {
	return append([]byte("workers"), []byte(id)...)
}

type TaskBadger struct {
	db   *badger.KV
	conf config.Config
}

func NewTaskBadger(conf config.Config) (*TaskBadger, error) {
	opt := badger.DefaultOptions
	opt.Dir = conf.WorkDir
	util.EnsureDir(conf.WorkDir)

	db, err := badger.NewKV(&opt)
	log.Info("BADGER", "db", db == nil, "err", err)
	if err != nil {
		return nil, err
	}

	return &TaskBadger{db: db, conf: conf}, nil
}

// ReadQueue returns a slice of queued Tasks. Up to "n" tasks are returned.
func (tb *TaskBadger) ReadQueue(n int) []*tes.Task {
	var tasks []*tes.Task
	opt := badger.IteratorOptions{
		PrefetchSize: n,
	}

	itr := tb.db.NewIterator(opt)
	defer itr.Close()

	for itr.Seek(queueKey("")); itr.Valid(); itr.Next() {
		item := itr.Item()
		key := item.Key()
		if !bytes.HasPrefix(key, queueKey("")) {
			break
		}
		taskID := string(bytes.TrimPrefix(key, queueKey("")))
		task, _ := tb.getTask(taskID)
		tasks = append(tasks, task)
	}
	return tasks
}

// GenTaskID generates a task ID string.
// IDs are globally unique and sortable.
func GenTaskID() string {
	log.Debug("GEN")
	id := xid.New()
	return id.String()
}

// CreateTask provides an HTTP/gRPC endpoint for creating a task.
// This is part of the TES implementation.
func (tb *TaskBadger) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
	log.Debug("CreateTask")

	if err := tes.Validate(task); err != nil {
		log.Error("Invalid task message", "error", err)
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	var err error
	taskID := GenTaskID()

	key := taskKey(taskID)
	task.Id = taskID
	task.State = tes.State_QUEUED

	v, err := proto.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = tb.db.Set(key, v)
	if err != nil {
		return nil, err
	}

	qkey := queueKey(taskID)
	err = tb.db.Set(qkey, nil)
	if err != nil {
		return nil, err
	}

	return &tes.CreateTaskResponse{Id: taskID}, nil
}

// GetTask gets a task, which describes a running task
func (tb *TaskBadger) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
	log.Debug("GetTask")
	return tb.getTask(req.Id)
}

func (tb *TaskBadger) getTask(id string) (*tes.Task, error) {
	key := taskKey(id)
	var item badger.KVItem
	err := tb.db.Get(key, &item)
	if err != nil {
		return nil, err
	}
	b := item.Value()
	var task tes.Task
	proto.Unmarshal(b, &task)
	return &task, nil
}

func (tb *TaskBadger) putTask(task *tes.Task) error {
	key := taskKey(task.Id)
	b, _ := proto.Marshal(task)
	err := tb.db.Set(key, b)
	return err
}

// ListTasks returns a list of taskIDs
func (tb *TaskBadger) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
	log.Debug("ListTasks")

	var tasks []*tes.Task

	opt := badger.DefaultIteratorOptions
	itr := tb.db.NewIterator(opt)
	defer itr.Close()

	for itr.Seek(taskKey("")); itr.Valid(); itr.Next() {
		item := itr.Item()
		key := item.Key()
		if !bytes.HasPrefix(key, taskKey("")) {
			break
		}
		val := item.Value()
		var task tes.Task
		proto.Unmarshal(val, &task)
		tasks = append(tasks, &task)
	}

	return &tes.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

// CancelTask cancels a task
func (tb *TaskBadger) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {
	log.Debug("CancelTask")
	tb.transitionTaskState(req.Id, tes.State_CANCELED)
	return &tes.CancelTaskResponse{}, nil
}

func (tb *TaskBadger) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{}, nil
}
