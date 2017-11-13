package datastore

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

/*
Datastore is missing support for some types we need:
- map[string]string
- uint32

So we need to do some extra work to map a task to/from a []datastore.Property.
It also allows us to be very selective about what gets saved and indexed.
*/

type task struct {
	Id, CreationTime  string
	Name, Description string     `datastore:",noindex"`
	Resources         *resources `datastore:",noindex,omitempty"`
	Executors         []executor `datastore:",noindex"`
	Inputs            []param    `datastore:",noindex,omitempty"`
	Outputs           []param    `datastore:",noindex,omitempty"`
	Volumes           []string   `datastore:",noindex,omitempty"`
	Tags              []kv
}

type executor struct {
	Image, Workdir, Stdin, Stdout, Stderr string   `datastore:",noindex"`
	Command                               []string `datastore:",noindex"`
	Env                                   []kv     `datastore:",noindex"`
}

type param struct {
	Name, Description, Url, Path, Type string `datastore:",noindex"`
}

type resources struct {
	CpuCores      int64    `datastore:",omitempty"`
	RamGb, DiskGb float64  `datastore:",omitempty"`
	Preemptible   bool     `datastore:",omitempty"`
	Zones         []string `datastore:",omitempty"`
}

func fromTask(t *tes.Task) *task {
	z := &task{
		Id:           t.Id,
		CreationTime: t.CreationTime,
		Name:         t.Name,
		Description:  t.Description,
		Volumes:      t.Volumes,
		Tags:         fromMap(t.Tags),
	}
	if t.Resources != nil {
		z.Resources = &resources{
			CpuCores:    int64(t.Resources.CpuCores),
			RamGb:       t.Resources.RamGb,
			DiskGb:      t.Resources.DiskGb,
			Preemptible: t.Resources.Preemptible,
			Zones:       t.Resources.Zones,
		}
	}
	for _, e := range t.Executors {
		z.Executors = append(z.Executors, executor{
			Image:   e.Image,
			Workdir: e.Workdir,
			Stdin:   e.Stdin,
			Stdout:  e.Stdout,
			Stderr:  e.Stderr,
			Command: e.Command,
			Env:     fromMap(e.Env),
		})
	}
	for _, i := range t.Inputs {
		z.Inputs = append(z.Inputs, param{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        i.Type.String(),
		})
	}
	for _, i := range t.Outputs {
		z.Outputs = append(z.Outputs, param{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        i.Type.String(),
		})
	}
	return z
}

func toTask(t *task) *tes.Task {
	z := &tes.Task{
		Id:           t.Id,
		CreationTime: t.CreationTime,
		Name:         t.Name,
		Description:  t.Description,
		Volumes:      t.Volumes,
		Tags:         toMap(t.Tags),
	}
	if t.Resources != nil {
		z.Resources = &tes.Resources{
			CpuCores:    uint32(t.Resources.CpuCores),
			RamGb:       t.Resources.RamGb,
			DiskGb:      t.Resources.DiskGb,
			Preemptible: t.Resources.Preemptible,
			Zones:       t.Resources.Zones,
		}
	}
	for _, e := range t.Executors {
		z.Executors = append(z.Executors, &tes.Executor{
			Image:   e.Image,
			Workdir: e.Workdir,
			Stdin:   e.Stdin,
			Stdout:  e.Stdout,
			Stderr:  e.Stderr,
			Command: e.Command,
			Env:     toMap(e.Env),
		})
	}
	for _, i := range t.Inputs {
		z.Inputs = append(z.Inputs, &tes.Input{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        tes.FileType(tes.FileType_value[i.Type]),
		})
	}
	for _, i := range t.Outputs {
		z.Outputs = append(z.Outputs, &tes.Output{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        tes.FileType(tes.FileType_value[i.Type]),
		})
	}
	return z
}

type state struct {
	State string
}

func fromState(s tes.State) *state {
	return &state{s.String()}
}

type tasklog struct {
	StartTime, EndTime string               `datastore:",noindex"`
	Outputs            []*tes.OutputFileLog `datastore:",noindex"`
	SystemLogs         []string             `datastore:",noindex,omitempty"`
	Metadata           []kv                 `datastore:",noindex,omitempty"`
}

func fromTaskLog(tl *tes.TaskLog) *tasklog {
	return &tasklog{
		StartTime:  tl.StartTime,
		EndTime:    tl.EndTime,
		Outputs:    tl.Outputs,
		SystemLogs: tl.SystemLogs,
		Metadata:   fromMap(tl.Metadata),
	}
}

type execlog struct {
	StartTime, EndTime string `datastore:",noindex"`
	ExitCode           int64  `datastore:",noindex"`
}

func fromExecLog(el *tes.ExecutorLog) *execlog {
	return &execlog{
		StartTime: el.StartTime,
		EndTime:   el.EndTime,
	}
}

type stdout struct {
	Stdout string
}

func fromStdout(s string) *stdout {
	return &stdout{s}
}

type stderr struct {
	Stderr string
}

func fromStderr(s string) *stderr {
	return &stderr{s}
}

type kv struct {
	Key, Value string
}

func fromMap(m map[string]string) []kv {
	var out []kv
	for k, v := range m {
		out = append(out, kv{k, v})
	}
	return out
}

func toMap(kvs []kv) map[string]string {
	out := map[string]string{}
	for _, kv := range kvs {
		out[kv.Key] = kv.Value
	}
	return out
}
