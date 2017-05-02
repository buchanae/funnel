package server

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
)

type TaskProxy struct {
  dts DtsClient
}

func NewTaskProxy(conf config.Config) (*TaskProxy, error) {
}

func (p *TaskProxy) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
  // Get the strategy from the task tags. Default to "file_routed".
  strategy := "file_routed"
  if s, ok := task.Tags["strategy"]; ok {
    strategy = s
  }

  // Track the locations where the inputs exist.
  // Used in routing decisions.
  locations := locationSet{}

  for _, input := range task.Inputs {
    // TODO fail on non-ccc protocols? e.g s3://

    resp, err := p.dts.Get(input.Url)
    /*
      {
        "cccId": "foo/bar",
        "name": "nohup.out",
        "size": 0,
        "location": [{
          "site": "http://10.73.127.6",
          "path": "/cluster_share/home/buchanae",
          "timestampUpdated": 1489431561,
          "user": {
            "name": "buchanae"
          }
        }]
      }
    */

    locations.Include(resp.CCCID, resp.Location)
  }

  bestSite, ok := locations.BestSite()
  if !ok {
    // TODO FAIL
    // TODO in a future version, this would take the next best site?
  }


}

// GetTask gets a task, which describes a running task
func (p *TaskProxy) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
}


type locationSet struct {
  counts map[string]int
  inputs []string
}
func (locset *locationSet) Include(id string, locs []Location) {
  locset.inputs = append(locset.inputs, id)
  for _, loc := range locs {
    locset.counts[loc.Site] += 1
  }
}
func (locset *locationSet) BestSite() (string, bool) {
  for site, count := locset.counts {
    if count == len(locset.inputs) {
      return site, true
    }
  }
  return "", false
}
