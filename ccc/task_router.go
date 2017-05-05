package ccc

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/ccc/dts"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"sort"
	"strings"
)

type RouterStrategy string

const (
	RoutedFile      RouterStrategy = "routed_file"
	PushFile                       = "push_file"
	FetchFile                      = "fetch_file"
	UnknownStrategy                = "UNKNOWN"
)

func routeTask(conf config.Config, dts dts.Client, task *tes.Task) (string, error) {

	strategy, err := getTaskStrategy(task)
	if err != nil {
		return "", err
	}

  // TODO move this hack
  if len(task.Inputs) == 0 {
    return conf.CCC.CentralSite, nil
  }

	// inputSites helps track which sites have input files available locally
	// and which are fetchable from central. This helps make decisions based
	// on the task strategy below.
	t := inputSites{}
	for _, input := range task.Inputs {
		// Ignore non-CCC input URLs.
		if url, ok := cccInputURL(input.Url); ok {
			// Get the DTS entry
			entry, err := dts.Get(url)
			if err != nil {
				return "", err
			}
			// Add the input site locations to the input sites tracker.
			for _, loc := range entry.Location {
				if loc.Site == conf.CCC.CentralSite {
					t.Fetchable(url)
				}
				t.Local(loc.Site, url)
			}
		}
	}

	var candidates []candidate

	switch strategy {
	// RoutedFile tasks:
	// - must satisfy all inputs in the local site
	//   i.e. must not fetch from a remote site.
	// - may run on the central site.
	case RoutedFile:
		t.RequireLocal = true
		candidates = t.Candidates()

	// PushFile tasks:
	// - must satisfy all inputs in the local site
	//   i.e. must not fetch from a remote site.
	// - must not be run on the central site.
	case PushFile:
		t.RequireLocal = true
		t.ExcludeSite(conf.CCC.CentralSite)
		candidates = t.Candidates()

	// FetchFile tasks:
	// - must satisfy all inputs either in the local site,
	//   or by fetching from the central site.
	// - must not run on the central site
	// - must fail if the inputs cannot be satisfied
	//   either locally or by fetching from the central site.
	// - should run on the site that needs to fetch the
	//   fewest files from central.
	case FetchFile:
		t.ExcludeSite(conf.CCC.CentralSite)
		candidates = t.Candidates()
		sort.Sort(byLocality(candidates))
	}

	// No best site was found, so return an error.
	if len(candidates) == 0 {
		return "", ErrNoSite
		//fmt.Errorf("no candidate found to match %s strategy", strategy)
	}
	best := candidates[0]
	return best.name, nil
}

type candidate struct {
	name     string
	locality int
}
type byLocality []candidate

func (b byLocality) Len() int           { return len(b) }
func (b byLocality) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byLocality) Less(i, j int) bool { return b[i].locality < b[j].locality }

// inputSites helps track input + site mappings from the DTS
// and helps make decisions in routing tasks/files.
type inputSites struct {
	RequireLocal bool
	allInputs    map[string]bool
	excluded     map[string]bool
	fetchable    map[string]bool
	sites        map[string]map[string]bool
}

func (i *inputSites) Fetchable(url string) {
	if i.fetchable == nil {
		i.fetchable = map[string]bool{}
	}
	i.fetchable[url] = true
}

func (i *inputSites) Local(site, url string) {
	if i.allInputs == nil {
		i.allInputs = map[string]bool{}
	}
	i.allInputs[url] = true

	if i.sites == nil {
		i.sites = make(map[string]map[string]bool)
	}
	local, ok := i.sites[site]
	if !ok {
		local = map[string]bool{}
		i.sites[site] = local
	}
	local[url] = true
}

func (i *inputSites) ExcludeSite(site string) {
	if i.excluded == nil {
		i.excluded = map[string]bool{}
	}
	i.excluded[site] = true
}

func (i *inputSites) Candidates() []candidate {
	var candidates []candidate
	// Determine which sites are valid candidates.
	for site, local := range i.sites {

		_, isExcluded := i.excluded[site]
		if isExcluded {
			continue
		}

		// Check all known inputs against this site.
		valid := true
		locality := 0
		for input := range i.allInputs {
			_, isLocal := local[input]
			_, isFetchable := i.fetchable[input]
			if (!isLocal && i.RequireLocal) || (!isLocal && !isFetchable) {
				valid = false
				break
			}
			if isLocal {
				locality += 1
			}
		}
		if valid {
			candidates = append(candidates, candidate{site, locality})
		}
	}
	return candidates
}

func getTaskStrategy(task *tes.Task) (RouterStrategy, error) {
	for key, value := range task.Tags {
		if strings.ToLower(key) == "strategy" {
			strategy := RouterStrategy(strings.ToLower(value))
			switch strategy {
			case RoutedFile, PushFile, FetchFile:
				return strategy, nil
			case "":
				// Default to RoutedFile when there's a key but no value.
				return RoutedFile, nil
			default:
				// Unknown strategy. Raise an error
				return UnknownStrategy, fmt.Errorf("unknown task strategy: %s", strategy)
			}
		}
	}
	// Default to RoutedFile
	return RoutedFile, nil
}

func cccInputURL(in string) (string, bool) {
	if !strings.HasPrefix(in, "ccc://") {
		return "", false
	}
	return strings.TrimPrefix(in, "ccc://"), true
}
