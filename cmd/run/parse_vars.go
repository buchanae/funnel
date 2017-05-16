package run

import (
	"errors"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrKeyFmt = errors.New("Arguments passed to --in, --out and --env must be of the form: KEY=VALUE")
var ErrStorageScheme = errors.New("File paths must be prefixed with one of:\n file://\n gs://\n s3://")

func DuplicateKeyErr(key string) error {
	return errors.New("Can't use the same KEY for multiple --in, --out, --env arguments: " + key)
}

// Parse CLI variable definitions (e.g "varname=value") into usable task values.
func valsToTask(vals flagVals) (*tes.Task, error) {
	var err error

	// Build the task message
	task := &tes.Task{
		Name:        vals.name,
		Project:     vals.project,
		Description: vals.description,
		Resources: &tes.Resources{
			CpuCores:    uint32(vals.cpu),
			RamGb:       vals.ram,
			SizeGb:      vals.disk,
			Zones:       vals.zones,
			Preemptible: vals.preemptible,
		},
		Executors: []*tes.Executor{
			{
				ImageName: vals.container,
				Cmd:       vals.cmd,
				Workdir:   vals.workdir,
				Stdout:    "/opt/funnel/outputs/stdout",
				Stderr:    "/opt/funnel/outputs/stderr",
				// TODO no ports
				Ports: nil,
			},
		},
		Volumes: vals.volumes,
	}

	// Only set the stdin path if the --stdin flag was used.
	if vals.stdin != "" {
		task.Executors[0].Stdin = "/opt/funnel/inputs/stdin"
	}

	// Any error occuring during parsing the variables an preparing the task
	// is a fatal error, so I'm using panic/recover to simplify error handling.
	defer func() {
		err = recover().(error)
	}()

	// Helper to make sure variable keys are unique.
	setenv := func(key, val string) {
		_, exists := task.Executors[0].Environ[key]
		if exists {
			panic(DuplicateKeyErr(key))
		}
		task.Executors[0].Environ[key] = val
	}

	for _, raw := range vals.inputs {
		k, v := parseCliVar(raw)
		setenv(k, v)
		url := resolvePath(v)
		task.Inputs = append(task.Inputs, &tes.TaskParameter{
			Name: k,
			Url:  url,
			Path: "/opt/funnel/inputs/" + stripStoragePrefix(url),
		})
	}

	for _, raw := range vals.inputDirs {
		k, v := parseCliVar(raw)
		setenv(k, v)
		url := resolvePath(v)
		task.Inputs = append(task.Inputs, &tes.TaskParameter{
			Name: k,
			Url:  url,
			Path: "/opt/funnel/inputs/" + stripStoragePrefix(url),
			Type: tes.FileType_DIRECTORY,
		})
	}

	for _, raw := range vals.contents {
		k, v := parseCliVar(raw)
		setenv(k, v)
		task.Inputs = append(task.Inputs, &tes.TaskParameter{
			Name:     k,
			Path:     "/opt/funnel/inputs/" + stripStoragePrefix(v),
			Contents: getContents(v),
		})
	}

	for _, raw := range vals.outputs {
		k, v := parseCliVar(raw)
		setenv(k, v)
		url := resolvePath(v)
		task.Outputs = append(task.Outputs, &tes.TaskParameter{
			Name: k,
			Url:  url,
			Path: "/opt/funnel/outputs/" + stripStoragePrefix(url),
		})
	}

	for _, raw := range vals.outputDirs {
		k, v := parseCliVar(raw)
		setenv(k, v)
		url := resolvePath(v)
		task.Outputs = append(task.Outputs, &tes.TaskParameter{
			Name: k,
			Url:  url,
			Path: "/opt/funnel/outputs/" + stripStoragePrefix(url),
			Type: tes.FileType_DIRECTORY,
		})
	}

	for _, raw := range vals.environ {
		k, v := parseCliVar(raw)
		setenv(k, v)
	}

	for _, raw := range vals.tags {
		k, v := parseCliVar(raw)
		task.Tags[k] = v
	}

	return task, err
}

func getContents(p string) string {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func parseCliVar(raw string) (string, string) {
	re := regexp.MustCompile("=")
	res := re.Split(raw, -1)

	if len(res) != 2 {
		panic(ErrKeyFmt)
	}

	key := res[0]
	val := res[1]
	return key, val
}

// Give a input/output URL "raw", return the path of the file
// relative to the container.
func containerPath(raw, base string) string {
	url := resolvePath(raw)
	p := stripStoragePrefix(url)
	return base + p
}

func stripStoragePrefix(url string) string {
	re := regexp.MustCompile("[a-z0-9]+://")
	if !re.MatchString(url) {
		panic(ErrStorageScheme)
	}
	path := re.ReplaceAllString(url, "")
	return strings.TrimPrefix(path, "/")
}

func resolvePath(url string) string {
	local := strings.HasPrefix(url, "/") || strings.HasPrefix(url, ".") ||
		strings.HasPrefix(url, "~")
	re := regexp.MustCompile("[a-z0-9]+://")
	prefixed := re.MatchString(url)

	switch {
	case local:
		path, err := filepath.Abs(url)
		if err != nil {
			panic(err)
		}
		return "file://" + path
	case prefixed:
		return url
	default:
		panic(fmt.Errorf("could not resolve filepath: %s", url))
	}
}
