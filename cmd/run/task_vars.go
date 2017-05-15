package run

func defaultVars() *taskvars {
  return &taskvars{
    workdir: "/opt/funnel",
    container: "alpine",
  }

	// Default name
	if vals.name == "" {
		vals.name = "Funnel run: " + strings.Join(vals.cmd, " ")
	}

	var stdin string
	if vals.stdin != "" {
		stdin = "/opt/funnel/inputs/stdin"
	}
  vars.server = "http://localhost:8000"
}


func (vals *taskvars) Task() *tes.Task {
	// Build the task message
	return &tes.Task{
		Name:        vals.name,
		Project:     vals.project,
		Description: vals.description,
		Inputs:      inputs,
		Outputs:     outputs,
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
				Environ:   environ,
				Workdir:   vals.workdir,
				Stdin:     stdin,
				Stdout:    "/opt/funnel/outputs/stdout",
				Stderr:    "/opt/funnel/outputs/stderr",
				// TODO no ports
				Ports: nil,
			},
		},
		Volumes: vals.volumes,
		Tags:    tagsMap,
	}
}
