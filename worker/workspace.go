
func NewWorkspace(root string, taskID string) (*Workspace, error) {
	baseDir := path.Join(conf.WorkDir, t.Task.Id)
	dir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
  util.EnsureDir(dir)
}

type Workspace struct {}

// CreateHostFile creates a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Create
//
// If the path can't be mapped or the file can't be created, an error is returned.
func (w *Workspace) Writer(p string) (io.Writer, error) {
  if p == "" {
    return nil, nil
  }
	p, perr := w.Path(src)
	if perr != nil {
		return nil, perr
	}
	f, oerr := os.Create(p)
	if oerr != nil {
		return nil, oerr
	}
	return f, nil
}

// OpenHostFile opens a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Open
//
// If the path can't be mapped or the file can't be opened, an error is returned.
func (w *Workspace) Reader(p string) (io.Reader, error) {
  if p == "" {
    return nil, nil
  }
	p, perr := w.Path(src)
	if perr != nil {
		return nil, perr
	}
	f, oerr := os.Open(p)
	if oerr != nil {
		return nil, oerr
	}
	return f, nil
}

// Path returns a mapped path.
//
// The path is concatenated to the w's base dir.
// e.g. If the w is configured with a base dir of "/tmp/mapped_files", then
// w.Path("/home/ubuntu/myfile") will return "/tmp/mapped_files/home/ubuntu/myfile".
//
// The mapped path is required to be a subpath of the w's base directory.
// e.g. w.Path("../../foo") should fail with an error.
func (w *Workspace) Path(src string) (string, error) {
	p := path.Join(w.dir, src)
	p = path.Clean(p)
  // Path must be a subpath of
  // TODO ensure this is the best way to check
	if !isSubpath(p, w.dir) {
		return "", fmt.Errorf("Invalid path: %s is not a valid subpath of %s", p, w.dir)
	}
	return p, nil
}

// isSubpath returns true if the given path "p" is a subpath of "base".
func isSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}

func (w *Workspace) MapVolumes(in []Volume) []Volume {
  var volumes []Volume
  for _, v := range in {
  }
  return volumes
}

func (w *Workspace) MapStorage(s storage.Storage) storage.Storage {
  // TODO wrap Storage so that paths are mapped to workspace
}
