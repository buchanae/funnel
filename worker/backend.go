
// TODO document behavior of slow consumer of task log updates
type Backend interface {
  TaskLogger
  Storage
  Executor(int) Executor
}

func defaultBackend(conf config.Worker, taskID string) Runner {
	// Map files into this baseDir
	baseDir := path.Join(conf.WorkDir, t.Task.Id)
  prepareDir(baseDir)
  ctrl := NewTaskControl()

  mapper.MapTask(r.wrapper.Task)
	r := &taskRunner{
		wrapper: t,
		mapper:  NewFileMapper(baseDir),
		store:   &storage.Storage{},
		conf:    conf,
		log:     ctrl,
	}
}

type DefaultBackend struct {
	mapper  *FileMapper
  *storage.Storage
}
func (b *DefaultBackend) Executor(i int, task tes.Task) (Executor, error) {
    return &DockerExecutor{
      ImageName:     d.ImageName,
      Cmd:           d.Cmd,
      Volumes:       r.mapper.Volumes,
      Workdir:       d.Workdir,
      Ports:         d.Ports,
      ContainerName: fmt.Sprintf("%s-%d", task.Id, i),
      RemoveContainer: r.conf.RemoveContainer,
			Environ:       d.Environ,
      Stdin: r.Stdin(i),
      Stdout: r.Stdout(i),
      Stderr: r.Stderr(i),
    }
}

// Create working dir
func prepareDir(path string) error {
	dir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return util.EnsureDir(dir)
}

// Configure a task-specific storage backend.
// This provides download/upload for inputs/outputs.
func (r *taskRunner) prepareStorage() error {
	var err error

	for _, conf := range r.conf.Storage {
		r.store, err = r.store.WithConfig(conf)
		if err != nil {
			return err
		}
	}

	return nil
}
