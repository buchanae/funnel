package util

import (
	"os"
	"path"
	"syscall"
)

const DefaultMode = os.ModePerm

// exists returns whether the given file or directory exists or not
func exists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// EnsureDir ensures a directory exists.
func EnsureDir(p string, mode os.FileMode) error {
	e, err := exists(p)
	if err != nil {
		return err
	}
	if !e {
		syscall.Umask(0000)
		err := os.MkdirAll(p, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsurePath ensures a directory exists, given a file path. This calls path.Dir(p)
// TODO probably just remove this
func EnsurePath(p string) error {
	dir := path.Dir(p)
	return EnsureDir(dir, DefaultMode)
}
