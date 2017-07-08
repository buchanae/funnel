package gcp

import (
  "bytes"
  "cloud.google.com/go/storage"
  "io"
  "github.com/golang/protobuf/proto"
  "fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	uuid "github.com/nu7hatch/gouuid"
)


type taskBucket struct {
  *storage.BucketHandle
  id string
}
func (tb *taskBucket) minimal() *storage.ObjectHandle {
  return tb.Object("tasks/" + tb.id + "/views/minimal")
}
func (tb *taskBucket) basic() *storage.ObjectHandle {
  return tb.Object("tasks/" + tb.id + "/views/basic")
}
func (tb *taskBucket) full() *storage.ObjectHandle {
  return tb.Object("tasks/" + tb.id + "/views/full")
}
func (tb *taskBucket) cancel() *storage.ObjectHandle {
  return tb.Object("tasks/" + tb.id + "/cancel")
}

func viewTask(t *tes.Task, v tes.TaskView) *tes.Task {
  switch v {
  case tes.TaskView_MINIMAL:
    return &tes.Task{
      Id: t.Id,
      State: t.State,
    }

  case tes.TaskView_BASIC:
    o := proto.Clone(t).(*tes.Task)
    // Clear input contents
    for _, i := range t.Inputs {
      i.Contents = ""
    }
    // Clear task/executor logs
    for _, l := range t.Logs {
      for _, e := range l.Logs {
        e.Stdout = ""
        e.Stderr = ""
      }
    }
    return o

  case tes.TaskView_FULL:
    fallthrough
  default:
    return proto.Clone(t).(*tes.Task)
  }
}

func readTask(r io.Reader) (*tes.Task, error) {
  var b bytes.Buffer
  if _, err := io.Copy(&b, r); err != nil {
    return nil, err
  }
  task := &tes.Task{}
  err := proto.Unmarshal(b.Bytes(), task)
  return task, err
}

func writeTask(t *tes.Task, w io.WriteCloser) error {
  b, err := proto.Marshal(t)
  if err != nil {
    return err
  }

  buf := bytes.NewBuffer(b)

  if _, err := io.Copy(w, buf); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
  return nil
}

func genWorkerID() string {
	u, _ := uuid.NewV4()
	return fmt.Sprintf("funnel-worker-%s", u.String())
}

// localize helps make a resource string zone-specific
func localize(zone, resourceType, val string) string {
	return fmt.Sprintf("zones/%s/%s/%s", zone, resourceType, val)
}
