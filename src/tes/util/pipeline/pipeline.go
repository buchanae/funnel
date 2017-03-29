package pipeline

import (
  "context"
  "sync"
)

func WithContext(ctx context.Context) *Pipeline {
  p := &Pipeline{ctx: ctx}
  go func() {
    <-ctx.Done()
    p.setErr(ctx.Err())
  }()
  return p
}

type Pipeline struct {
  done bool
  err error
  ctx context.Context
  mtx sync.Mutex
}

func (p *Pipeline) setErr(err error) {
  p.mtx.Lock()
  defer p.mtx.Unlock()
  if !p.done {
    p.err = err
  }
}

func (p *Pipeline) Err() error {
  p.mtx.Lock()
  defer p.mtx.Unlock()
  return p.err
}

func (p *Pipeline) Run(runfunc func() error) {
	// If the pipeline is already complete (perhaps because a previous step failed)
	// skip the step.
  if !p.done {
    err := runfunc()
    p.setErr(err)
	}
}
