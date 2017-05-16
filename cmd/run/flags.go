package run

import (
	"github.com/spf13/pflag"
	"strings"
)

// flagVals captures values from CLI flag parsing
type flagVals struct {
	// Top-level flag values. These are not allowed to be redefined
	// by scattered tasks or extra args, to avoid complexity in avoiding
	// circular imports or nested scattering
	printTask    bool
	server       string
	extra        []string
	extraFiles   []string
	scatterFiles []string

	// Per-task flag values. These may be overridden by scattered tasks.
	name        string
	workdir     string
	container   string
	project     string
	description string
	stdin       string
	stdout      string
	stderr      string
	preemptible bool
	wait        bool
	waitFor     []string
	inputs      []string
	inputDirs   []string
	outputs     []string
	outputDirs  []string
	contents    []string
	environ     []string
	tags        []string
	volumes     []string
	zones       []string
	cpu         int
	ram         float64
	disk        float64
	cmd         []string
}

func addTopLevelFlags(f *pflag.FlagSet, v *flagVals) {
	// These flags are separate because they are not allowed
	// in scattered tasks.
	//
	// Scattering and loading extra args is currently only allowed
	// at the top level in order to avoid any issues with circular
	// includes. If we want this to be per-task, it's possible,
	// but more work.
	f.StringVarP(&v.server, "server", "S", v.server, "")
	f.BoolVarP(&v.printTask, "print", "p", v.printTask, "")
	f.StringSliceVarP(&v.extra, "extra", "x", v.extra, "")
	f.StringSliceVarP(&v.extraFiles, "extra-file", "X", v.extraFiles, "")
	f.StringSliceVar(&v.scatterFiles, "scatter", v.scatterFiles, "")

	// Add per-task flags.
	addTaskFlags(f, v)
}

func addTaskFlags(f *pflag.FlagSet, v *flagVals) {
	// General
	f.StringVarP(&v.container, "container", "c", v.container, "")
	f.StringVarP(&v.workdir, "workdir", "w", v.workdir, "")

	// Input/output
	f.StringSliceVarP(&v.inputs, "in", "i", v.inputs, "")
	f.StringSliceVarP(&v.inputDirs, "in-dir", "I", v.inputDirs, "")
	f.StringSliceVarP(&v.outputs, "out", "o", v.outputs, "")
	f.StringSliceVarP(&v.outputDirs, "out-dir", "O", v.outputDirs, "")
	f.StringVar(&v.stdin, "stdin", v.stdin, "")
	f.StringVar(&v.stdout, "stdout", v.stdout, "")
	f.StringVar(&v.stderr, "stderr", v.stderr, "")
	f.StringSliceVarP(&v.contents, "contents", "C", v.contents, "")

	// Resoures
	f.IntVar(&v.cpu, "cpu", v.cpu, "")
	f.Float64Var(&v.ram, "ram", v.ram, "")
	f.Float64Var(&v.disk, "disk", v.disk, "")
	f.StringSliceVar(&v.zones, "zone", v.zones, "")
	f.BoolVar(&v.preemptible, "preemptible", v.preemptible, "")

	// Other
	f.StringVarP(&v.name, "name", "n", v.name, "")
	f.StringVar(&v.description, "description", v.description, "")
	f.StringVar(&v.project, "project", v.project, "")
	f.StringSliceVar(&v.volumes, "vol", v.volumes, "")
	f.StringSliceVar(&v.tags, "tag", v.tags, "")
	f.StringSliceVarP(&v.environ, "env", "e", v.environ, "")

	// TODO
	//f.StringVar(&cmdFile, "cmd-file", cmdFile, "Read cmd template from file")
	f.BoolVar(&v.wait, "wait", v.wait, "")
	f.StringSliceVar(&v.waitFor, "wait-for", v.waitFor, "")
}

// Set default flagVals
func defaultVals(vals *flagVals) {
	if vals.workdir == "" {
		vals.workdir = "/opt/funnel"
	}

	if vals.container == "" {
		vals.container = ""
	}

	// Default name
	if vals.name == "" {
		vals.name = "Funnel run: " + strings.Join(vals.cmd, " ")
	}

	if vals.server == "" {
		vals.server = "http://localhost:8000"
	}
}
