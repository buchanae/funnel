package mapper

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
  "github.com/golang/protobuf/proto"
	"path/filepath"
  "strings"
)

func MapTask(base string, task *tes.Task) (*tes.Task, error) {
  var err error

  // Create the base task working directory.
  _, err = prepPath(base, ".", tes.FileType_DIRECTORY)

  // Clone the task into a mapped task which will be returned.
  mapped := proto.Clone(task).(*tes.Task)

  // Create and map task inputs.
	for i, input := range task.Inputs {
    if err == nil {
      mapped.Inputs[i].Path, err = prepPath(base, input.Path, input.Type)
    }
  }

  // Create and map task volumes.
	for i, vol := range task.Volumes {
    if err == nil {
      mapped.Volumes[i], err = prepPath(base, vol, tes.FileType_DIRECTORY)
    }
  }

  // Create and map task outputs.
	for i, output := range task.Outputs {
    if err == nil {
      mapped.Outputs[i].Path, err = prepPath(base, output.Path, output.Type)
    }
  }

  if err != nil {
    return nil, err
  }
  return mapped, nil
}

func MapPath(base, src string) (string, error) {
	p := filepath.Join(base, src)
	p, err := filepath.Abs(p)

  if err != nil {
    return "", err
  }

  // Path must be a subpath of
	if !isSubpath(p, base) {
		return "", fmt.Errorf("Invalid path: %s is not a valid subpath of %s", src, base)
	}
  return p, nil
}

// prepPath calls mapPath and also ensures the directories in the path exist.
// "err" is passed as an argument so that there's less error checking in
// MapTask().
func prepPath(base, src string, t tes.FileType) (string, error) {
  p, err := MapPath(base, src)
  if err != nil {
    return "", err
  }

  // Ensure the path directory exists.
  d := p
  if t == tes.FileType_FILE {
    d = filepath.Dir(d)
  }
  eerr := util.EnsureDir(p)
	return p, eerr
}

// isSubpath returns true if the given path "p" is a subpath of "base".
// TODO ensure this is the best way to check
func isSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}
