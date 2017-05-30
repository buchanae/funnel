package perf

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/tests/testutils"
	"testing"
)

var log = logger.New("perf")

func BenchmarkRunTinyTask(b *testing.B) {
	var fun = testutils.NewFunnel()
	fun.Conf.LogLevel = "info"
	fun.StartServer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fun.RunTask(&tes.Task{
			Executors: []*tes.Executor{
				{
					ImageName: "alpine",
					Cmd:       []string{"sleep 1"},
				},
			},
		})
	}
}

func BenchmarkRunMidsizeTask(b *testing.B) {
	var fun = testutils.NewFunnel()
	fun.Conf.LogLevel = "info"
	fun.StartServer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fun.RunTask(&tes.Task{
			Name:        "test-midsize-perf",
			Description: "some reasonable long description of the task to be run.",
			Project:     "project-id",
			Resources: &tes.Resources{
				CpuCores: 10,
				RamGb:    10.0,
				SizeGb:   100.0,
			},
			Volumes: []string{"/tmp"},
			Tags:    map[string]string{"one": "two"},
			Inputs: []*tes.TaskParameter{
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
			},
			Outputs: []*tes.TaskParameter{
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
				{
					Url:  "/path/to/nowhere",
					Path: "/tmp/path/to/nowhere",
				},
			},
			Executors: []*tes.Executor{
				{
					ImageName: "alpine",
					Cmd:       []string{"echo", "foo", "bar", "baz"},
					Stdout:    "/tmp/stdout",
					Stderr:    "/tmp/stderr",
				},
			},
		})
	}
}
