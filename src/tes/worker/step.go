package worker

import (
	"context"
	"io"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	pbr "tes/server/proto"
	"tes/util/ring"
	"time"
)

type stepRunner struct {
	JobID   string
	Conf    config.Worker
	Num     int
	Cmd     *DockerCmd
	Log     logger.Logger
	Updates updateChan
	IP      string
}

func (s *stepRunner) Run(ctx context.Context) error {
	log.Debug("Running step", "jobID", s.JobID, "stepNum", s.Num)

	// Send update for host IP address.
	s.update(&pbe.JobLog{
		HostIP: s.IP,
	})

	// subctx helps ensure that these goroutines are cleaned up,
	// even when the job is canceled.
	subctx, cleanup := context.WithCancel(ctx)
	defer cleanup()

	// tailLogs modifies the cmd Stdout/err fields, so should be called before Run.
	done := make(chan error, 1)
	go s.tailLogs(subctx)
	go func() {
		done <- s.Cmd.Run()
	}()
	go s.inspectContainer(subctx)

	select {
	case <-ctx.Done():
		// Likely the job was canceled.
		s.Cmd.Stop()
		return ctx.Err()
	case result := <-done:
		s.update(&pbe.JobLog{
			ExitCode: getExitCode(result),
		})
		return result
	}
}

// tailLogs starts a ticker and loop that sends updates of step stdout/err logs.
// This modifies "dcmd" to replace the Stdout/err fields with wrapped io.Writers.
func (s *stepRunner) tailLogs(ctx context.Context) {
	ticker := time.NewTicker(s.Conf.LogUpdateRate)
	defer ticker.Stop()

	var err error
	var stdoutTail, stderrTail *ring.Buffer
	stdoutTail, err = ring.NewBuffer(s.Conf.LogTailSize)
	stderrTail, err = ring.NewBuffer(s.Conf.LogTailSize)
	if err != nil {
		s.Log.Error("Can't create stdout/err tail buffer.", err)
		return
	}

	if s.Cmd.Stdout != nil {
		s.Cmd.Stdout = io.MultiWriter(stdoutTail, s.Cmd.Stdout)
	}
	if s.Cmd.Stderr != nil {
		s.Cmd.Stderr = io.MultiWriter(stderrTail, s.Cmd.Stderr)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.update(&pbe.JobLog{
				Stdout: string(stdoutTail.Bytes()),
				Stderr: string(stderrTail.Bytes()),
			})
			stdoutTail.Reset()
			stderrTail.Reset()
		}
	}
}

// inspectContainer calls Inspect on the DockerCmd, and sends an update with the results.
func (s *stepRunner) inspectContainer(ctx context.Context) {
	ports, err := s.Cmd.Inspect(ctx)
	if err != nil {
		s.Log.Error("Error inspecting container", err)
		return
	}

	s.update(&pbe.JobLog{
		Ports: ports,
	})
}

// update sends an update of the JobLog of the currently running step.
// Used to update stdout/err logs, port mapping, etc.
func (s *stepRunner) update(log *pbe.JobLog) {
	s.Updates <- &pbr.UpdateStatusRequest{
		Id:   s.JobID,
		Step: int64(s.Num),
		Log:  log,
	}
}