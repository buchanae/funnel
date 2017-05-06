type Executor struct {
  *tes.Executor
  Name   string
	// Volumes         []Volume
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
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

type DockerFactory struct {
  TaskLogger
  task *tes.Task
  mapped *tes.Task
  conf config.Worker
}

func (b *DockerFactory) Executor(i int) (Executor, error) {
  e := b.task.Executors[i]
  m := b.mapped.Executors[i]
  return &Executor{
    e,
    Name:   fmt.Sprintf("%s-%d", b.task.Id, i),
    //Volumes:         MappedVolumes(b.task, b.mapped),
    //Stdin:           util.ReaderOrEmpty(m.Stdin),
    Stdout:          b.TaskWriter.ExecutorStdout(m.Stdout),
    Stderr:          b.TaskWriter.ExecutorStderr(m.Stderr),
  }, nil
}



func MappedVolumes(task, mapped *tes.Task) []Volume {
  var volumes []Volume

	for i, _ := range task.Inputs {
    volumes = append(volumes, Volume{
      HostPath: mapped.Inputs[i].Path,
      ContainerPath: task.Inputs[i].Path,
      Readonly: true,
    }
	}

	for i, _ := range task.Volumes {
    volumes = append(volumes, Volume{
      HostPath: mapped.Volumes[i],
      ContainerPath: task.Volumes[i],
      Readonly: false,
    }
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
    }
	}

  return volumes
}
