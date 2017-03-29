package storage

// NOTE!
// It's important that Storage instances be immutable!
// We don't want storage authentication to be accidentally shared between jobs.
// If they are mutable, there's more chance that storage config can leak
// between separate processes.

import (
	"context"
	"fmt"
	"tes/config"
)

const (
	// File represents the file type
	File string = "File"
	// Directory represents the directory type
	Directory = "Directory"
)

// Storage provides an interface for a storage backend.
// New storage backends must support this interface.
type Storage interface {
	Get(ctx context.Context, url string, path string, class string, readonly bool) error
	Put(ctx context.Context, url string, path string, class string) error
	// Determines whether this backends supports the given request (url/path/class).
	// A backend normally uses this to match the url prefix (e.g. "s3://")
	// TODO would it be useful if this included the request type (Get/Put)?
	Supports(url string, path string, class string) bool
}

// Storage provides a client for accessing multiple storage systems,
// i.e. for downloading/uploading job files from S3, GS, local disk, etc.
//
// For a given storage url, the storage backend is usually determined by the url prefix,
// e.g. "s3://my-bucket/file" will access the S3 backend.
type storage struct {
	backends []Backend
}

// Get downloads a file from a storage system at the given "url".
// The file is downloaded to the given local "path".
// "class" is either "File" or "Directory".
func (storage Storage) Get(ctx context.Context, url, path, class string, readonly bool) error {
	backend, err := storage.findBackend(url, path, class)
	if err != nil {
		return err
	}
	return backend.Get(ctx, url, path, class, readonly)
}

// Put uploads a file to a storage system at the given "url".
// The file is uploaded from the given local "path".
// "class" is either "File" or "Directory".
func (storage Storage) Put(ctx context.Context, url, path, class string) error {
	backend, err := storage.findBackend(url, path, class)
	if err != nil {
		return err
	}
	return backend.Put(ctx, url, path, class)
}

// Supports indicates whether the storage supports the given request.
func (storage Storage) Supports(url string, path string, class string) bool {
	b, _ := storage.findBackend(url, path, class)
	return b != nil
}

// findBackend tries to find a backend that matches the given url/path/class.
// This is how a url gets matched to a backend, for example by the url prefix "s3://".
func (storage Storage) findBackend(url string, path string, class string) (Backend, error) {
	for _, backend := range storage.backends {
		if backend.Supports(url, path, class) {
			return backend, nil
		}
	}
	return nil, fmt.Errorf("Could not find matching storage system for %s", url)
}

// FromConfig returns a new Storage instance with the given backend configurations.
func FromConfig(conf []*config.StorageConfig) (*Storage, error) {
  backends := []Backend

	for _, conf := range r.conf.Storage {
    if conf.Local.Valid() {
      local, err := NewLocalBackend(conf.Local)
      if err != nil {
        return nil, err
      }
      backends = append(backends, local)
    }

    if conf.S3.Valid() {
      s3, err := NewS3Backend(conf.S3)
      if err != nil {
        return nil, err
      }
      backends = append(backends, s3)
    }

    if conf.GS.Valid() {
      gs, err := NewGSBackend(conf.GS)
      if err != nil {
        return nil, err
      }
      backends = append(backends, gs)
    }
	}

  // If no valid backends were found, return an error.
  if len(backends) == 0 {
    return nil, fmt.Error("No valid storage backends could be configured.")
  }
  return &Storage{backends}, nil
}
