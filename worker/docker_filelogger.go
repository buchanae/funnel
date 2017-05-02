package worker

func NewFileBackend(conf config.Worker, taskID string) (*FileBackend, error) {
  workspace, werr := NewWorkspace(conf.WorkDir, taskID)
  NewTaskState()

  store, serr := storage.FromConfig(conf.Storage)

  if serr != nil {
    return nil, serr
  }

  return &FileBackend{
    Logger: log.WithFields("task", taskID),
    FileTaskLogger: &FileTaskLogger{workspace, taskID},
    FileTaskReader: &FileTaskReader{},
    Storage: store,
    taskID: taskID,
    workspace: workspace,
  }, nil
}


type FileTaskReader struct {}
func (b *FileTaskReader) Task() *tes.Task {
  return b.task
}
func (f *FileTaskReader) State() tes.State {
}


type FileBackend struct {
  logger.Logger
  *FileTaskLogger
  *FileTaskReader
  storage.Storage
  taskID string
  workspace *Workspace
}


func (b *FileBackend) Close() {
}

func (b *FileBackend) Executor(i int, d *tes.Executor) Executor {

  stdin, ierr := b.workspace.Reader(d.Stdin)
  stdout, oerr := b.workspace.Writer(d.Stdout)
  stderr, eerr := b.workspace.Writer(d.Stderr)

  if err := util.Check(ierr, oerr, eerr); err != nil {
    return nil, err
  }

  return &Docker{
    ImageName:       d.ImageName,
    Cmd:             d.Cmd,
    Volumes:         r.mapper.Volumes,
    Workdir:         d.Workdir,
    Ports:           d.Ports,
    ContainerName:   fmt.Sprintf("%s-%d", b.taskID, i),
    RemoveContainer: r.conf.RemoveContainer,
    Environ:         d.Environ,
    Stdin: stdin,
    Stdout: stdout,
    Stderr: stderr,
  }, nil
}
