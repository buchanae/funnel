package worker

import (
	"io/ioutil"
	"os"
	"path"
	pbe "tes/ga4gh"
	"testing"
  "tes/config"
  "tes/logger"
  "tes/storage"
)

func newTestJobRunner() *jobRunner {
	f, _ := ioutil.TempDir("", "funnel-test-resolve-links-")
  c := config.WorkerDefaultConfig()
  c.Storage.Local.AllowedDirs = append(c.Storage.Local.AllowedDirs, f)
	m := NewFileMapper(f)
	r := jobRunner{
		mapper: m,
    conf: c,
    store: &storage.Storage{},
    log: logger,
    ctrl: JobControl{},
    updates: make(logUpdateChan),
	}
  return &r
}

func TestResolveLinks(t *testing.T) {
  r := newTestJobRunner()

  r.wrapper = &pbr.JobWrapper{
    Job: &pbe.Job{
      Task: &pbe.Task{
        Outputs: []*pbe.TaskParameter{
        },
      },
    },
  }

	c, e := ioutil.ReadFile(m.Outputs[0].Path)

	if e != nil {
		log.Error("Error reading file", e)
		t.Error("Error reading file")
		return
	}

	if string(c) != "foo\n" {
		t.Error("Error: unexpected content")
	}
}
