package worker

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"path/filepath"
)

func MapTaskFiles(base string, task *tes.Task) (*tes.Task, error) {
  var err *maperr
  base := filepath.Join(base, task.Id),

  // Create the base task working directory.
  h.mapPath("/", true)

  // Clone the task into a mapped task which will be returned.
  mapped := proto.Clone(task).(*tes.Task)

	for i, input := range task.Inputs {
    mapped.Inputs[i].Path = h.mapPath(input.Path, input.Type, err)
  }

	for i, vol := range task.Volumes {
    mapped.Volumes[i] = h.mapPath(vol, tes.FileType_DIRECTORY, err)
  }

	for i, output := range task.Outputs {
    mapped.Outputs[i].Path = h.mapPath(output.Path, output.Type, err)
  }

  if h.err != nil {
    return nil, h.err
  }
  return mapped, nil
}

func MapTaskStorage(base string, store storage.Storage) storage.Storage {
  return &mappedStorage{base, store}
}

// TODO maybe another argument for Storage.DownloadInputs()?
//      other use cases:
//      - fix links
//      - parallel download
type mappedStorage struct {
  storage.Storage
  base string
}
func (m *MappedStorage) Get(url, path, filetype string) error {
  // Map the path. If there's an error, return it.
  var err *maperr
  m, err := mapPath(path, filetype, err)
  if err != nil {
    return err
  }
  // No error, so continue with Get on mapped path.
  return m.Storage.Get(url, m, filetype)
}

func (m *mappedStorage) Put(url, path, filetype string) error {
  // Map the path. If there's an error, return it.
  var err *maperr
  m, err := mapPath(path, filetype, err)
  if err != nil {
    return err
  }
  // No error, so continue with Put on mapped path.
  return m.Storage.Put(url, m, filetype)
}

func mapPath(base, src string, t tes.FileType, err *maperr) string {
  // If there was a previous error, don't do anything.
  if err != nil {
    return ""
  }

	p := filepath.Join(base, src)
	p = filepath.Abs(p)

  // Path must be a subpath of
  // TODO ensure this is the best way to check
	if !isSubpath(p, base) {
    *err = maperr{p, base}
		return ""
	}

  // Ensure the path directory exists.
  d := p
  if t == tes.FileType_FILE {
    d = path.Dir(d)
  }
  ederr := util.EnsureDir(p)
  if ederr != nil {
    *err = maperr{err: ederr}
  }

	return p
}

// isSubpath returns true if the given path "p" is a subpath of "base".
func isSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}

type maperr struct {
  path string
  base string
  err error
}
func (m *maperr) Error() string {
  if m.err != nil {
    return m.err.Error()
  }
  return fmt.Sprintf("Invalid path: %s is not a valid subpath of %s", m.path, m.base)
}
