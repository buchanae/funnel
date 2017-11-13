package datastore

import (
  "fmt"
  "cloud.google.com/go/datastore"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

/*
Datastore is missing support for some types we need:
- map[string]string
- uint32

So we need to do some extra work to map a task to/from a []datastore.Property.
It also allows us to be very selective about what gets saved and indexed.
*/

// whether or not a field will be indexed.
var index = true
var noindex = false

// shortcut for datastore.Property{}
func p(name string, val interface{}, index bool) datastore.Property {
  return datastore.Property{name, val, !index}
}

type resources struct {
  CpuCores int64
}

type task struct {
  Id string `datastore:"noindex"`
  Name string `datastore:"noindex"`
  Description string `datastore:"noindex"`
  CreationTime string `datastore:"noindex"`
  Executors []executor
}

func saveTask(t *tes.Task) datastore.PropertyList {

  if t.Resources !=  nil {
    ps = append(ps, saveResources(t.Resources))
  }

  if t.Volumes != nil {
    ps = append(ps, p("Volumes", t.Volumes, noindex))
  }
  /*
  if t.Inputs != nil {
    ps = append(ps, saveInputs(t.Inputs))
  }
  if t.Outputs != nil {
    ps = append(ps, saveOutputs(t.Outputs))
  }
  */
  if t.Executors != nil {
    ps = append(ps, saveExecutors(t.Executors))
  }
  //ps = append(ps, saveMap("Tags.", t.Tags, index)...)
  return ps
}

func saveState(s tes.State) datastore.PropertyList {
  return []datastore.Property{
    p("State", s.String(), index),
  }
}

func saveResources(r *tes.Resources) datastore.Property {
  ent := &datastore.Entity{
    Key: datastore.NameKey("Resources", "Resources", nil),
  }
  ent.Properties = append(ent.Properties,
    p("CpuCores", int64(r.CpuCores), noindex),
    p("Preemptible", r.Preemptible, noindex),
    p("RamGb", r.RamGb, noindex),
    p("DiskGb", r.DiskGb, noindex),
  )
  if r.GetZones() != nil {
    ent.Properties = append(ent.Properties,
      p("Zones", r.Zones, noindex))
  }
  return p("Resources", ent, noindex)
}

func saveExecutors(exs []*tes.Executor) datastore.Property {
  ent := &datastore.Entity{
    Key: datastore.NameKey("Resources", "Resources", nil),
  }
  for i, v := range exs {
    ps = append(ps,
      p(k + "Image", v.Image, noindex),
      p(k + "Workdir", v.Workdir, noindex),
      p(k + "Stdin", v.Stdin, noindex),
      p(k + "Stdout", v.Stdout, noindex),
      p(k + "Stderr", v.Stderr, noindex),
      p(k + "Command", v.Command, noindex),
    )
    ps = append(ps, saveMap(k + "Env.", v.Env, noindex)...)
  }
  return p("Executors", ent, noindex)
}

/*
func saveInputs(in []*tes.Input) *datastore.Entity {
  ent := &datastore.Entity{
    Key: datastore.NameKey("Inputs", "Inputs", nil),
  }
  for i, v := range in {
    k := fmt.Sprintf("Inputs.%d.", i)
    ps = append(ps,
      p(k + "Name", v.Name, noindex),
      p(k + "Description", v.Description, noindex),
      p(k + "Url", v.Url, noindex),
      p(k + "Path", v.Path, noindex),
      p(k + "Type", v.Type.String(), noindex),
    )
  }
  return ent
}

func saveOutputs(t *tes.Task) datastore.PropertyList {
  var ps []datastore.Property
  for i, v := range t.Outputs {
    k := fmt.Sprintf("Outputs.%d.", i)
    ps = append(ps,
      p(k + "Name", v.Name, noindex),
      p(k + "Description", v.Description, noindex),
      p(k + "Url", v.Url, noindex),
      p(k + "Path", v.Path, noindex),
      p(k + "Type", v.Type.String(), noindex),
    )
  }
  return ps
}
*/

func saveTaskLog(t *tes.TaskLog) datastore.PropertyList {
  ps := []datastore.Property{
    p("StartTime", t.StartTime, noindex),
    p("EndTime", t.EndTime, noindex),
    p("Outputs", t.Outputs, noindex),
    p("SystemLogs", t.SystemLogs, noindex),
  }
  ps = append(ps, saveMap("Metadata.", t.Metadata, noindex)...)
  return ps
}

func saveExecutorLog(e *tes.ExecutorLog) datastore.PropertyList {
  return []datastore.Property{
    p("StartTime", e.StartTime, noindex),
    p("EndTime", e.EndTime, noindex),
    p("ExitCode", e.ExitCode, noindex),
  }
}

func saveStdout(stdout string) datastore.PropertyList {
  return []datastore.Property{
    p("Stdout", stdout, noindex),
  }
}

func saveStderr(stderr string) datastore.PropertyList {
  return []datastore.Property{
    p("Stderr", stderr, noindex),
  }
}

func saveMap(prefix string, m map[string]string, index bool) datastore.PropertyList {
  var ps []datastore.Property
  for k, v := range m {
    ps = append(ps, p(prefix + k, v, index))
  }
  return ps
}
