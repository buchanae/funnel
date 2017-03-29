package worker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"tes/config"
	pbe "tes/ga4gh"
	"tes/logger"
	pbr "tes/server/proto"
	"tes/storage"
	"tes/util"
)

type ExecutorMeta struct {
  Ports
  IP
}

type Backend interface {
  Validate(*pbe.Job) error
  Storage(*pbe.Job) (*storage.Storage, error)
  Execute(context.Context) error
  Inspect(context.Context) ExecutorMeta
}


func (r *backend) Executor(*pbe.Job, ???) error {
    ID: fmt.Sprintf("%s-%d", job.JobID, i),

    exec := DockerExecutor{
      // TODO make RemoveContainer configurable
      RemoveContainer: true,
    }

    // Opens stdin/out/err files and updates those fields on "cmd".
    err := r.openStepLogs(s, d)
    if err != nil {
      stepLog.Error("Couldn't prepare log files", err)
      return err
    }
  }
}


// openLogs opens/creates the logs files for a step and updates those fields.
func (r *backend) openStepLogs(d *pbe.DockerExecutor) error {

	// Find the path for job stdin
	var err error
	if d.Stdin != "" {
		s.Cmd.Stdin, err = r.mapper.OpenHostFile(d.Stdin)
		if err != nil {
			return err
		}
	}

	// Create file for job stdout
	if d.Stdout != "" {
		s.Cmd.Stdout, err = CreateHostFile(d.Stdout)
		if err != nil {
			return err
		}
	}

	// Create file for job stderr
	if d.Stderr != "" {
		s.Cmd.Stderr, err = CreateHostFile(d.Stderr)
		if err != nil {
			return err
		}
	}
	return nil
}


func (r *backend) Validate(job *pbe.Job) error {
  err := validate.Task(job)
  if err != nil {
    return err
  }

  // Validate volumes
	for _, vol := range job.Task.Volumes {
    err := validate.Volume(vol)
    if err != nil {
      return err
    }

    // TODO find good home for directory/file init
    err := util.EnsureDir(hostPath)
    if err != nil {
      return err
    }
  }

  // Validate outputs
	for _, output := range m.Outputs {

    // Create the file if needed, as per the TES spec
    // TODO find a good home for directory prep
    if output.Create {
      err := util.EnsureFile(p, output.Class)
      if err != nil {
        return err
      }
    }
	}
	return nil
}







type jobStorage struct {
  store storage.Storage
}

func (j *jobStorage) hostPath(p string) string {
}

func (j *jobStorage) Get(ctx context.Context, url, path, class string, readonly bool) error {
  mapped := j.hostPath(path)
  return j.store.Get(ctx, url, mapped, class, readonly)
}

func (j *jobStorage) Put(ctx context.Context, url, path, class string) error {
  mapped := r.hostPath(path)
  j.fixLinks(mapped)
  return j.store.Put(ctx, url, mapped, class)
}

func (j *jobStorage) Supports(url, path, class string) bool {
  mapped := r.hostPath(job, input.Path)
  return isSubpath(mapped, j.baseDir)
}

// fixLinks walks the output paths, fixing cases where a symlink is
// broken because it's pointing to a path inside a container volume.
func (j *jobStorage) fixLinks(basepath string) {
  // TODO what happens when basepath is a file?
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
        dst, err := os.Readlink(p)
        if err != nil {
          return nil
        }

        mapped, err := hostPath(dst)
        if err != nil {
          return nil
        }

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
