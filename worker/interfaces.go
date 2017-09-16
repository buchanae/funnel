package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
	"time"
)

// Worker is a type which runs a task.
type Worker interface {
	Run(context.Context)
}

// TaskReader is a type which reads and writes task information
// during task execution.
type TaskReader interface {
	Task() (*tes.Task, error)
	State() tes.State
}

// Logger provides write access to a worker's logs.
type Logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})

	State(tes.State)
	StartTime(t time.Time)
	EndTime(t time.Time)
	Outputs(o []*tes.OutputFileLog)
	Metadata(m map[string]string)

	ExecutorExitCode(i int, code int)
	ExecutorPorts(i int, ports []*tes.Ports)
	ExecutorHostIP(i int, ip string)
	ExecutorStartTime(i int, t time.Time)
	ExecutorEndTime(i int, t time.Time)

	ExecutorStdout(i int) io.Writer
	ExecutorStderr(i int) io.Writer
}
