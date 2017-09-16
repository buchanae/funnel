package worker

import (
  "fmt"
  "io"
  "io/ioutil"
  "bufio"
  "os"
)

type Stdio struct {
  In  io.Reader
  Out io.Writer
  Err io.Writer
}

func OpenStdio(stdin, stdout, stderr string) (*Stdio, error) {
  s := Stdio{
    In: ioutil.NopCloser(bufio.NewReader(nil)),
    Out: ioutil.Discard,
    Err: ioutil.Discard,
  }
  var err error

	if stdin != "" {
    s.In, err = os.Open(stdin)
		if err != nil {
			return nil, fmt.Errorf("couldn't open stdin: %s", err)
		}
  }

	if stdout != "" {
		s.Out, err = os.Create(stdout)
		if err != nil {
			return nil, fmt.Errorf("couldn't open stdout: %s", err)
		}
	}

	if stderr != "" {
		s.Err, err = os.Create(stderr)
		if err != nil {
			return nil, fmt.Errorf("couldn't open stderr: %s", err)
		}
	}
	return &s, nil
}

/*
type ExecTailer struct {
  TaskLogger
  BufferSize int64
  FlushRate time.Duration
}
func (e *ExecTailer) Tail(ctx context.Context, i int, s Stdio) (*Stdio, error) {

  // Tail stdout
  stdout, err := newTailer(e.BufferSize, func(chunk string) {
    // When flushed, write an event to the task logger
    e.TaskLogger.AppendExecutorStdout(i, chunk)
  })
  if err != nil {
    return nil, err
  }

  // Tail stderr
  stderr, err := newTailer(e.BufferSize, func(chunk string) {
    // When flushed, write an event to the task logger
    e.TaskLogger.AppendExecutorStderr(i, chunk)
  })
  if err != nil {
    return nil, err
  }

  // Start a goroutine to periodically flush the tailers created above.
  go func() {
    ticker := time.NewTicker(e.FlushRate)
    defer ticker.Stop()
    defer stdout.Flush()
    defer stderr.Flush()

    for {
      select {
      case <-ctx.Done():
        return
      case <-ticker.C:
        stdout.Flush()
        stderr.Flush()
      }
    }
  }()

  return &s, nil
}

func OpenExecStdioWithTailer(e *tes.Executor, i int, t *ExecTailer) (*Stdio, error) {
  // Tail the stdout/err streams, periodically flushing events
  // to the task logger.
  stdio, err = tailer.Tail(ctx, i, stdio)
  if err != nil {
    return fmt.Errorf("couldn't prepare stdio tails: %s", err.Error())
  }
}
*/
