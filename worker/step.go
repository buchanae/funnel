package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"io"
	"time"
)

type stepRunner struct {
	TaskID     string
	Conf       config.Worker
	Num        int
  Exec       ExecutorBackend
	Log        logger.Logger
	TaskLogger TaskLogger
	IP         string
}

func (s *stepRunner) Run(ctx context.Context) error {
	log.Debug("Running step", "taskID", s.TaskID, "stepNum", s.Num)

	// Send update for host IP address.
	s.TaskLogger.ExecutorStartTime(s.Num, time.Now())
	s.TaskLogger.ExecutorHostIP(s.Num, s.IP)

	// subctx helps ensure that these goroutines are cleaned up,
	// even when the task is canceled.
	subctx, cleanup := context.WithCancel(ctx)
	defer cleanup()

	done := make(chan int64, 1)

	// tailLogs modifies the cmd Stdout/err fields, so should be called before Run.
	stdout, stderr := s.logTails()
	defer stdout.Flush()
	defer stderr.Flush()

	ticker := time.NewTicker(s.Conf.LogUpdateRate)
	defer ticker.Stop()

	go func() {
    code, err := s.Exec.Run(ctx)
    if err != nil {
      code = int64(-999)
    }
		done <- code
	}()
	go s.inspectContainer(subctx)

	for {
		select {
		case <-ctx.Done():
			// Likely the task was canceled.
			s.Exec.Stop(ctx)
			s.TaskLogger.ExecutorEndTime(s.Num, time.Now())
			return ctx.Err()

		case <-ticker.C:
			stdout.Flush()
			stderr.Flush()

		case exitCode := <-done:
			s.TaskLogger.ExecutorEndTime(s.Num, time.Now())
			s.TaskLogger.ExecutorExitCode(s.Num, exitCode)
			return result
		}
	}
}

func (s *stepRunner) logTails() (*tailer, *tailer) {
	stdout, _ := newTailer(s.Conf.LogTailSize, func(c string) {
		s.TaskLogger.AppendExecutorStdout(s.Num, c)
	})
	stderr, _ := newTailer(s.Conf.LogTailSize, func(c string) {
		s.TaskLogger.AppendExecutorStderr(s.Num, c)
	})
	if s.Cmd.Stdout != nil {
		s.Cmd.Stdout = io.MultiWriter(s.Cmd.Stdout, stdout)
	}
	if s.Cmd.Stderr != nil {
		s.Cmd.Stderr = io.MultiWriter(s.Cmd.Stderr, stderr)
	}
	return stdout, stderr
}

// inspectContainer calls Inspect on the DockerCmd, and sends an update with the results.
func (s *stepRunner) inspectContainer(ctx context.Context) {
	ports, err := s.Exec.Inspect(ctx)
	if err != nil {
		s.Log.Error("Error inspecting container", err)
		return
	}
	s.TaskLogger.ExecutorPorts(s.Num, ports)
}
