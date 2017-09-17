package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDockerWorker(t *testing.T) {
	// collect events
	col := events.Collector{}
	mw := events.MultiWriter(&col, events.NewEventLogger("test"))
	log := NewEventLogger("task-id-0", 0, mw)

	// Prepare temp workdir and input files
	dir, err := ioutil.TempDir("", "funnel-worker-stdio-test-")
	dir, err = filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "download-file.txt"), []byte("download foo"), os.ModePerm)

	// Build task and TaskReader mock
	read := &mockTaskReader{
		task: &tes.Task{
			Id: "task-id-0",
			Inputs: []*tes.TaskParameter{
				{
					Contents: "foo contents",
					Path:     "/opt/test/contents-file.txt",
				},
				{
					Url:  "file://" + filepath.Join(dir, "download-file.txt"),
					Path: "/opt/test/download-file.txt",
				},
			},
			Outputs: []*tes.TaskParameter{
				{
					Url:  "file://" + filepath.Join(dir, "output-file.txt"),
					Path: "/out/test/output-file.txt",
				},
			},
			Executors: []*tes.Executor{
				{
					ImageName: "alpine",
					Cmd:       []string{"/bin/cat"},
					Stdin:     "/opt/test/contents-file.txt",
					Stdout:    "/opt/test/stdout-0.txt",
				},
				{
					ImageName: "alpine",
					Cmd:       []string{"/bin/cat", "/opt/test/stdout-0.txt"},
					Stdout:    "/opt/test/stdout-1.txt",
				},
				{
					ImageName: "alpine",
					Cmd:       []string{"/bin/sh", "-c", `cat /opt/test/download-file.txt >&2`},
					Stderr:    "/opt/test/stderr-2.txt",
				},
				{
					ImageName: "alpine",
					Cmd:       []string{"/bin/sh", "-c", `sleep 100`},
					Stderr:    "/opt/test/stderr-2.txt",
				},
			},
			Volumes: []string{"/opt/test"},
		},
		state: tes.State_QUEUED,
	}

	// Build worker
	w := DockerWorker{
		conf: DockerConfig{
			Storage: storage.Config{
				Local: storage.LocalConfig{
					AllowedDirs: []string{dir},
				},
			},
			UpdateRate: time.Millisecond,
			WorkDir:    dir,
		},
		read: read,
		log:  log,
	}
	// Run
	w.Run(context.Background())
}

type mockTaskReader struct {
	task    *tes.Task
	taskerr error
	state   tes.State
}

func (m *mockTaskReader) Task() (*tes.Task, error) {
	if m.taskerr != nil {
		return nil, m.taskerr
	}
	return m.task, nil
}
func (m *mockTaskReader) State() tes.State {
	return m.state
}
