package worker

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
  "os"
	"os/exec"
	"strings"
)

// DockerCmd is responsible for configuring and running a docker container.
type DockerCmd struct {
	ImageName       string
	Cmd             []string
	Volumes         []Volume
	Workdir         string
	Ports           []*tes.Ports
	ContainerName   string
	RemoveContainer bool
  ExtraFlags      []string
	Environ         map[string]string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
	Event           *events.ExecutorWriter
}

// Run runs the Docker command and blocks until done.
func (dcmd DockerCmd) Run(ctx context.Context) error {

	dcmd.Event.Info("Running command", "cmd", strings.Join(dcmd.Cmd, " "))
	cmd := exec.Command(dcmd.Cmd[0], dcmd.Cmd[1:]...)

  cmd.Env = append(cmd.Env, os.Environ()...)
  for k, v := range dcmd.Environ {
    cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
  }

	if dcmd.Workdir != "" {
    cmd.Dir = dcmd.Workdir
	}

	if dcmd.Stdin != nil {
		cmd.Stdin = dcmd.Stdin
	}
	if dcmd.Stdout != nil {
		cmd.Stdout = dcmd.Stdout
	}
	if dcmd.Stderr != nil {
		cmd.Stderr = dcmd.Stderr
	}
	return cmd.Run()
}

// Inspect returns metadata about the container (calls "docker inspect").
func (dcmd DockerCmd) Inspect(ctx context.Context) ([]*tes.Ports, error) {
  return nil, nil
}

// Stop stops the container.
func (dcmd DockerCmd) Stop() error {
  return nil
}

func formatVolumeArg(v Volume) string {
	// `o` is structed as "HostPath:ContainerPath:Mode".
	mode := "rw"
	if v.Readonly {
		mode = "ro"
	}
	return fmt.Sprintf("%s:%s:%s", v.HostPath, v.ContainerPath, mode)
}
