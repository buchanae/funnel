package datastore

import (
  "context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/events"
)

/*
Datastore is missing support for some types we need:
- map[string]string
- uint32

So we need to do some extra work to map a task to/from a []datastore.Property.
It also allows us to be very selective about what gets saved and indexed.
*/

type chunktype int

const (
	base chunktype = iota
	state
	logevent
	stdout
	stderr
	content
  syslog
)

type event struct {
  *events.Event
  Attempt, Index int
}

type chunk struct {
	Type              chunktype  `datastore:",noindex"`
	Id, CreationTime  string     `datastore:",omitempty"`
	State             string     `datastore:",omitempty"`
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

  Event         *event `datastore:",noindex,omitempty"`
	Stdout, Stderr     string               `datastore:",noindex,omitempty"`
}

type executor struct {
	Image, Workdir, Stdin, Stdout, Stderr string   `datastore:",noindex,omitempty"`
	Command                               []string `datastore:",noindex,omitempty"`
	Env                                   []kv     `datastore:",noindex,omitempty"`
}

type param struct {
	Name, Description, Url, Path, Type, Content string `datastore:",noindex,omitempty"`
}

func fromTask(t *tes.Task) *chunk {
	z := &chunk{
		Type:         base,
		Id:           t.Id,
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

func toTask(c *chunk, z *tes.Task) {
  switch c.Type {
  case state:
    z.State = tes.State(tes.State_value[c.State])
  case logevent:
    tb := events.TaskBuilder{z}
    ev := c.Event.Event
    ev.Attempt = uint32(c.Event.Attempt)
    ev.Index = uint32(c.Event.Index)
    tb.WriteEvent(context.Background(), ev)
  case base:
    z.Id =           c.Id
    z.CreationTime = c.CreationTime
    z.Name =         c.Name
    z.Description =  c.Description
    z.Volumes =      c.Volumes
    z.Tags =         toMap(c.Tags)
    z.Resources =    &tes.Resources{
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
        Type:        tes.FileType(tes.FileType_value[i.Type]),
      })
    }
    for _, i := range c.Outputs {
      z.Outputs = append(z.Outputs, &tes.Output{
        Name:        i.Name,
        Description: i.Description,
        Url:         i.Url,
        Path:        i.Path,
        Type:        tes.FileType(tes.FileType_value[i.Type]),
      })
    }
  }
}

func fromState(s tes.State) *chunk {
	return &chunk{Type: state, State: s.String()}
}

func fromLogEvent(e *events.Event) *chunk {
	return &chunk{
		Type:       logevent,
    Event: &event{
      Event: e,
      Attempt: int(e.Attempt),
      Index: int(e.Index),
    },
	}
}

func fromStdout(s string) *chunk {
	return &chunk{Type: stdout, Stdout: s}
}
func fromStderr(s string) *chunk {
	return &chunk{Type: stderr, Stderr: s}
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
