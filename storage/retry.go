package storage

import (
  "context"
  "github.com/ohsu-comp-bio/funnel/proto/tes"
  "github.com/ohsu-comp-bio/funnel/config"
)

// Retrier retries storage operations a fixed about of times.
// If MaxAttempts is not set, it defaults to 3.
type Retrier struct {
  Storage
  MaxAttempts int
}

// Get calls Retrier.Storage.Get up to Retrier.MaxAttempts times.
func (r Retrier) Get(ctx context.Context, url, path string, class tes.FileType) error {
  max := r.MaxAttempts
  if max == 0 {
    max = 3
  }
  var err error
  for i := 0; i < max; i++ {
    err = r.Storage.Get(ctx, url, path, class)
    if err == nil {
      return nil
    }
  }
  return err
}

// Put calls Retrier.Storage.Put up to Retrier.MaxAttempts times.
// If all calls fail, the output file log is empty.
func (r Retrier) Put(ctx context.Context, url, path string, class tes.FileType) ([]*tes.OutputFileLog, error) {
  max := r.MaxAttempts
  if max == 0 {
    max = 3
  }
  var err error
  var log []*tes.OutputFileLog
  for i := 0; i < max; i++ {
    log, err = r.Storage.Put(ctx, url, path, class)
    if err == nil {
      return log, nil
    }
  }
  return nil, err
}
