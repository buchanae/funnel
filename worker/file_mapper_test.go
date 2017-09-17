package worker

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	logger.Configure(logger.DebugConfig())
}

func TestMapTask(t *testing.T) {
	tmp, err := ioutil.TempDir("", "funnel-test-mapper")
	if err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	task := &tes.Task{
		Inputs: []*tes.TaskParameter{
			{
				Name: "f1",
				Url:  "file://" + cwd + "/testdata/f1.txt",
				Path: "/opt/funnel/inputs/testdata/f1.txt",
			},
			{
				Name: "f4",
				Url:  "file://" + cwd + "/testdata/f4",
				Path: "/opt/funnel/inputs/testdata/f4",
				Type: tes.FileType_DIRECTORY,
			},
			{
				Name:     "c1",
				Path:     "/opt/funnel/inputs/testdata/contents.txt",
				Contents: "test content\n",
			},
		},
		Outputs: []*tes.TaskParameter{
			{
				Name: "stdout-0",
				Url:  "file://" + cwd + "/testdata/stdout-first",
				Path: "/opt/funnel/outputs/stdout-0",
			},
			{
				Name: "o9",
				Url:  "file://" + cwd + "/testdata/o9",
				Path: "/opt/funnel/outputs/sub/o9",
				Type: tes.FileType_DIRECTORY,
			},
		},
		Executors: []*tes.Executor{
			{
				Stdin:  "/opt/funnel/execs/stdindir/stdin.txt",
				Stdout: "/opt/funnel/execs/stdoutdir/stdout.txt",
				Stderr: "/opt/funnel/execs/stderrdir/stderr.txt",
			},
		},
		Volumes: []string{"/volone", "/voltwo"},
	}

	f, err := NewFileMapper(tmp, task)
	if err != nil {
		t.Fatal(err)
	}

	ei := []*tes.TaskParameter{
		{
			Name: "f1",
			Url:  "file://" + cwd + "/testdata/f1.txt",
			Path: tmp + "/opt/funnel/inputs/testdata/f1.txt",
		},
		{
			Name: "f4",
			Url:  "file://" + cwd + "/testdata/f4",
			Path: tmp + "/opt/funnel/inputs/testdata/f4",
			Type: tes.FileType_DIRECTORY,
		},
	}

	eo := []*tes.TaskParameter{
		{
			Name: "stdout-0",
			Url:  "file://" + cwd + "/testdata/stdout-first",
			Path: tmp + "/opt/funnel/outputs/stdout-0",
		},
		{
			Name: "o9",
			Url:  "file://" + cwd + "/testdata/o9",
			Path: tmp + "/opt/funnel/outputs/sub/o9",
			Type: tes.FileType_DIRECTORY,
		},
	}

	ev := []Volume{
		{
			HostPath:      tmp + "/volone",
			ContainerPath: "/volone",
			Readonly:      false,
		},
		{
			HostPath:      tmp + "/voltwo",
			ContainerPath: "/voltwo",
			Readonly:      false,
		},
		{
			HostPath:      tmp + "/opt/funnel/inputs/testdata/f1.txt",
			ContainerPath: "/opt/funnel/inputs/testdata/f1.txt",
			Readonly:      true,
		},
		{
			HostPath:      tmp + "/opt/funnel/inputs/testdata/f4",
			ContainerPath: "/opt/funnel/inputs/testdata/f4",
			Readonly:      true,
		},
		{
			HostPath:      tmp + "/opt/funnel/inputs/testdata/contents.txt",
			ContainerPath: "/opt/funnel/inputs/testdata/contents.txt",
			Readonly:      true,
		},
		{
			HostPath:      tmp + "/opt/funnel/outputs",
			ContainerPath: "/opt/funnel/outputs",
			Readonly:      false,
		},
	}

	if diff := deep.Equal(f.Inputs, ei); diff != nil {
		t.Log("Expected", fmt.Sprintf("%+v", ei))
		t.Log("Actual", fmt.Sprintf("%+v", f.Inputs))
		for _, d := range diff {
			t.Log("Diff", d)
		}
		t.Fatal("unexpected mapper inputs")
	}

	c, err := ioutil.ReadFile(tmp + "/opt/funnel/inputs/testdata/contents.txt")
	if err != nil {
		t.Fatal(err)
	}
	if string(c) != "test content\n" {
		t.Fatal("unexpected content")
	}

	if diff := deep.Equal(f.Outputs, eo); diff != nil {
		t.Log("Expected", fmt.Sprintf("%+v", eo))
		t.Log("Actual", fmt.Sprintf("%+v", f.Outputs))
		for _, d := range diff {
			t.Log("Diff", d)
		}
		t.Fatal("unexpected mapper outputs")
	}

	if diff := deep.Equal(f.Volumes, ev); diff != nil {
		t.Log("Expected", fmt.Sprintf("%+v", ev))
		t.Log("Actual", fmt.Sprintf("%+v", f.Volumes))
		for _, d := range diff {
			t.Log("Diff", d)
		}
		t.Fatal("unexpected mapper volumes")
	}

	// executor stdin/out/err directory paths should be created
	if !isDir(tmp + "/opt/funnel/execs/stdindir") {
		t.Error("expected stdin directory path to be created")
	}
	if !isDir(tmp + "/opt/funnel/execs/stdoutdir") {
		t.Error("expected stdout directory path to be created")
	}
	if !isDir(tmp + "/opt/funnel/execs/stderrdir") {
		t.Error("expected stderr directory path to be created")
	}
	// ... but the files should not be created
	if exists(tmp + "/opt/funnel/execs/stdindir/stdin.txt") {
		t.Error("did not expect stdin file to be created")
	}
	if exists(tmp + "/opt/funnel/execs/stdoutdir/stdout.txt") {
		t.Error("did not expect stdout file to be created")
	}
	if exists(tmp + "/opt/funnel/execs/stderrdir/stderr.txt") {
		t.Error("did not expect stderr file to be created")
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return info.Mode().IsRegular()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}
