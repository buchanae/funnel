package ccc

import (
	"golang.org/x/net/context"
  "testing"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/config"
  //"github.com/ohsu-comp-bio/funnel/ccc/dts"
  dtsmocks "github.com/ohsu-comp-bio/funnel/ccc/dts/mocks"
  "github.com/ohsu-comp-bio/funnel/proto/tes"
  tesmocks "github.com/ohsu-comp-bio/funnel/proto/tes/mocks"
  mock "github.com/stretchr/testify/mock"
)

func init() {
  logger.ForceColors()
}

func TestCreateTask(t *testing.T) {
  dtsMock := new(dtsmocks.Client)

  task := taskWithInputs(
    "ccc://ccc-id/file-one",
    "ccc://ccc-id/file-two",
    "s3://bkt/file-three",
  )

  getSiteClient := func(address string) (tes.TaskServiceClient, error) {
    if address != "site-one:9090" {
      log.Debug("SITE CONN", address)
      t.Fatal("Unexpected site connection")
    }
    c := new(tesmocks.TaskServiceClient)

    c.On("CreateTask", mock.Anything, mock.Anything, mock.Anything).
      Return(&tes.CreateTaskResponse{
        Id: "test-task-id",
      }, nil)

    return c, nil
  }

  conf := config.DefaultConfig()
  m := &siteMapper{conf, getSiteClient}
  p := TaskProxy{dtsMock, m}

  dtsMock.SetFileSites("ccc-id/file-one", []string{
    "http://site-one", "http://site-two",
  })
  dtsMock.SetFileSites("ccc-id/file-two", []string{
    "http://site-one", "http://site-three",
  })

  resp, err := p.CreateTask(context.Background(), task)
  if err != nil {
    t.Fatal(err)
  }
  if resp.Id != "http://site-one/test-task-id" {
    log.Debug("TASK ID", resp.Id)
    t.Fatal("Unexpected task id")
  }
}

func TestCreateNoSharedSite(t *testing.T) {
  dtsMock := new(dtsmocks.Client)

  task := taskWithInputs(
    "ccc://ccc-id/file-one",
    "ccc://ccc-id/file-two",
    "ccc://ccc-id/file-three",
    "s3://bkt/file-three",
  )
  dtsMock.SetFileSites("ccc-id/file-one", []string{
    "http://site-one",
  })
  dtsMock.SetFileSites("ccc-id/file-two", []string{
    "http://site-two",
  })
  dtsMock.SetFileSites("ccc-id/file-three", []string{
    "http://site-one",
  })

  getSiteClient := func(address string) (tes.TaskServiceClient, error) {
    log.Debug("SITE CONN", address)
    t.Fatal("Unexpected site connection")
    return nil, nil
  }

  conf := config.DefaultConfig()
  m := &siteMapper{conf, getSiteClient}
  p := TaskProxy{dtsMock, m}
  _, err := p.CreateTask(context.Background(), task)

  if err == nil {
    t.Fatal("Expected error")
  }
  if err != ErrNoSite {
    t.Fatal("Unexpected error value")
  }
}

func TestGetTask(t *testing.T) {
  t.Error("FUU")
}

func taskWithInputs(urls ...string) *tes.Task {
  var params []*tes.TaskParameter
  for _, url := range urls {
    params = append(params, &tes.TaskParameter{
      Url: url,
    })
  }
  return &tes.Task{Inputs: params}
}
