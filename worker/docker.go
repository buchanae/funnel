package worker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Docker is responsible for configuring and running a docker container.
type Docker struct {
	ContainerName   string
	ImageName       string
	Cmd             []string
	Volumes         []Volume
	Workdir         string
	Ports           []*tes.Ports
	RemoveContainer bool
	Environ         map[string]string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
}

// Run runs the Docker command and blocks until done.
func (dcmd Docker) Run(ctx context.Context) error {
	args := []string{"run", "-i"}

	if dcmd.RemoveContainer {
		args = append(args, "--rm")
	}

	if dcmd.Environ != nil {
		for k, v := range dcmd.Environ {
			args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
		}
	}

	if dcmd.Ports != nil {
		for i := range dcmd.Ports {
			hostPort := dcmd.Ports[i].Host
			containerPort := dcmd.Ports[i].Container
			// TODO move to validation?
			if hostPort <= 1024 && hostPort != 0 {
				return fmt.Errorf("Error cannot use restricted ports")
			}
			args = append(args, "-p", fmt.Sprintf("%d:%d", hostPort, containerPort))
		}
	}

	if dcmd.ContainerName != "" {
		args = append(args, "--name", dcmd.ContainerName)
	}

	if dcmd.Workdir != "" {
		args = append(args, "-w", dcmd.Workdir)
	}

	for _, vol := range dcmd.Volumes {
		arg := formatVolumeArg(vol)
		args = append(args, "-v", arg)
	}

	args = append(args, dcmd.ImageName)
	args = append(args, dcmd.Cmd...)

	// Roughly: `docker run --rm -i -w [workdir] -v [bindings] [imageName] [cmd]`
	d.Log.Info("Running command", "cmd", "docker "+strings.Join(args, " "))
	cmd := exec.Command("docker", args...)

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
  // TODO watch context and call stop
}

// Inspect returns metadata about the container (calls "docker inspect").
func (dcmd Docker) Inspect(ctx context.Context) ExecutorMetadata {
	d.Log.Info("Fetching container metadata")
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
				d.Log.Error("Error inspecting container", err)
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
						d.Log.Debug("Found port mapping:", "host", hostPort, "container", containerPort)
					}
				}
				return portMap, nil
			}
		}
	}
}

// Stop stops the container.
func (dcmd Docker) stop() error {
	log.Info("Stopping container", "container", dcmd.ContainerName)
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
