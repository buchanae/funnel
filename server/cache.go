package server

import (
	"context"
	"sync"

	"github.com/ohsu-comp-bio/funnel/tes"
)

// TaskCache caches tasks. CreateTask will calculate a cache key for new tasks.
// If the task is cached, it the ID of the cached task will be returned instead
// of creating a new task.
//
// This is an experimental proof-of-concept. Tasks are cached in memory, and 
// inputs/outputs storage objects are not checked as part of the caching.
type TaskCache struct {
	tes.TaskServiceServer
	cache sync.Map
}

/*
TODO

- should failed tasks be cached? what if system error?
- should task output storage object etags be included in the task cache key?
  this would significantly complicate things, since caching could not be completely
  handled in CreateTask
- Will CreateTask need to make storage calls before returning? This could be a
  substantial hit to the speed of CreateTask.
- Docker image hashes are critical as well.
*/

// CreateTask provides an HTTP/gRPC endpoint for creating a task.
// This is part of the TES implementation.
func (ts *TaskCache) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
	hash, hasherr := tes.Hash(task)
	if hasherr == nil {
		val, ok := ts.cache.Load(hash)
		if ok {
			return val.(*tes.CreateTaskResponse), nil
		}
	}
	resp, err := ts.TaskServiceServer.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	ts.cache.Store(hash, resp)
	return resp, nil
}
