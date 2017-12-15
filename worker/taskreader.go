package worker

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"os"
)

// GenericTaskReader provides read access to tasks.
type GenericTaskReader struct {
	get    func(ctx context.Context, in *tes.GetTaskRequest) (*tes.Task, error)
	taskID string
}

// NewGenericTaskReader returns a new generic task reader.
func NewGenericTaskReader(get func(ctx context.Context, in *tes.GetTaskRequest) (*tes.Task, error), taskID string) *GenericTaskReader {
	return &GenericTaskReader{get, taskID}
}

// Task returns the task descriptor.
func (r *GenericTaskReader) Task() (*tes.Task, error) {
	return r.get(context.Background(), &tes.GetTaskRequest{
		Id:   r.taskID,
		View: tes.TaskView_FULL,
	})
}

// State returns the current state of the task.
func (r *GenericTaskReader) State() (tes.State, error) {
	t, err := r.get(context.Background(), &tes.GetTaskRequest{
		Id:   r.taskID,
		View: tes.TaskView_MINIMAL,
	})
	return t.GetState(), err
}

type FileTaskReader struct {
	task *tes.Task
}

func NewFileTaskReader(path string) (*FileTaskReader, error) {
	var task tes.Task

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	err = jsonpb.Unmarshal(f, &task)
	if err != nil {
		return nil, fmt.Errorf("can't load task: %s", err)
	}
	task.Id = tes.GenerateID()
	return &FileTaskReader{&task}, nil
}

func (f *FileTaskReader) Task() (*tes.Task, error) {
	return f.task, nil
}

func (f *FileTaskReader) State() (tes.State, error) {
	return tes.Unknown, nil
}
