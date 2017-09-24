package tes

import (
  "cloud.google.com/go/datastore"
)

type kvprop struct {
  Key string
  Val string
}

func loadMap(ps []datastore.Property, n string, m map[string]string) ([]datastore.Property, error) {
  var auto []datastore.Property
  for _, p := range ps {
    if p.Name == n {
      kvs := p.Value.([]*kvprop)
      for _, kv := range kvs {
        m[kv.Key] = kv.Val
      }
    } else {
      auto = append(auto, p)
    }
  }
  return auto, nil
}

func saveMap(m map[string]string, n string) datastore.Property {
  var kvs []*kvprop
  for k, v := range m {
    kvs = append(kvs, &kvprop{k, v})
  }
  return datastore.Property{
    Name: n,
    Value: kvs,
  }
}

func (t *Task) Load(ps []datastore.Property) error {
  auto, _ := loadMap(ps, "Tags", t.Tags)
  return datastore.LoadStruct(t, auto)
}
func (t *Task) Save() ([]datastore.Property, error) {
  return datastore.SaveStruct(t)
  // TODO props = append(props, saveMap(t.Tags, "Tags"))
  //return props, nil
}

func (e *Executor) Load(ps []datastore.Property) error {
  auto, _ := loadMap(ps, "Environ", e.Environ)
  return datastore.LoadStruct(e, auto)
}
func (e *Executor) Save() ([]datastore.Property, error) {
  props, _ := datastore.SaveStruct(e)
  // TODO props = append(props, saveMap(e.Environ, "Environ"))
  return props, nil
}

func (t *TaskLog) Load(ps []datastore.Property) error {
  auto, _ := loadMap(ps, "Metadata", t.Metadata)
  return datastore.LoadStruct(t, auto)
}
func (t *TaskLog) Save() ([]datastore.Property, error) {
  props, _ := datastore.SaveStruct(t)
  // TODO props = append(props, saveMap(t.Metadata, "Metadata"))
  return props, nil
}


func (r *Resources) Load(ps []datastore.Property) error {
  var auto []datastore.Property
  for _, p := range ps {
    if p.Name == "CpuCores" {
      r.CpuCores = uint32(p.Value.(int64))
    } else {
      auto = append(auto, p)
    }
  }
  return datastore.LoadStruct(r, auto)
}

func (r *Resources) Save() ([]datastore.Property, error) {
  props, _ := datastore.SaveStruct(r)
  props = append(props, datastore.Property{
    Name: "CpuCores",
    Value: int64(r.CpuCores),
  })
  return props, nil
}

func (p *Ports) Save() ([]datastore.Property, error) {
  return []datastore.Property{
    {
      Name: "Container",
      Value: int64(p.Container),
    },
    {
      Name: "Host",
      Value: int64(p.Host),
    },
  }, nil
}

func (p *Ports) Load(ps []datastore.Property) error {
  for _, prop := range ps {
    switch prop.Name {
    case "Container":
      p.Container = uint32(prop.Value.(int64))
    case "Host":
      p.Host = uint32(prop.Value.(int64))
    }
  }
  return nil
}
