package tests

import (
	"errors"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"golang.org/x/net/context"
	"testing"
	"time"
)

var log = logger.New("tests")

func init() {
	logger.ForceColors()
}

func TestHelloWorld(t *testing.T) {
	srv := NewFunnel(NewConfig())
	srv.Start()
	defer srv.Stop()

	// Run task
	taskID := srv.RunCmd("echo", "hello world")

	// Check task state in DB
	r, _ := srv.DB.GetTask(ctx, &tes.GetTaskRequest{Id: taskID})
}
