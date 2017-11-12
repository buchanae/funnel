package scheduler

import (
	"context"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

func (s *Scheduler) GetOffer(j *tes.Task) *Offer {
	offers := []*Offer{}

	// Get the nodes from the funnel server
  var nodes []*pbs.Node
	resp, err := s.DB.ListNodes(context.Background(), &pbs.ListNodesRequest{})

	// If there's an error, return an empty list
	if err == nil {
    nodes = resp.Nodes
	}

	for _, w := range nodes {
		// Filter out nodes that don't match the task request.
		// Checks CPU, RAM, disk space, etc.
		if !Match(w, j, s.Predicates) {
			continue
		}

		sc := DefaultScores(w, j)
		offer := NewOffer(w, j, sc)
		offers = append(offers, offer)
	}

	// No matching nodes were found.
	if len(offers) == 0 {
		return nil
	}

	SortByAverageScore(offers)
	return offers[0]
}
