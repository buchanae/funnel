package worker

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"os"
	"path"
	"path/filepath"
	"strings"
	pbe "tes/ga4gh"
	"tes/util"
)

// Volume represents a volume mounted into a docker container.
// This includes a HostPath, the path on the host file system,
// and a ContainerPath, the path on the container file system,
// and whether the volume is read-only.
type Volume struct {
	// The path in tes worker.
	HostPath string
	// The path in Docker.
	ContainerPath string
	Readonly      bool
}


// HostPath returns a mapped path.
//
// The path is concatenated to the mapper's base dir.
// e.g. If the mapper is configured with a base dir of "/tmp/mapped_files", then
// mapper.HostPath("/home/ubuntu/myfile") will return "/tmp/mapped_files/home/ubuntu/myfile".
//
// The mapped path is required to be a subpath of the mapper's base directory.
// e.g. mapper.HostPath("../../foo") should fail with an error.
func hostPath(base, src string) string {
	p := path.Join(base, src)
	p = path.Clean(p)
	if !isSubpath(p, base) {
    return base
	}
	return p
}

// CreateHostFile creates a file on the host file system at a mapped path.
// "src" is an unmapped path. This function will handle mapping the path.
//
// This function calls os.Create
//
// If the path can't be mapped or the file can't be created, an error is returned.
func CreateHostFile(src string) (*os.File, error) {
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

// isSubpath returns true if the given path "p" is a subpath of "base".
func isSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base) && p != base
}

// findVolume finds the volume that contains the given input/output parameter..
func findVolume(p *pbe.TaskParamter, task *pbe.Task) *pbe.Volume {
	for _, vol := range task.Volumes {
		if isSubpath(p.Path, vol.MountPoint) {
			return vol
		}
	}
	return nil
}
