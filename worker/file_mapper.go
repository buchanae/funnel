package worker

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"path/filepath"
)

func MapTaskFiles(base string, task *tes.Task) (*tes.Task, error) {
  h := helper{
    base: filepath.Join(base, task.Id),
  }
  // Create the base task working directory.
  h.mapPath("/", true)

  // Clone the task into a mapped task which will be returned.
  mapped := proto.Clone(task).(*tes.Task)

	for i, input := range task.Inputs {
    mapped.Inputs[i].Path = h.mapPath(input.Path, input.Type)
  }

	for i, vol := range task.Volumes {
    mapped.Volumes[i] = h.mapPath(vol, tes.FileType_DIRECTORY)
  }

	for i, output := range task.Outputs {
    mapped.Outputs[i].Path = h.mapPath(output.Path, output.Type)
  }

  if h.err != nil {
    return nil, h.err
  }
  return mapped, nil
}


type helper {
  base string
  err error
}
func (h *helper) mapPath(src string, t tes.FileType) string {
  if h.err != nil {
    return ""
  }

	p := filepath.Join(h.base, src)
	p = filepath.Abs(p)

  // Path must be a subpath of
  // TODO ensure this is the best way to check
	if !isSubpath(p, h.base) {
    h.err = fmt.Errorf("Invalid path: %s is not a valid subpath of %s", p, h.base)
		return ""
	}

  // Ensure the path directory exists.
  d := p
  if t == tes.FileType_FILE {
    d = path.Dir(d)
  }
  f.err = util.EnsureDir(p)

	return p
}

// isSubpath returns true if the given path "p" is a subpath of "base".
func isSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}
