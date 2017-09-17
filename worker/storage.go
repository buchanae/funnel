package worker

import (
	"context"
  "fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/storage"
	"os"
  "io/ioutil"
	"path/filepath"
)

func Upload(ctx context.Context, tp []*tes.TaskParameter, s storage.Storage) ([]*tes.OutputFileLog, error) {
	var outputs []*tes.OutputFileLog

	for _, output := range tp {
		var out []*tes.OutputFileLog
		out, err := s.Put(ctx, output.Url, output.Path, output.Type)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, out...)
	}

	return outputs, nil
}

// FixLinks walks the output paths, fixing cases where a symlink is
// broken because it's pointing to a path inside a container volume.
//
// "mapper" is used to map the broken link (e.g. containerized path)
// to a correct path (e.g. host path).
func FixLinks(basepath string, mapper func(string) string) {
	filepath.Walk(basepath, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			// There's an error, so be safe and give up on this file
			return nil
		}

		// Only bother to check symlinks
		if f.Mode()&os.ModeSymlink != 0 {
			// Test if the file can be opened because it doesn't exist
			fh, rerr := os.Open(p)
			fh.Close()

			if rerr != nil && os.IsNotExist(rerr) {

				// Get symlink source path
				src, err := os.Readlink(p)
				if err != nil {
					return nil
				}
				// Map symlink source (possible container path) to host path
				mapped := mapper(src)

				// Check whether the mapped path exists
				fh, err := os.Open(mapped)
				fh.Close()

				// If the mapped path exists, fix the symlink
				if err == nil {
					err := os.Remove(p)
					if err != nil {
						return nil
					}
					os.Symlink(mapped, p)
				}
			}
		}
		return nil
	})
}

func LogUpload(ctx context.Context, out []*tes.TaskParameter, s storage.Storage, l Logger) error {
	outputs, err := Upload(ctx, out, s)
	if err != nil {
		return err
	}
	l.Outputs(outputs)
	return nil
}

func Download(ctx context.Context, in []*tes.TaskParameter, s storage.Storage) error {
	for _, input := range in {

    // If 'contents' field is set create the file
    if input.Contents != "" {
      err := ioutil.WriteFile(input.Path, []byte(input.Contents), 0775)
      if err != nil {
        return fmt.Errorf("Error writing contents of input to file %v", err)
      }

    } else {
      err := s.Get(ctx, input.Url, input.Path, input.Type)
      if err != nil {
        return err
      }
    }
	}
	return nil
}

func ValidateStorage(s storage.Storage, inputs, outputs []*tes.TaskParameter) error {
	for _, input := range inputs {
    if input.Contents != "" {
      continue
    }
		if !s.Supports(input.Url, input.Path, input.Type) {
			return fmt.Errorf("input not supported by storage: %v", input)
		}
	}
	for _, output := range outputs {
		if !s.Supports(output.Url, output.Path, output.Type) {
			return fmt.Errorf("output not supported by storage: %v", output)
		}
	}
	return nil
}
