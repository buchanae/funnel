package worker

func NewFileBackend(conf config.Worker, taskID string) (*FileBackend, error) {
  workspace, werr := NewWorkspace(conf.WorkDir, taskID)
  NewTaskState()

  task := // TODO get task from file
  store, serr := storage.FromConfig(conf.Storage)

  if serr != nil {
    return nil, serr
  }

  canceler := NewFileTaskCanceler(workspace)

  return &FileBackend{
    Logger: log.WithFields("task", taskID),
    FileTaskLogger: &FileTaskLogger{workspace, taskID},
    Storage: store,
    task: task,
    workspace: workspace,
  }, nil
}

type FileBackend struct {
  logger.Logger
  *FileTaskLogger
  *FileTaskCanceler
  storage.Storage
  task *tes.Task
  workspace *Workspace
}

func (b *FileBackend) Task() *tes.Task {
  return b.task
}

func (b *FileBackend) Close() {
}

func (b *FileBackend) Executor(i int, d *tes.Executor) Executor {
  log := &FileExecutorLogger{
    client: b.client,
    taskID: b.task.Id,
    executor: i,
  }

  stdin, ierr := b.workspace.Reader(d.Stdin)
  stdout, oerr := b.workspace.Writer(d.Stdout)
  stderr, eerr := b.workspace.Writer(d.Stderr)

  if err := util.Check(ierr, oerr, eerr); err != nil {
    return nil, err
  }

  return &Docker{
    log,
    ImageName:       d.ImageName,
    Cmd:             d.Cmd,
    Volumes:         r.mapper.Volumes,
    Workdir:         d.Workdir,
    Ports:           d.Ports,
    ContainerName:   fmt.Sprintf("%s-%d", task.Id, i),
    RemoveContainer: r.conf.RemoveContainer,
    Environ:         d.Environ,
    Stdin: stdin,
    Stdout: stdout,
    Stderr: stderr,
  }, nil
}
