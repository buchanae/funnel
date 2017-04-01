package run

import (
  "bytes"
  "text/template"
  "strings"
  "funnel/proto/tes"
  "net/url"
)

type parseResult struct {
  Inputs []*tes.TaskParameter
  Outputs []*tes.TaskParameter
  Volumes []*tes.Volume
  Cmd []string
  // The template functions must return a string
  // so map urls to TaskParameters so other functions
  // can set properties such as name, path, etc.
  inputsMap map[string]*tes.TaskParameter
  outputsMap map[string]*tes.TaskParameter
}

func (res *parseResult) AddInput(rawurl string) string {
  u, _ := url.Parse(rawurl)
  // TODO linux specific path?
  // TODO raw/encoded path is best?
  p := "/opt/funnel/inputs" + u.EscapedPath()
  in := &tes.TaskParameter{
    Location: rawurl,
    Path: p,
    Class: "File",
  }
  res.inputsMap[p] = in
  res.Inputs = append(res.Inputs, in)
  return p
}

func (res *parseResult) AddOutput(rawurl string) string {
  u, _ := url.Parse(rawurl)
  log.Debug("OUTPATH", u.EscapedPath())
  // TODO linux specific path?
  // TODO raw/encoded path is best?
  p := "/opt/funnel/outputs" + u.EscapedPath()
  out := &tes.TaskParameter{
    Location: rawurl,
    Path: p,
    Class: "File",
  }
  res.outputsMap[p] = out
  res.Outputs = append(res.Outputs, out)
  return p
}

func parseTpl(raw string, vars map[string]string) (*parseResult, error) {
  res := &parseResult{
    // TODO
    Volumes: []*tes.Volume{
      {
        Name: "TODO Default funnel run volume",
        // TODO
        SizeGb: 10,
        MountPoint: "/opt/funnel/outputs",
      },
    },
    inputsMap: map[string]*tes.TaskParameter{},
    outputsMap: map[string]*tes.TaskParameter{},
  }
  funcs := template.FuncMap{
    "in": res.AddInput,
    "out": res.AddOutput,
  }

  t, err := template.New("cmd").
    Delims("{", "}").
    Funcs(funcs).
    Parse(raw)

  if err != nil {
    return nil, err
  }

  buf := &bytes.Buffer{}
  xerr := t.Execute(buf, vars)
  if xerr != nil {
    return nil, xerr
  }

  cmd := buf.String()
  log.Debug("CMD", cmd)
  // TODO shell splitting needed
  res.Cmd = strings.Split(cmd, " ")
  return res, nil
}
