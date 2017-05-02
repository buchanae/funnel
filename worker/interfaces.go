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

type ExecutorLogger interface {
	ExitCode(int)
	Ports([]*tes.Ports)
	HostIP(string)
	StartTime(t string)
	EndTime(t string)
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
}

type TaskReader interface {
  Task() *tes.Task
  State() tes.State
}

// TODO document behavior of slow consumer of task log updates
type Backend interface {
  logger.Logger
	TaskLogger
  TaskReader
	storage.Storage

  Executor(int, *tes.Executor) (Executor, error)
  Close()
}


      // TODO move to storage wrapper
			//r.fixLinks(output.Path)
