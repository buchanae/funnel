package util

// CreateHostFile creates a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Create
//
// If the path can't be mapped or the file can't be created, an error is returned.
func (w *FileMapper) FileWriter(p string) (io.Writer, error) {
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
func (w *FileMapper) FileReader(p string) (io.Reader, error) {
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
