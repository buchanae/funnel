package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
  "github.com/docker/go-connections/nat"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var log = logger.Sub("docker")

type Executor struct {
  client *client.Client
  conf worker.ExecConfig
}

// NewExecutor returns a new DockerClient instance.
// This util will attempt to negotiate a Docker API version mismatch.
func NewExecutor(conf worker.ExecConfig) (*Executor, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	// If the api version is not set test if the client can communicate with the
	// server; if not infer API version from error message and inform the client
	// to use that version for future communication
	if os.Getenv("DOCKER_API_VERSION") == "" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err := c.ServerVersion(ctx)
		if err != nil {
			re := regexp.MustCompile(`([0-9\.]+)`)
			version := re.FindAllString(err.Error(), -1)
			if version == nil {
				return nil, errors.New("Can't connect docker client")
			}
			// Error message example:
			//   Error getting metadata for container: Error response from daemon:
      //   client is newer than server (client API version: 1.26, server API version: 1.24)
			os.Setenv("DOCKER_API_VERSION", version[1])
			return NewExecutor(conf)
		}
	}
	return &Executor{c, conf}, nil
}

// Run runs the Docker command and blocks until done.
func (d *Executor) Run(ctx context.Context) (exitCode int64, err error) {

  c := container.Config{
    Cmd: d.conf.Cmd,
    Image: d.conf.ImageName,
    WorkingDir: d.conf.Workdir,
    AttachStdin: cmd.Stdin != nil,
    OpenStdin: cmd.Stdin != nil,
    AttachStdout: cmd.Stdout != nil,
    AttachStderr: cmd.Stderr != nil,
  }
  h := container.HostConfig{
    AutoRemove: d.conf.RemoveContainer,
    ReadonlyRootfs: true,
  }

  // Environment variables
  for k, v := range d.conf.Environ {
    c.Env = append(c.Env, fmt.Sprintf("%s=%s", k, v))
  }

  // Ports
  for _, x := range d.conf.Ports {
    // TODO move to validation?
    if p.Host <= 1024 && p.Host != 0 {
      err = fmt.Errorf("Error cannot use restricted ports")
      return
    }

    cont := nat.Port(fmt.Sprint(x.Container))
    // Set exposed port on container config.
    c.ExposedPorts[cont] = struct{}{}
    // Set bound port on host config.
    h.PortBindings[cont] = append(h.PortBindings[cont], nat.PortBinding{
      HostPort: fmt.Sprint(x.Host),
    })
  }

  // Volume binding
	for _, vol := range d.conf.Volumes {
    h.Binds = append(h.Binds, formatVolumeArg(vol))
	}

  // Create
  resp, err := d.client.ContainerCreate(ctx, c, h, nil, d.conf.ContainerName)


  // Attach stdout stream
  if d.conf.Stdout != nil {
    var so types.HijackedResponse
    so, err = d.client.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
      Stream: true,
      Stdout: true,
    })

    if err != nil {
      return
    }

    defer so.Close()
    io.Pipe(so.Reader, d.conf.Stdout)
  }

  // Attach stderr stream
  if d.conf.Stderr != nil {
    var se types.HijackedResponse
    se, err = d.client.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
      Stream: true,
      Stderr: true,
    })
    defer se.Close()

    if err != nil {
      return
    }

    io.Pipe(se.Reader, d.conf.Stderr)
  }

  // Attach logs stream
  var logs types.HijackedResponse
  logs, lerr := d.client.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
    Stream: true,
    Logs: true,
  })

  if lerr != nil {
    defer logs.Close()
    // TODO
    io.Pipe(logs.Reader, os.Stderr)
  }

  // Block, wait for container
  var wait types.ContainerWaitOKBody
  wait, waiterr = d.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

  select {
  case x := <-wait:
    exitCode = x.StatusCode
    return
  case e := <-waiterr:
    err = e
    return
  }
}

// Inspect returns metadata about the container (calls "docker inspect").
func (d *Executor) Inspect(ctx context.Context) ([]*tes.Ports, error) {
	log.Info("Fetching container metadata")
	// close the docker client connection
	defer d.client.Close()

	for {
		select {
		case <-ctx.Done():
			return nil, nil

		default:
			metadata, err := d.client.ContainerInspect(ctx, d.conf.ContainerName)
			if client.IsErrContainerNotFound(err) {
				break
			}
			if err != nil {
				log.Error("Error inspecting container", err)
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
						log.Debug("Found port mapping:", "host", hostPort, "container", containerPort)
					}
				}
				return portMap, nil
			}
		}
	}
}

// Stop stops the container.
func (d *Executor) Stop(ctx context.Context) error {
	log.Info("Stopping container", "container", d.conf.ContainerName)

	// Set timeout
  // TODO why timeout if we're using context?
	timeout := time.Second * 10
	// Issue stop call
  return d.client.ContainerStop(ctx, d.conf.ContainerName, &timeout)
}

func formatVolumeArg(v Volume) string {
	// `o` is structed as "HostPath:ContainerPath:Mode".
	mode := "rw"
	if v.Readonly {
		mode = "ro"
	}
	return fmt.Sprintf("%s:%s:%s", v.HostPath, v.ContainerPath, mode)
}
