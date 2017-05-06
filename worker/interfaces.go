package worker

type ExecutorMetadata struct {
	HostIP string
	Ports  []*tes.Ports
}

type Executor interface {
	Run(context.Context) error
	Inspect(context.Context) ExecutorMetadata
  Close()
}

type TaskLogger interface {
	StartTime(t string)
	EndTime(t string)
	OutputFile(f string)
	Metadata(map[string]string)
  Running()
  Result(error)
  ExecutorExitCode(int, int)
  ExecutorPorts(int, []*tes.Ports)
  ExecutorHostIP(int, string)
  ExecutorStartTime(int, string)
  ExecutorEndTime(int, string)
  // TODO should these get access to the tes.Executor ?
  //      or even the whole task?
  ExecutorStdout(int) io.Writer
  ExecutorStderr(int) io.Writer
}

type TaskReader interface {
  Task() (*tes.Task, error)
  State() tes.State
}

type ExecutorFactory interface {
  Executor(int) (Executor, error)
}

type TaskRunner interface {
  RunTask(context.Context, *tes.Task)
}

type Backend interface {
  logger.Logger
	TaskLogger
  TaskReader
	storage.Storage
  ExecutorFactory
  TaskRunner

  Close()
}
