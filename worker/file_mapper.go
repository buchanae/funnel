package worker

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"path"
	"path/filepath"
	"strings"
)

// FileMapper is responsible for mapping paths into a working directory on the
// worker's host file system.
//
// Every task needs it's own directory to work in. When a file is downloaded for
// a task, it needs to be stored in the task's working directory. Similar for task
// outputs, uploads, stdin/out/err, etc. FileMapper helps the worker engine
// manage all these paths.
type FileMapper struct {
	Volumes   []Volume
	Inputs    []*tes.TaskParameter
	Outputs   []*tes.TaskParameter
	Executors []*tes.Executor
	dir       string
}

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

// NewFileMapper returns a new FileMapper, which maps files into the given
// base directory.
func NewFileMapper(dir string, task *tes.Task) (*FileMapper, error) {

	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	mapper := &FileMapper{dir: dir}

	// Add all the volumes to the mapper
	for _, vol := range task.Volumes {
		err := mapper.AddTmpVolume(vol)
		if err != nil {
			return nil, err
		}
	}

	// Add all the inputs to the mapper
	for _, input := range task.Inputs {
		err := mapper.AddInput(input)
		if err != nil {
			return nil, err
		}
	}

	// Add all the outputs to the mapper
	for _, output := range task.Outputs {
		err := mapper.AddOutput(output)
		if err != nil {
			return nil, err
		}
	}

	// Check the executor paths.
	for _, exec := range task.Executors {
		err := mapper.AddExecutor(exec)
		if err != nil {
			return nil, err
		}
	}

	return mapper, nil
}

func (mapper *FileMapper) AddExecutor(exec *tes.Executor) error {
	exec = proto.Clone(exec).(*tes.Executor)

	if exec.Stdin != "" {
		// Ensure the path is valid
		exec.Stdin = mapper.HostPath(exec.Stdin)
		err := mapper.CheckPath(exec.Stdin)
		if err != nil {
			return err
		}

		// Ensure the directory exists.
		err = util.EnsurePath(exec.Stdin)
		if err != nil {
			return err
		}
	}

	if exec.Stdout != "" {
		// Ensure the path is valid
		exec.Stdout = mapper.HostPath(exec.Stdout)
		err := mapper.CheckPath(exec.Stdout)
		if err != nil {
			return err
		}

		// Ensure the directory exists.
		err = util.EnsurePath(exec.Stdout)
		if err != nil {
			return err
		}
	}

	if exec.Stderr != "" {
		// Ensure the path is valid
		exec.Stderr = mapper.HostPath(exec.Stderr)
		err := mapper.CheckPath(exec.Stderr)
		if err != nil {
			return err
		}

		// Ensure the directory exists.
		err = util.EnsurePath(exec.Stderr)
		if err != nil {
			return err
		}
	}
	mapper.Executors = append(mapper.Executors, exec)
	return nil
}

// AddVolume adds a mapped volume to the mapper. A corresponding Volume record
// is added to mapper.Volumes.
//
// If the volume paths are invalid or can't be mapped, an error is returned.
func (mapper *FileMapper) AddVolume(hostPath string, mountPoint string, readonly bool) error {
	vol := Volume{
		HostPath:      hostPath,
		ContainerPath: mountPoint,
		Readonly:      readonly,
	}

	for i, v := range mapper.Volumes {
		// check if this volume is already present in the mapper
		if vol == v {
			return nil
		}

		// If the proposed RW Volume is a subpath of an existing RW Volume
		// do not add it to the mapper
		// If an existing RW Volume is a subpath of the proposed RW Volume, replace it with
		// the proposed RW Volume
		if !vol.Readonly && !v.Readonly {
			if IsSubpath(vol.ContainerPath, v.ContainerPath) {
				return nil
			} else if IsSubpath(v.ContainerPath, vol.ContainerPath) {
				mapper.Volumes[i] = vol
				return nil
			}
		}
	}

	mapper.Volumes = append(mapper.Volumes, vol)
	return nil
}

// HostPath returns a mapped path.
//
// The path is concatenated to the mapper's base dir.
// e.g. If the mapper is configured with a base dir of "/tmp/mapped_files", then
// mapper.HostPath("/home/ubuntu/myfile") will return "/tmp/mapped_files/home/ubuntu/myfile".
//
// The mapped path is required to be a subpath of the mapper's base directory.
// e.g. mapper.HostPath("../../foo") should fail with an error.
func (mapper *FileMapper) HostPath(src string) string {
	if src == "" {
		return src
	}
	p := path.Join(mapper.dir, src)
	p = path.Clean(p)
	return p
}

func (mapper *FileMapper) CheckPath(p string) error {
	if !IsSubpath(p, mapper.dir) {
		return fmt.Errorf("Invalid path: %s is not a valid subpath of %s", p, mapper.dir)
	}
	return nil
}

func CreateWorkDir(dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	return util.EnsureDir(dir)
}

// AddTmpVolume creates a directory on the host based on the declared path in
// the container and adds it to mapper.Volumes.
//
// If the path can't be mapped, an error is returned.
func (mapper *FileMapper) AddTmpVolume(mountPoint string) error {
	hostPath := mapper.HostPath(mountPoint)

	err := util.EnsureDir(hostPath)
	if err != nil {
		return err
	}

	err = mapper.AddVolume(hostPath, mountPoint, false)
	if err != nil {
		return err
	}
	return nil
}

// AddInput adds an input to the mapped files for the given TaskParameter.
// A copy of the TaskParameter will be added to mapper.Inputs, with the
// "Path" field updated to the mapped host path.
//
// If the path can't be mapped an error is returned.
func (mapper *FileMapper) AddInput(input *tes.TaskParameter) error {
	hostPath := mapper.HostPath(input.Path)

	err := mapper.CheckPath(hostPath)
	if err != nil {
		return err
	}

	err = util.EnsurePath(hostPath)
	if err != nil {
		return err
	}

	// Add input volumes
	err = mapper.AddVolume(hostPath, input.Path, true)
	if err != nil {
		return err
	}

	// Create a TaskParameter for the input with a path mapped to the host
	hostIn := proto.Clone(input).(*tes.TaskParameter)
	hostIn.Path = hostPath
	mapper.Inputs = append(mapper.Inputs, hostIn)
	return nil
}

// AddOutput adds an output to the mapped files for the given TaskParameter.
// A copy of the TaskParameter will be added to mapper.Outputs, with the
// "Path" field updated to the mapped host path.
//
// If the path can't be mapped, an error is returned.
func (mapper *FileMapper) AddOutput(output *tes.TaskParameter) error {
	hostPath := mapper.HostPath(output.Path)

	err := mapper.CheckPath(hostPath)
	if err != nil {
		return err
	}

	hostDir := hostPath
	mountDir := output.Path
	if output.Type == tes.FileType_FILE {
		hostDir = path.Dir(hostPath)
		mountDir = path.Dir(output.Path)
	}

	err = util.EnsureDir(hostDir)
	if err != nil {
		return err
	}

	// Add output volumes
	err = mapper.AddVolume(hostDir, mountDir, false)
	if err != nil {
		return err
	}

	// Create a TaskParameter for the out with a path mapped to the host
	hostOut := proto.Clone(output).(*tes.TaskParameter)
	hostOut.Path = hostPath
	mapper.Outputs = append(mapper.Outputs, hostOut)
	return nil
}

// IsSubpath returns true if the given path "p" is a subpath of "base".
func IsSubpath(p string, base string) bool {
	return strings.HasPrefix(p, base)
}
