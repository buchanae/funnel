package sh

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestServerManualBackendPanic(t *testing.T) {
	_, stderr := run(t, "server_manual_backend.sh")

	if strings.Contains(stderr, "panic") {
		t.Error("server panic")
	}
}

func run(t *testing.T, file string) (string, string) {
	var stdout, stderr bytes.Buffer

	// Capture signals so that cleanup happens properly
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", file)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		t.Error(ctx.Err())
	}

	t.Log("STDOUT\n", stdout.String())
	t.Log("STDERR\n", stderr.String())
	return stdout.String(), stderr.String()
}
