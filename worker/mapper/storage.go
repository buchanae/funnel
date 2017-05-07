package mapper

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
)

func MapStorage(base string, store storage.Backend) storage.Backend {
  return &mappedStorage{base, store}
}

type mappedStorage struct {
  base string
  storage.Backend
}

func (m *mappedStorage) Get(ctx context.Context, url, path string, ft tes.FileType) error {
  // Map the path. If there's an error, return it.
  p, err := MapPath(m.base, path)
  if err != nil {
    return err
  }
  // No error, so continue with Get on mapped path.
  return m.Backend.Get(ctx, url, p, ft)
}

func (m *mappedStorage) Put(ctx context.Context, url, path string, ft tes.FileType) error {
  // Map the path. If there's an error, return it.
  p, err := MapPath(m.base, path)
  if err != nil {
    return err
  }

  util.FixSymlinks(p, func(p string) (string, error) {
    return MapPath(m.base, p)
  })

  // No error, so continue with Put on mapped path.
  return m.Backend.Put(ctx, url, p, ft)
}
