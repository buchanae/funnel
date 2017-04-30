package worker

type ExecutorMetadata struct {
	HostIP string
	Ports  []*tes.Ports
}

type Executor interface {
  ExecutorLogger

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
  Close()
}

type TaskLogger interface {
	StartTime(t string)
	EndTime(t string)
	OutputFile(f string)
	Metadata(map[string]string)
  Running()
  Result(error)
  Close()
}

type TaskCanceler interface {
  WatchForCancel(context.Context) context.Context
}

// TODO document behavior of slow consumer of task log updates
type Backend interface {
  logger.Logger
	TaskLogger
  TaskCanceler
	storage.Storage

  Task(id string)
  Executor(int, *tes.Executor) (Executor, error)
  WithContext(context.Context) context.Context
  Close()
}


      // TODO move to storage wrapper
			//r.fixLinks(output.Path)
