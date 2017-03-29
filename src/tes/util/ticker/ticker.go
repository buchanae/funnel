package ticker

import (
  "context"
  "time"
)

type Ticker struct {
  ticker time.Ticker
  C <-chan time.Time
}

func (t *Ticker) Stop() {
  t.ticker.Stop()
}

func (t *Ticker) run() {
  for {
    select {
    case t.C <- t.ticker.C:
    case <-ctx.Done():
      t.Stop()
      return
    }
  }
}

func WithContext(ctx context.Context, rate time.Duration) *Ticker {
  t := Ticker{time.NewTicker(rate)}
  go t.run()
  return &t
}
