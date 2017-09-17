package worker

import (
	"bytes"
	"github.com/ohsu-comp-bio/funnel/events"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Test events generated by executor stdio io.Reader/Writer interfaces
func TestLogStdio(t *testing.T) {
	col := events.Collector{}
	log := NewEventLogger("task-id", 0, &col)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	a := Stdio{Out: stdout, Err: stderr}
	b := LogStdio(a, 1, log)

	if len(col) != 0 {
		t.Error("expected zero events: %d", len(col))
	}

	table := [][2]string{
		// lines are pair of [stdout, stderr] to write
		{"a", "b"},
		{"out foo bar baz", "err foo bar baz"},
	}

	for i, x := range table {
		// Write to stdio
		b.Out.Write([]byte(x[0]))
		b.Err.Write([]byte(x[1]))

		// Check generated events
		count := (i + 1) * 2
		if len(col) != count {
			t.Fatalf("expected %d events but got %d", count, len(col))
		}

		outEvent := col[len(col)-2]
		errEvent := col[len(col)-1]

		if outEvent.Type != events.Type_STDOUT {
			t.Error("expected stdout event")
		} else if outEvent.Stdout != x[0] {
			t.Errorf("expected stdout event '%s' but got '%s'", x[0], outEvent.Stdout)
		}

		if errEvent.Type != events.Type_STDERR {
			t.Error("expected stderr event")
		} else if errEvent.Stderr != x[1] {
			t.Errorf("expected stderr event '%s' but got '%s'", x[1], errEvent.Stderr)
		}

		if stdout.String() != x[0] {
			t.Errorf("expected stdout buffer '%s' but got '%s'", x[0], stdout.String())
		}
		if stderr.String() != x[1] {
			t.Errorf("expected stderr buffer '%s' but got '%s'", x[1], stdout.String())
		}
		// Reset buffer so it's easier to check the next iteration
		stdout.Reset()
		stderr.Reset()
	}
}

// Test that empty paths result in noop reader/writers.
func TestNoopStdio(t *testing.T) {
	s, err := OpenStdio("", "", "")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(s.In)
	if len(b) > 0 {
		t.Error("expected empty bytes", b)
	}
	if err != nil {
		t.Error(err)
	}

	c := bytes.NewBufferString("foo")
	n, err := io.Copy(s.Out, c)
	if err != nil {
		t.Error(err)
	}
	if n != int64(3) {
		t.Error("expected all bytes to be written", n)
	}

	c = bytes.NewBufferString("foo")
	n, err = io.Copy(s.Err, c)
	if err != nil {
		t.Error(err)
	}
	if n != int64(3) {
		t.Error("expected all bytes to be written", n)
	}
}

func TestStdio(t *testing.T) {
	dir, err := ioutil.TempDir("", "funnel-worker-stdio-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	stdin := filepath.Join(dir, "stdin")
	stdout := filepath.Join(dir, "stdout")
	stderr := filepath.Join(dir, "stderr")

	ioutil.WriteFile(stdin, []byte("foo"), os.ModePerm)

	s, err := OpenStdio(stdin, stdout, stderr)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(s.In)
	if err != nil {
		t.Error(err)
	}
	if string(b) != "foo" {
		t.Error("expected stdin content 'foo' but got %s", string(b))
	}

	s.Out.Write([]byte("stdout foo"))
	s.Err.Write([]byte("stderr foo"))

	b, err = ioutil.ReadFile(stdout)
	if err != nil {
		t.Error(err)
	}
	if string(b) != "stdout foo" {
		t.Error("expected stdin content 'stdout foo' but got %s", string(b))
	}

	b, err = ioutil.ReadFile(stderr)
	if err != nil {
		t.Error(err)
	}
	if string(b) != "stderr foo" {
		t.Error("expected stdin content 'stderr foo' but got %s", string(b))
	}
}
