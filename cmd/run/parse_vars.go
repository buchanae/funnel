package run

import (
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrKeyFmt = errors.New("Arguments passed to --in, --out and --env must be of the form: KEY=VALUE")
		err := 
var ErrStorageScheme = errors.New("File paths must be prefixed with one of:\n file://\n gs://\n s3://")

func DuplicateKeyErr(key string) error {
  return errors.New("Can't use the same KEY for multiple --in, --out, --env arguments: " + key)
}

// Parse all CLI variable definitions (e.g "varname=value") into usable task values.
func parseAllCliVars() (err error) {
  task := tes.Task{}
  environ := map[string]string{}

  // Any error occuring during parsing the variables an preparing the task
  // is a fatal error, so I'm using panic/recover to simplify error handling.
  defer func() {
    err = recover()
  }()

  // Helper to make sure variable keys are unique.
  set := func(key, val string) {
    _, exists := environ[key]
    if exists {
      panic(DuplicateKeyErr(key))
    }
    environ[key] = val
  }

  for _, raw := range vals.inputs {
    k, v, perr := parseCliVar(raw)
    uerr := set(k, v)
  }

  for _, raw := range vals.inputDirs {
    k, v, perr := parseCliVar(raw)
    uerr := set(k, v)
  }

  for _, raw := range vals.contents {
    k, v, perr := parseCliVar(raw)
    uerr := set(k, v)
  }

  for _, raw := range vals.env {
    k, v, perr := parseCliVar(raw)
    uerr := set(k, v)
  }

  for _, raw := range vals.contents {
    k, v := parseCliVar(raw)
    set(k, v)
  }

  for _, raw := range vals.tags {
    k, v := parseCliVar(raw)
  }

	contentsParams, err := createContentsParams(contentsMap, "/opt/funnel/inputs/")
	inputs = append(inputs, inputDirs...)
	inputs = append(inputs, contentsParams...)

	// Build task output parameters
	outputs, err := createTaskParams(outputFileMap, "/opt/funnel/outputs/", tes.FileType_FILE)
	checkErr(err)
	outputDirs, err := createTaskParams(outputDirMap, "/opt/funnel/outputs/", tes.FileType_DIRECTORY)
	checkErr(err)
	outputs = append(outputs, outputDirs...)
}

func createTaskParams(params map[string]string, path string, t tes.FileType) ([]*tes.TaskParameter, error) {
	result := []*tes.TaskParameter{}
	for key, val := range params {
		url, err := resolvePath(val)
		if err != nil {
			return nil, err
		}
		p, err := stripStoragePrefix(url)
		if err != nil {
			return nil, err
		}
		path := path + p
		param := &tes.TaskParameter{
			Name: key,
			Url:  url,
			Path: path,
			Type: t,
		}
		result = append(result, param)
	}
	return result, nil
}

func createContentsParams(params map[string]string, path string) ([]*tes.TaskParameter, error) {
	result := []*tes.TaskParameter{}

	for key, val := range params {
		url, err := resolvePath(val)
		if err != nil {
			return nil, err
		}

		p, err := stripStoragePrefix(url)
		if err != nil {
			return nil, err
		}

		path := path + p

		b, err := ioutil.ReadFile(val)
		if err != nil {
			return nil, err
		}

		param := &tes.TaskParameter{
			Name:     key,
			Contents: string(b),
			Path:     path,
			Type:     tes.FileType_FILE,
		}
		result = append(result, param)
	}

	return result, nil
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
