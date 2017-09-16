package worker

import (
	"github.com/ohsu-comp-bio/funnel/util/ring"
	"sync"
)

func newTailer(size int64, out func(string)) (*tailer, error) {
	buf, err := ring.NewBuffer(size)
	if err != nil {
		return nil, err
	}
	return &tailer{buf: buf, out: out}, nil
}

type tailer struct {
	out func(string)
	buf *ring.Buffer
	mtx sync.Mutex
}

func (t *tailer) Write(b []byte) (int, error) {
	t.mtx.Lock()
	t.mtx.Unlock()
	w, err := t.buf.Write(b)
	if err != nil {
		return w, err
	}
	// This is suspicious. I think this helps flush the first small amount
	// of content written?
	if t.buf.TotalWritten() > 100 {
		t.Flush()
	}
	return w, nil
}

func (t *tailer) Flush() {
	t.mtx.Lock()
	t.mtx.Unlock()
	if t.buf.TotalWritten() > 0 {
		t.out(t.buf.String())
		t.buf.Reset()
	}
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
