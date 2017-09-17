package worker

import (
	"context"
	"fmt"
	tl "github.com/ohsu-comp-bio/funnel/proto/tasklogger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"io"
	"time"
)

// TODO document behavior of slow consumer of task log updates

// RPCTask provides access to writing task logs over gRPC to the funnel server.
type RPCTask struct {
	Reader *RPCReader
	Logger *RPCLogger
}

func newRPCTask(conf rpc.Config, taskID string) (*RPCTask, error) {
	conn, err := rpc.Dial(conf)
	if err != nil {
		return nil, err
	}

	return &RPCTask{
		&RPCReader{tes.NewTaskServiceClient(conn), taskID},
		&RPCLogger{tl.NewTaskLoggerServiceClient(conn), taskID, conf.Timeout},
	}, nil
}

type RPCReader struct {
	client tes.TaskServiceClient
	taskID string
}

// Task returns the task descriptor.
func (r *RPCReader) Task() (*tes.Task, error) {
	return r.client.GetTask(context.Background(), &tes.GetTaskRequest{
		Id:   r.taskID,
		View: tes.TaskView_FULL,
	})
}

// State returns the current state of the task.
func (r *RPCReader) State() tes.State {
	t, _ := r.client.GetTask(context.Background(), &tes.GetTaskRequest{
		Id: r.taskID,
	})
	return t.GetState()
}

type RPCLogger struct {
	client        tl.TaskLoggerServiceClient
	taskID        string
	updateTimeout time.Duration
}

func (e RPCLogger) Debug(msg string, fields map[string]string) error {
  return nil
}
func (e RPCLogger) Info(msg string, fields map[string]string) error {
  return nil
}
func (e RPCLogger) Error(msg string, fields map[string]string) error {
  return nil
}

// SetState sets the state of the task.
func (r *RPCLogger) State(s tes.State) error {
	_, err := r.client.UpdateTaskState(context.Background(), &tl.UpdateTaskStateRequest{
		Id:    r.taskID,
		State: s,
	})
	return err
}

// StartTime updates the task's start time log.
func (r *RPCLogger) StartTime(t time.Time) error {
	return r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			StartTime: t.Format(time.RFC3339),
		},
	})
}

// EndTime updates the task's end time log.
func (r *RPCLogger) EndTime(t time.Time) error {
	return r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			EndTime: t.Format(time.RFC3339),
		},
	})
}

// Outputs updates the task's output file log.
func (r *RPCLogger) Outputs(f []*tes.OutputFileLog) error {
	return r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			Outputs: f,
		},
	})
}

// Metadata updates the task's metadata log.
func (r *RPCLogger) Metadata(m map[string]string) error {
	return r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			Metadata: m,
		},
	})
}

// ExecutorStartTime updates an executor's start time log.
func (r *RPCLogger) ExecutorStartTime(i int, t time.Time) error {
	return r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			StartTime: t.Format(time.RFC3339),
		},
	})
}

// ExecutorEndTime updates an executor's end time log.
func (r *RPCLogger) ExecutorEndTime(i int, t time.Time) error {
	return r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			EndTime: t.Format(time.RFC3339),
		},
	})
}

// ExitCode updates an executor's exit code log.
func (r *RPCLogger) ExitCode(i int, x int) error {
	return r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			ExitCode: int32(x),
		},
	})
}

// Ports updates an executor's ports log.
func (r *RPCLogger) Ports(i int, ports []*tes.Ports) error {
	return r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			Ports: ports,
		},
	})
}

// HostIP updates an executor's host IP log.
func (r *RPCLogger) HostIP(i int, ip string) error {
	return r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			HostIp: ip,
		},
	})
}

func (r *RPCLogger) Stdout(i int) io.Writer {
	return &stdoutWriter{r, i}
}
func (r *RPCLogger) Stderr(i int) io.Writer {
	return &stderrWriter{r, i}
}

type stdoutWriter struct {
	r *RPCLogger
	i int
}

func (s *stdoutWriter) Write(p []byte) (int, error) {
	err := s.r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   s.r.taskID,
		Step: int64(s.i),
		Log: &tes.ExecutorLog{
			Stdout: string(p),
		},
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type stderrWriter struct {
	r *RPCLogger
	i int
}

func (s *stderrWriter) Write(p []byte) (int, error) {
	err := s.r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   s.r.taskID,
		Step: int64(s.i),
		Log: &tes.ExecutorLog{
			Stderr: string(p),
		},
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (r *RPCLogger) updateExecutorLogs(up *tl.UpdateExecutorLogsRequest) error {
	ctx, cleanup := context.WithTimeout(context.Background(), r.updateTimeout)
	_, err := r.client.UpdateExecutorLogs(ctx, up)
	if err != nil {
		return fmt.Errorf("Couldn't update executor logs: %s", err)
	}
	cleanup()
	return err
}

func (r *RPCLogger) updateTaskLogs(up *tl.UpdateTaskLogsRequest) error {
	ctx, cleanup := context.WithTimeout(context.Background(), r.updateTimeout)
	_, err := r.client.UpdateTaskLogs(ctx, up)
	if err != nil {
		return fmt.Errorf("Couldn't update task logs: %s", err)
	}
	cleanup()
	return err
}
