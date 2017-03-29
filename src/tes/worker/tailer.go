package worker

import (
	"sync"
	"tes/util/ring"
)

func newTailer(size int64) (*tailer, error) {
	buf, err := ring.NewBuffer(size)
	if err != nil {
		return nil, err
	}
	return &tailer{buf: buf}, nil
}

type tailer struct {
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
	if t.buf.TotalWritten() > 100 {
		t.Flush()
	}
	return w, nil
}

func (t *tailer) Flush() string {
	t.mtx.Lock()
	t.mtx.Unlock()
	if t.buf.TotalWritten() > 0 {
    return t.buf.String()
		t.buf.Reset()
	}
  return ""
}
