package run

import (
	"bufio"
	"fmt"
	"github.com/kballard/go-shellquote"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"strings"
)

// cmdvars capture values from CLI flag parsing
type flagVals struct {
  // Top-level flag values. These are not allowed to be redefined
  // by scattered tasks or extra args, to avoid complexity in avoiding
  // circular imports or nested scattering
  printTask bool
  server string
  extra []string
  extraFiles []string
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
	env         []string
	tags        []string
	volumes     []string
	zones       []string
	cpu         int
	ram         float64
	disk        float64
  cmd         []string
}

func addTopLevelFlags(f *pflag.FlagSet, v *cmdvars) {
	// These flags are separate because they are not allowed
	// in scattered tasks.
  //
  // Scattering and loading extra args is currently only allowed
  // at the top level in order to avoid any issues with circular
  // includes. If we want this to be per-task, it's possible,
  // but more work.
	f.StringVarP(&server, "server", "S", v.server, "")
	f.BoolVarP(&printTask, "print", "p", v.printTask, "")
	f.StringSliceVarP(&extra, "extra", "x", v.extra, "")
	f.StringSliceVarP(&extraFiles, "extra-file", "X", v.extraFiles, "")
	f.StringSliceVar(&scatterFiles, "scatter", v.scatterFiles, "")

	// Add per-task flags.
	addTaskFlags(f, &v.task)
}

func addTaskFlags(f *pflag.FlagSet, v *taskvars) {
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
	f.StringSliceVarP(&v.envVars, "env", "e", v.envVars, "")

	// TODO
	//f.StringVar(&cmdFile, "cmd-file", cmdFile, "Read cmd template from file")
	f.BoolVar(&v.wait, "wait", v.wait, "")
	f.StringSliceVar(&v.waitFor, "wait-for", v.waitFor, "")
}
