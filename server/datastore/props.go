package datastore

import (
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

/*
Datastore is missing support for some types we need:
- map[string]string
- uint32

So we need to do some extra work to map a task to/from a []datastore.Property.
It also allows us to be very selective about what gets saved and indexed.
*/

type event struct {
	Type                       int32
	Attempt, Index             int    `datastore:",noindex,omitempty"`
	Stdout, Stderr, Msg, Level string `datastore:",noindex,omitempty"`
	Fields                     []kv   `datastore:",noindex,omitempty"`
}

type task struct {
	Id, CreationTime  string `datastore:",omitempty"`
	State             int32
	Name, Description string     `datastore:",noindex,omitempty"`
	Executors         []executor `datastore:",noindex,omitempty"`
	Inputs            []param    `datastore:",noindex,omitempty"`
	Outputs           []param    `datastore:",noindex,omitempty"`
	Volumes           []string   `datastore:",noindex,omitempty"`
	Tags              []kv

	CpuCores      int64    `datastore:",noindex,omitempty"`
	RamGb, DiskGb float64  `datastore:",noindex,omitempty"`
	Preemptible   bool     `datastore:",noindex,omitempty"`
	Zones         []string `datastore:",noindex,omitempty"`

	TaskLogs []tasklog `datastore:",noindex,omitempty"`
}

type tasklog struct {
	*tes.TaskLog
	Metadata []kv `datastore:",noindex,omitempty"`
}

type executor struct {
	Image, Workdir, Stdin, Stdout, Stderr string   `datastore:",noindex,omitempty"`
	Command                               []string `datastore:",noindex,omitempty"`
	Env                                   []kv     `datastore:",noindex,omitempty"`
}

type param struct {
	Name, Description, Url, Path, Content string `datastore:",noindex,omitempty"`
	Type                                  int32  `datastore:",noindex,omitempty"`
}

func fromTask(t *tes.Task) *task {
	z := &task{
		Id:           t.Id,
		State:        int32(t.State),
		CreationTime: t.CreationTime,
		Name:         t.Name,
		Description:  t.Description,
		Volumes:      t.Volumes,
		Tags:         fromMap(t.Tags),
	}
	if t.Resources != nil {
		z.CpuCores = int64(t.Resources.CpuCores)
		z.RamGb = t.Resources.RamGb
		z.DiskGb = t.Resources.DiskGb
		z.Preemptible = t.Resources.Preemptible
		z.Zones = t.Resources.Zones
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
			Type:        int32(i.Type),
		})
	}
	for _, i := range t.Outputs {
		z.Outputs = append(z.Outputs, param{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        int32(i.Type),
		})
	}
	for _, i := range t.Logs {
		z.TaskLogs = append(z.TaskLogs, tasklog{
			TaskLog:  i,
			Metadata: fromMap(i.Metadata),
		})
	}
	return z
}

func toTask(c *task, z *tes.Task) {
	z.Id = c.Id
	z.CreationTime = c.CreationTime
	z.State = tes.State(c.State)
	z.Name = c.Name
	z.Description = c.Description
	z.Volumes = c.Volumes
	z.Tags = toMap(c.Tags)
	z.Resources = &tes.Resources{
		CpuCores:    uint32(c.CpuCores),
		RamGb:       c.RamGb,
		DiskGb:      c.DiskGb,
		Preemptible: c.Preemptible,
		Zones:       c.Zones,
	}
	for _, e := range c.Executors {
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
	for _, i := range c.Inputs {
		z.Inputs = append(z.Inputs, &tes.Input{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        tes.FileType(i.Type),
		})
	}
	for _, i := range c.Outputs {
		z.Outputs = append(z.Outputs, &tes.Output{
			Name:        i.Name,
			Description: i.Description,
			Url:         i.Url,
			Path:        i.Path,
			Type:        tes.FileType(i.Type),
		})
	}
	for _, i := range c.TaskLogs {
		tl := i.TaskLog
		tl.Metadata = toMap(i.Metadata)
		z.Logs = append(z.Logs, tl)
	}
}

func fromEvent(e *events.Event) *event {
	z := &event{
		Type:    int32(e.Type),
		Attempt: int(e.Attempt),
		Index:   int(e.Index),
	}
	switch e.Type {
	case events.Type_SYSTEM_LOG:
		l := e.GetSystemLog()
		z.Msg = l.Msg
		z.Level = l.Level
		z.Fields = fromMap(l.Fields)
	case events.Type_EXECUTOR_STDOUT:
		z.Stdout = e.GetStdout()
	case events.Type_EXECUTOR_STDERR:
		z.Stdout = e.GetStderr()
	}
	return z
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
