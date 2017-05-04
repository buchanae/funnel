package worker

func NewFileBackend(conf config.Worker, taskID string) (*FileBackend, error) {
  workspace, werr := NewWorkspace(conf.WorkDir, taskID)
  store, serr := storage.FromConfig(conf.Storage)
  filetask, fterr := NewFileTask(conf, taskID)
  task, terr := filetask.Task()
  docker := DockerExecutor{
    RemoveContainer: conf.RemoveContainer,
    task: task,
    logger: filetask,
    workspace: workspace,
  }

  if err := util.Check(werr, terr, serr); err != nil {
    return nil, err
  }

  return &FileBackend{
    Logger: log.WithFields("task", taskID),
    FileTaskLogger: filetask,
    FileTaskReader: filetask,
    Storage: store,
    DockerExecutor: docker,
  }, nil
}

type FileBackend struct {
  logger.Logger
  *FileTaskLogger
  *FileTaskReader
  storage.Storage
  *DockerExecutor
}

func (b *FileBackend) Close() {}
