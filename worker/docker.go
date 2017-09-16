package worker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// DockerCmd is responsible for configuring and running a docker container.
type DockerCmd struct {
  TaskLogger
  *Stdio
  ExecIndex int
  Exec *tes.Executor
  Volumes []Volume
  RemoveContainer bool
  ContainerName string
}

func (dcmd *DockerCmd) Run(ctx context.Context) error {

	done := make(chan error, 1)

	go func() {
		done <- dcmd.runcmd()
	}()
  go dcmd.Inspect(ctx)

  select {
  case <-ctx.Done():
    // Likely the task was canceled.
    dcmd.Stop()
    return ctx.Err()

  case result := <-done:
    code := GetExitCode(result)
    dcmd.TaskLogger.ExecutorExitCode(dcmd.ExecIndex, code)

    if result != nil {
      return ErrExecFailed(result)
    }
    return nil
  }
}


// Run runs the Docker command and blocks until done.
func (dcmd *DockerCmd) runcmd() error {
	// (Hopefully) temporary hack to sync docker API version info.
	// Don't need the client here, just the logic inside NewDockerClient().
	_, derr := util.NewDockerClient()
	if derr != nil {
		log.Error("Can't connect to Docker", derr)
		return derr
	}

	cmd := exec.Command("docker", dcmd.Args()...)
  cmd.Stdin = dcmd.Stdio.In
  cmd.Stdout = dcmd.Stdio.Out
  cmd.Stderr = dcmd.Stdio.Err

	return cmd.Run()
}

// Inspect returns metadata about the container (calls "docker inspect").
func (dcmd *DockerCmd) Inspect(ctx context.Context) {
	t := time.NewTimer(time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:

			ports, err := dcmd.inspect(ctx)
			if err != nil && !client.IsErrContainerNotFound(err) {
				break
			}
			dcmd.TaskLogger.ExecutorPorts(dcmd.ExecIndex, ports)
			return
		}
	}
}

func (dcmd *DockerCmd) inspect(ctx context.Context) ([]*tes.Ports, error) {
	dclient, derr := util.NewDockerClient()
	if derr != nil {
		return nil, derr
	}
	// close the docker client connection
	defer dclient.Close()

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:

			metadata, err := dclient.ContainerInspect(ctx, dcmd.ContainerName)
			if client.IsErrContainerNotFound(err) {
				break
			}
			if err != nil {
				break
			}

			if metadata.State.Running == true {
				var portMap []*tes.Ports
				// extract exposed host port from
				// https://godoc.org/github.com/docker/go-connections/nat#PortMap
				for k, v := range metadata.NetworkSettings.Ports {
					// will end up taking the last binding listed
					for i := range v {
						p := strings.Split(string(k), "/")
						containerPort, err := strconv.Atoi(p[0])
						if err != nil {
							return nil, err
						}
						hostPort, err := strconv.Atoi(v[i].HostPort)
						if err != nil {
							return nil, err
						}
						portMap = append(portMap, &tes.Ports{
							Container: uint32(containerPort),
							Host:      uint32(hostPort),
						})
					}
				}
				return portMap, nil
			}
		}
	}
}

// Stop stops the container.
func (dcmd *DockerCmd) Stop() error {
	dclient, derr := util.NewDockerClient()
	if derr != nil {
		return derr
	}
	// close the docker client connection
	defer dclient.Close()
	// Set timeout
	timeout := time.Second * 10
	// Issue stop call
	// TODO is context.Background right?
	err := dclient.ContainerStop(context.Background(), dcmd.ContainerName, &timeout)
	return err
}

func formatVolumeArg(v Volume) string {
	// `o` is structed as "HostPath:ContainerPath:Mode".
	mode := "rw"
	if v.Readonly {
		mode = "ro"
	}
	return fmt.Sprintf("%s:%s:%s", v.HostPath, v.ContainerPath, mode)
}

func (dcmd *DockerCmd) Args() []string {
  e := dcmd.Exec
	args := []string{"run", "-i"}

	if dcmd.RemoveContainer {
		args = append(args, "--rm")
	}

  for k, v := range e.Environ {
    args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
  }

  for _, p := range e.Ports {
    args = append(args, "-p", fmt.Sprintf("%d:%d", p.Host, p.Container))
  }

	if dcmd.ContainerName != "" {
		args = append(args, "--name", dcmd.ContainerName)
	}

	if e.Workdir != "" {
		args = append(args, "-w", e.Workdir)
	}

	for _, vol := range dcmd.Volumes {
		arg := formatVolumeArg(vol)
		args = append(args, "-v", arg)
	}

	args = append(args, e.ImageName)
	args = append(args, e.Cmd...)
  return args
}
