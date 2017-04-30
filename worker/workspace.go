
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
	p, perr := mapper.HostPath(src)
	if perr != nil {
		return nil, perr
	}
	err := util.EnsurePath(p)
	if err != nil {
		return nil, err
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
	p, perr := mapper.HostPath(src)
	if perr != nil {
		return nil, perr
	}
	f, oerr := os.Open(p)
	if oerr != nil {
		return nil, oerr
	}
	return f, nil
}

// HostPath returns a mapped path.
//
// The path is concatenated to the mapper's base dir.
// e.g. If the mapper is configured with a base dir of "/tmp/mapped_files", then
// mapper.HostPath("/home/ubuntu/myfile") will return "/tmp/mapped_files/home/ubuntu/myfile".
//
// The mapped path is required to be a subpath of the mapper's base directory.
// e.g. mapper.HostPath("../../foo") should fail with an error.
func (mapper *FileMapper) HostPath(src string) (string, error) {
	p := path.Join(mapper.dir, src)
	p = path.Clean(p)
	if !mapper.IsSubpath(p, mapper.dir) {
		return "", fmt.Errorf("Invalid path: %s is not a valid subpath of %s", p, mapper.dir)
	}
	return p, nil
}

// IsSubpath returns true if the given path "p" is a subpath of "base".
func (mapper *FileMapper) IsSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}
