package util

import (
  "github.com/ohsu-comp-bio/util"
  "syscall"
)

var Ctx = SignalContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
