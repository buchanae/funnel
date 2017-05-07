package util

import (
  "io"
  "io/ioutil"
  "os"
)

func WriterOrDiscard(p string) io.Writer {
  if p == "" {
    return ioutil.Discard
  }
	f, err := os.Create(p)
	if err != nil {
		return ioutil.Discard
	}
	return f
}

func ReaderOrEmpty(p string) io.Reader {
  if p == "" {
    return nil
  }
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	return f
}
