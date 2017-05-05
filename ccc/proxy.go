package ccc

import (
	"errors"
	dtsmocks "github.com/ohsu-comp-bio/funnel/ccc/dts/mocks"
	"github.com/ohsu-comp-bio/funnel/ccc/dts"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
)

var log = logger.New("ccc")
var ErrNoSite = errors.New("no site found")
var ErrBadSite = errors.New("can't connect to site")

type TaskProxy struct {
	conf   config.Config
	dts    dts.Client
	mapper SiteMapper
}

// NewTaskProxy
func NewTaskProxy(conf config.Config) (*TaskProxy, error) {
	dtsClient, err := dts.NewClient(conf.CCC.DTSAddress)
	if err != nil {
		return nil, err
	}
	mapper := &siteMapper{conf: conf}
	return &TaskProxy{conf, dtsClient, mapper}, nil
}

func NewDemoProxy(conf config.Config) *TaskProxy {
	mapper := &siteMapper{conf: conf}
	dtsMock := new(dtsmocks.Client)
  for fileID, siteIDs := range conf.CCC.DTSDemo {
    dtsMock.SetFileSites(fileID, siteIDs)
  }
  proxy := &TaskProxy{conf, dtsMock, mapper}
  return proxy
}

// CreateTask
func (p *TaskProxy) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
	bestSite, err := routeTask(p.conf, p.dts, task)
	if err != nil {
		return nil, err
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

func (p *TaskProxy) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{}, nil
}
