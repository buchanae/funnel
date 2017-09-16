package worker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ExecutorStdioEvents(s *Stdio, i int, log Logger) *Stdio {
	s.Out = io.MultiWriter(s.Out, log.ExecutorStdout(i))
	s.Err = io.MultiWriter(s.Err, log.ExecutorStderr(i))
	return s
}

type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func NewStdio(stdin, stdout, stderr string) (*Stdio, error) {
	s := Stdio{
		In:  ioutil.NopCloser(bufio.NewReader(nil)),
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
