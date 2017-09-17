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

func TestOSExecWorker(t *testing.T) {
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
  util.EnsureDir(tmp + "/opt/test")

	// Build task and TaskReader mock
	read := &mockTaskReader{
		task: &tes.Task{
			Id: id,
			Inputs: []*tes.TaskParameter{
				{
					Contents: "foo contents",
					Path:     tmp + "/opt/test/contents-file.txt",
				},
				{
					Url:  "file://" + tmp + "/download-file.txt",
					Path: tmp + "/opt/test/download-file.txt",
				},
			},
			Outputs: []*tes.TaskParameter{
				{
					Url:  "file://" + tmp + "/output-upload.txt",
					Path: tmp + "/opt/test/output-file.txt",
				},
			},
			Executors: []*tes.Executor{
				{
					Cmd:       []string{"/bin/cat"},
					Stdin:     tmp + "/opt/test/contents-file.txt",
					Stdout:    tmp + "/opt/test/stdout-0.txt",
				},
				{
					Cmd:       []string{"/bin/cat", tmp + "/opt/test/stdout-0.txt"},
					Stdout:    tmp + "/opt/test/stdout-1.txt",
				},
				{
					Cmd:       []string{"/bin/sh", "-c", "cat " + tmp + "/opt/test/download-file.txt >&2"},
					Stderr:    tmp + "/opt/test/stderr-2.txt",
				},
				{
					Cmd:       []string{"/bin/sh", "-c", "wc " + tmp + "/opt/test/stderr-2.txt"},
					Stdout:    tmp + "/opt/test/output-file.txt",
				},
				/*
									{
										Cmd:       []string{"/bin/sh", "-c", `sleep 100`},
					        },
				*/
			},
		},
		state: tes.State_QUEUED,
	}

	// Build and run the worker
	w := OSExecWorker{
		conf: OSExecConfig{
			Storage: storage.Config{
				Local: storage.LocalConfig{
					AllowedDirs: []string{tmp},
				},
			},
			UpdateRate: time.Millisecond,
		},
		read: read,
		log:  log,
	}
	w.Run(context.Background())

  expect := "0       2      12 " + tmp + "/opt/test/stderr-2.txt"

	b, err := ioutil.ReadFile(tmp + "/output-upload.txt")
	if err != nil {
		t.Error("error reading uploaded output file", err)
	} else if strings.TrimSpace(string(b)) != expect {
		t.Errorf("expected output '%s' but got '%s'", expect, strings.TrimSpace(string(b)))
	}
}
