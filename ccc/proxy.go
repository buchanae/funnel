package ccc

import (
  "errors"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/ccc/dts"
	"golang.org/x/net/context"
  "strings"
)

var log = logger.New("ccc")
var ErrNoSite = errors.New("no site found")
var ErrBadSite = errors.New("can't connect to site")

type TaskProxy struct {
  dts dts.Client
  mapper SiteMapper
}

// NewTaskProxy
func NewTaskProxy(conf config.Config) (*TaskProxy, error) {
  dtsClient, err := dts.NewClient(conf.CCC.DTSAddress)
  if err != nil {
    return nil, err
  }
  mapper := &siteMapper{conf: conf}
  return &TaskProxy{dtsClient, mapper}, nil
}

// CreateTask
func (p *TaskProxy) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
  // Track the locations where the inputs exist.
  // Used in routing decisions.
  locations := locationSet{}

  // Add all CCC task inputs to the location tracking set.
  for _, input := range task.Inputs {
    if strings.HasPrefix(input.Url, "ccc://") {
      url := strings.TrimPrefix(input.Url, "ccc://")
      resp, err := p.dts.GetFile(url)
      if err != nil {
        return nil, err
      }
      locations.Include(resp.ID, resp.Location)
    }
  }

  // Get the site that contains all the inputs.
  // If there is not site that contains ALL inputs, "ok" will be false.
  bestSite, ok := locations.BestSite()
  if !ok {
    // No appropriate site could be found, return an error.
    return nil, ErrNoSite
  }

  // Get a client for the best site
  client, err := p.mapper.Client(bestSite)
  if err != nil {
    return nil, err
  }

  // Call CreateTask on the best site
  resp, err := client.CreateTask(ctx, task)

  // Transform the task ID into a global task ID,
  // which allows the proxy to easily map the task back to the site
  // in future calls to Get/CancelTask.
  if resp != nil {
    resp.Id = p.mapper.GlobalID(bestSite, resp.Id)
  }
  return resp, err
}

func (p *TaskProxy) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {

  // Get a client for the task's site based on the global ID.
  site, _ := p.mapper.Site(req.Id)
  client, err := p.mapper.Client(site)
  if err != nil {
    return nil, err
  }

  // Save the global ID for later.
  // Transform the request to have the local ID.
  gid := req.Id
  req.Id, err = p.mapper.LocalID(req.Id)
  if err != nil {
    return nil, err
  }

  // Call GetTask on the site.
  resp, err := client.GetTask(ctx, req)

  // Transform the response back to the global ID.
  if resp != nil {
    resp.Id = gid
  }
  return resp, err
}

func (p *TaskProxy) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
  resp := &tes.ListTasksResponse{}

  // Loop over all the sites, calling ListTasks on each.
  // Concantenate the results.
  //
  // TODO this is shortcut and doesn't properly manage the sizes
  //      of the respones, pagination, etc.
  for _, site := range p.mapper.Sites() {
    client, err := p.mapper.Client(site)

    // Call ListTasks on the site
    r, err := client.ListTasks(ctx, req)
    if err != nil {
      return nil, err
    }

    for _, task := range r.Tasks {
      // Transform to global ID
      task.Id = p.mapper.GlobalID(site, task.Id)
      resp.Tasks = append(resp.Tasks, task)
    }
  }
  return resp, nil
}


// CancelTask
func (p *TaskProxy) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {

  // Get a client for the task's site based on the global ID.
  site, _ := p.mapper.Site(req.Id)
  client, err := p.mapper.Client(site)
  if err != nil {
    return nil, err
  }

  // Transform global ID to local ID
  req.Id, err = p.mapper.LocalID(req.Id)
  if err != nil {
    return nil, err
  }

  return client.CancelTask(ctx, req)
}


// locationSet helps track the locations for all task inputs
// in order to choose the best site to run the task on
// (which would be the site that contains all inputs)
type locationSet struct {
  counts map[string]int
  inputs []string
}
// Include a set of locations for the input ID
func (locset *locationSet) Include(id string, locs []dts.Location) {
  if locset.counts == nil {
    locset.counts = map[string]int{}
  }
  locset.inputs = append(locset.inputs, id)
  for _, loc := range locs {
    locset.counts[loc.Site] += 1
  }
}

// Get the site that includes all the inputs
func (locset *locationSet) BestSite() (string, bool) {
  for site, count := range locset.counts {
    if count == len(locset.inputs) {
      return site, true
    }
  }
  return "", false
}
