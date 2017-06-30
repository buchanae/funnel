package gce

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"google.golang.org/api/compute/v1"
	"time"
)

// Helper for creating a wrapper before creating a client
func newClientFromConfig(conf config.Config) (Client, error) {
	w, err := newWrapper(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return &gceClient{
		wrapper:  w,
	}, nil
}

func newCache(conf config.Config) {
  return cache{
		cacheTTL: conf.Backends.GCE.CacheTTL,
  }
}

type gceClient struct {
	// GCE API wrapper
	wrapper Wrapper
}

// Templates queries the GCE API to get details about GCE instance templates.
// If the API client fails to connect, this returns an empty list.
func (s *gceClient) Template(id string) []pbf.Worker {
	s.loadTemplates()
	workers := []pbf.Worker{}

	for id, tpl := range s.templates {

	}
	return workers
}

// loadTemplates loads all the project's instance templates from the GCE API
func (s *gceClient) Templates(project, zone string) {

	// Get the machine types available
	mtresp, mterr := s.wrapper.ListMachineTypes(project, zone)
	if mterr != nil {
		log.Error("Couldn't get GCE machine list",
			"error", mterr,
			"project", project,
			"zone", zone)
		return
	}

	// Get the instance template from the GCE API
	itresp, iterr := s.wrapper.ListInstanceTemplates(project)
	if iterr != nil {
		log.Error("Couldn't get GCE instance templates", iterr)
		return
	}

	s.machineTypes = map[string]*compute.MachineType{}
	s.templates = map[string]*compute.InstanceTemplate{}

	for _, m := range mtresp.Items {
		s.machineTypes[m.Name] = m
	}

	for _, t := range itresp.Items {
		// Only include instance templates with a "funnel" tag
		if hasTag(t) {
			s.templates[t.Name] = t
		}
	}
}

func hasTag(t *compute.InstanceTemplate) bool {
	for _, t := range t.Properties.Tags.Items {
		if t == "funnel" {
			return true
		}
	}
	return false
}


type cache struct {
	// Last time the cache was updated
	cacheTime time.Time
	// How long before expiring the cache
	cacheTTL time.Duration
	// cached templates list
	templates map[string]*compute.InstanceTemplate
	// cached machine types
	machineTypes map[string]*compute.MachineType
}

func (cache) Templates() {
	// Don't query the GCE API if we have cache results
	if s.cacheTime.IsZero() && time.Since(s.cacheTime) < s.cacheTTL {
		return
	}
	s.cacheTime = time.Now()
}
