package worker

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// LogStdio wraps executor stdout/err writers to write events to the given Logger.
// "i" is the index of the executor being logged.
func LogStdio(s Stdio, i int, log Logger) Stdio {
	s.Out = io.MultiWriter(s.Out, log.Stdout(i))
	s.Err = io.MultiWriter(s.Err, log.Stderr(i))
	return s
}

// Stdio makes it easier to reference a group of stdin/out/err handles.
type Stdio struct {
	In    io.Reader
	Out   io.Writer
	Err   io.Writer
	files []*os.File
}

func (s Stdio) Close() error {
	var errs []string

	// Close open file handles.
	for _, f := range s.files {
		err := f.Close()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// If there are errors, join them into a single error message.
	if errs != nil {
		return fmt.Errorf("error(s) closing stdio:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func OpenStdio(stdin, stdout, stderr string) (Stdio, error) {
	s := Stdio{
		In:  bytes.NewBuffer(nil),
		Out: ioutil.Discard,
		Err: ioutil.Discard,
	}

	if stdin != "" {
		f, err := os.Open(stdin)
		if err != nil {
			return s, fmt.Errorf("couldn't open stdin: %s", err)
		}
		s.files = append(s.files, f)
		s.In = f
	}

	if stdout != "" {
		f, err := os.Create(stdout)
		if err != nil {
			return s, fmt.Errorf("couldn't open stdout: %s", err)
		}
		s.files = append(s.files, f)
		s.Out = f
	}

	if stderr != "" {
		f, err := os.Create(stderr)
		if err != nil {
			return s, fmt.Errorf("couldn't open stderr: %s", err)
		}
		s.files = append(s.files, f)
		s.Err = f
	}
	return s, nil
}
