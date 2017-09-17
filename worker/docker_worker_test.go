package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/util"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDockerWorker(t *testing.T) {
	id := "test-id-" + util.GenTaskID()
	// collect events
	col := events.Collector{}
	mw := events.MultiWriter(&col, events.NewEventLogger("test"))
	log := NewEventLogger(id, 0, mw)

	// Prepare temp workdir and input files
	tmp, err := ioutil.TempDir("/tmp", "funnel-worker-stdio-test-")
	if err != nil {
		t.Fatal(err)
	}
	//defer os.RemoveAll(tmp)
	ioutil.WriteFile(tmp+"/download-file.txt", []byte("download foo"), os.ModePerm)

	// Build task and TaskReader mock
	read := &mockTaskReader{
		task: &tes.Task{
			Id: id,
			Inputs: []*tes.TaskParameter{
				{
					Contents: "foo contents",
					Path:     "/opt/test/contents-file.txt",
				},
				{
					Url:  "file://" + tmp + "/download-file.txt",
					Path: "/opt/test/download-file.txt",
				},
			},
			Outputs: []*tes.TaskParameter{
				{
					Url:  "file://" + tmp + "/output-upload.txt",
					Path: "/opt/test/output-file.txt",
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
					Cmd:       []string{"/bin/sh", "-c", `md5sum /opt/test/stderr-2.txt`},
					Stdout:    "/opt/test/output-file.txt",
				},
				/*
									{
										ImageName: "alpine",
										Cmd:       []string{"/bin/sh", "-c", `sleep 100`},
					        },
				*/
			},
			Volumes: []string{"/opt/test"},
		},
		state: tes.State_QUEUED,
	}

	// Build and run the worker
	w := DockerWorker{
		conf: DockerConfig{
			Storage: storage.Config{
				Local: storage.LocalConfig{
					AllowedDirs: []string{tmp},
				},
			},
			UpdateRate: time.Millisecond,
			WorkDir:    tmp,
		},
		read: read,
		log:  log,
	}
	w.Run(context.Background())

	expect := `a435f4c56d3985df68305a5240d76eac  /opt/test/stderr-2.txt`

	b, err := ioutil.ReadFile(tmp + "/output-upload.txt")
	if err != nil {
		t.Error("error reading uploaded output file", err)
	} else if strings.TrimSpace(string(b)) != expect {
		t.Errorf("expected output '%s' but got '%s'", expect, string(b))
	}
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
