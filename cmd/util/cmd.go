package util

import (
  "io"
)

type Interface interface {
  Run(ctx context.Context) error
}

type Command struct {
  Names []string
  Short string
  Args []string
  Out, Err io.Writer
}
