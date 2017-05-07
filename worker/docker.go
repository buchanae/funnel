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

type DockerExecutor struct {
  TaskLogger
  task *tes.Task
  mapped *tes.Task
  RemoveContainer bool
}

func (b *DockerExecutor) Execute(ctx context.Context, i int) error {
  e := b.task.Executors[i]
  m := b.mapped.Executors[i]

  name := fmt.Sprintf("%s-%d", b.task.Id, i)
  stdin := util.ReaderOrEmpty(m.Stdin)

  stdout := util.WriterOrDiscard(m.Stdout)
  stdout = io.MultiWriter(b.TaskLogger.ExecutorStdout(i), stdout)

  stderr := util.WriterOrDiscard(m.Stderr)
  stderr = io.MultiWriter(b.TaskLogger.ExecutorStderr(i), stderr)

  cmd := Docker{
    ContainerName:   name,
    ImageName:       e.ImageName,
    Cmd:             e.Cmd,
    Volumes:         MappedVolumes(b.task, b.mapped),
    Workdir:         e.Workdir,
    Ports:           e.Ports,
    RemoveContainer: b.RemoveContainer,
    Environ:         e.Environ,
    Stdin:           stdin,
    Stdout:          stdout,
    Stderr:          stderr,
  }

  go func() {
	  d.Log.Info("Fetching container metadata")
    ports := inspect(ctx, name)
    b.TaskLogger.ExecutorPorts(i, ports)
		log.Debug("Executor ports:", "ports", ports)
  }()

  ip, iperr := externalIP()
  if iperr == nil {
    b.TaskLogger.ExecutorHostIP(i, ip)
  }

  defer stop(name)

  select {
  case err := <-ctx.Done():
    return err
  case err := <-cmd.Run():
    return err
  }
}


// Volume represents a volume mounted into a docker container.
// This includes a HostPath, the path on the host file system,
// and a ContainerPath, the path on the container file system,
// and whether the volume is read-only.
type Volume struct {
	// The path in tes worker.
	HostPath string
	// The path in Docker.
	ContainerPath string
	Readonly      bool
}

func MappedVolumes(task, mapped *tes.Task) []Volume {
  var volumes []Volume

	for i, _ := range task.Inputs {
    volumes = append(volumes, Volume{
      HostPath: mapped.Inputs[i].Path,
      ContainerPath: task.Inputs[i].Path,
      Readonly: true,
    })
	}

	for i, _ := range task.Volumes {
    volumes = append(volumes, Volume{
      HostPath: mapped.Volumes[i],
      ContainerPath: task.Volumes[i],
      Readonly: false,
    })
	}

	for i, output := range task.Outputs {
    hp := mapped.Outputs[i].Path
    cp := output.Path

    if output.Type == tes.FileType_FILE {
      hp = filepath.Dir(hp)
      containterPath = filepath.Dir(cp)
    }

    volumes = append(volumes, Volume{
      HostPath: hp,
      ContainerPath: cp,
      Readonly: false,
    })
	}

  return volumes
}

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
func (dcmd Docker) Run() error {
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
}

// Inspect returns metadata about the container (calls "docker inspect").
func inspect(ctx context.Context, name string) []*tes.Ports {
	dclient := setupDockerClient()
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
			metadata, err := dclient.ContainerInspect(ctx, name)

			if err == nil && metadata.State.Running == true {
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
func stop(name string) error {
  log.Info("Stopping container", "container", name)
	log.Info("Stopping container", "container", dcmd.ContainerName)
	dclient, derr := util.NewDockerClient()
	if derr != nil {
		return derr
	}
	defer dclient.Close()
	timeout := time.Second * 10
  return dclient.ContainerStop(context.Background(), name, &timeout)
}

func formatVolumeArg(v Volume) string {
	mode := "rw"
	if v.Readonly {
		mode = "ro"
	}
	return fmt.Sprintf("%s:%s:%s", v.HostPath, v.ContainerPath, mode)
}
