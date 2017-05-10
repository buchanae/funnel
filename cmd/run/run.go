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

// Cmd represents the run command
var Cmd = &cobra.Command{
	Use:   "run 'CMD' [flags]",
	Short: "Run a task.",
	RunE:  run,
}

var log = logger.New("run")

type taskvars struct {
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
	envVars     []string
	tags        []string
	volumes     []string
	zones       []string
	cpu         int
	ram         float64
	disk        float64
}

var printTask bool
var server = "http://localhost:8000"

// Scattering and loading extra args is currently only allowed
// at the top level in order to avoid any issues with circular
// includes. If we want this to be per-task, it's possible,
// but more work.
var extra, extraFiles, scatterFiles []string
var vals taskvars

func init() {
	vals.workdir = "/opt/funnel"
	vals.container = "alpine"
}

//
// Usage/help docs are defined in usage.go.
// If you're updating flags, you probably need to update that file.
//

func init() {
	f := Cmd.Flags()

	// These flags are separate because they are not allowed
	// in scattered tasks.
	f.StringVarP(&server, "server", "S", server, "")
	f.BoolVarP(&printTask, "print", "p", printTask, "")
	f.StringSliceVarP(&extra, "extra", "x", extra, "")
	f.StringSliceVarP(&extraFiles, "extra-file", "X", extraFiles, "")
	f.StringSliceVar(&scatterFiles, "scatter", scatterFiles, "")

	// Add per-task flags.
	addTaskFlags(f, &vals)
	Cmd.SetUsageTemplate(usage)
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

func valsToTask(cmd []string, vals taskvars) (*tes.Task, error) {

	if vals.container == "" {
		return nil, fmt.Errorf("you must specify a container")
	}

	// Get template variables from the command line.
	inputFileMap, err := parseCliVars(vals.inputs)
	checkErr(err)
	inputDirMap, err := parseCliVars(vals.inputDirs)
	checkErr(err)
	outputFileMap, err := parseCliVars(vals.outputs)
	checkErr(err)
	outputDirMap, err := parseCliVars(vals.outputDirs)
	checkErr(err)
	envVarMap, err := parseCliVars(vals.envVars)
	checkErr(err)
	tagsMap, err := parseCliVars(vals.tags)
	checkErr(err)
	contentsMap, err := parseCliVars(vals.contents)
	checkErr(err)

	// check for key collisions
	err = compareKeys(inputFileMap, inputDirMap, outputFileMap, outputDirMap, envVarMap, contentsMap)
	checkErr(err)

	// Create map of enviromental variables to be passed to the executor
	inputEnvVars, err := fileMapToEnvVars(inputFileMap, "/opt/funnel/inputs/")
	checkErr(err)
	inputDirEnvVars, err := fileMapToEnvVars(inputDirMap, "/opt/funnel/inputs/")
	checkErr(err)
	outputEnvVars, err := fileMapToEnvVars(outputFileMap, "/opt/funnel/outputs/")
	checkErr(err)
	outputDirEnvVars, err := fileMapToEnvVars(outputDirMap, "/opt/funnel/outputs/")
	checkErr(err)
	contentsEnvVars, err := fileMapToEnvVars(contentsMap, "/opt/funnel/inputs/")
	environ, err := mergeVars(inputEnvVars, inputDirEnvVars, outputEnvVars, outputDirEnvVars, envVarMap, contentsEnvVars)
	checkErr(err)

	// Build task input parameters
	inputs, err := createTaskParams(inputFileMap, "/opt/funnel/inputs/", tes.FileType_FILE)
	checkErr(err)
	inputDirs, err := createTaskParams(inputDirMap, "/opt/funnel/inputs/", tes.FileType_DIRECTORY)
	checkErr(err)
	contentsParams, err := createContentsParams(contentsMap, "/opt/funnel/inputs/")
	checkErr(err)
	inputs = append(inputs, inputDirs...)
	inputs = append(inputs, contentsParams...)

	// Build task output parameters
	outputs, err := createTaskParams(outputFileMap, "/opt/funnel/outputs/", tes.FileType_FILE)
	checkErr(err)
	outputDirs, err := createTaskParams(outputDirMap, "/opt/funnel/outputs/", tes.FileType_DIRECTORY)
	checkErr(err)
	outputs = append(outputs, outputDirs...)

	// Default name
	if vals.name == "" {
		vals.name = "Funnel run: " + strings.Join(cmd, " ")
	}

	var stdin string
	if vals.stdin != "" {
		stdin = "/opt/funnel/inputs/stdin"
	}

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
				Cmd:       cmd,
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
	}, nil
}

func run(cmd *cobra.Command, args []string) error {

	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("you must specify a command to run")
	}

	if len(args) > 1 {
		cmd.Usage()
		return fmt.Errorf("--in, --out and --env args should have the form 'KEY=VALUE' not 'KEY VALUE'. Extra args: %s", args[1:])
	}

	for _, xf := range extraFiles {
		b, _ := ioutil.ReadFile(xf)
		extra = append(extra, string(b))
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, _ := ioutil.ReadAll(os.Stdin)
		if len(b) > 0 {
			extra = append(extra, string(b))
		}
	}

	for _, ex := range extra {
		sp, _ := shellquote.Split(ex)
		cmd.ParseFlags(sp)
	}

	executorCmd, _ := shellquote.Split(args[0])

	task, err := valsToTask(executorCmd, vals)
	if err != nil {
		cmd.Usage()
		return err
	}

	cli := client.NewClient(server)

	if len(scatterFiles) > 0 {
		tg := taskGroup{}
		for _, sc := range scatterFiles {
			fh, _ := os.Open(sc)
			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				// Copy task vals to act as base task
				tv := vals
				// Per-scatter flags
				fl := &pflag.FlagSet{}
				addTaskFlags(fl, &tv)
				sp, _ := shellquote.Split(scanner.Text())
				fl.Parse(sp)
				t, _ := valsToTask(executorCmd, tv)
				tg.runTask(t, cli, tv.wait, tv.waitFor)
			}
		}
		return tg.wait()
	}

	return runTask(task, cli, vals.wait, vals.waitFor)
}

func runTask(task *tes.Task, cli *client.Client, wait bool, waitFor []string) error {
	// Marshal message to JSON
	taskJSON, merr := cli.Marshaler.MarshalToString(task)
	if merr != nil {
		return merr
	}

	if printTask {
		fmt.Println(taskJSON)
		return nil
	}

	if len(waitFor) > 0 {
		for _, tid := range waitFor {
			cli.WaitForTask(tid)
		}
	}

	resp, rerr := cli.CreateTask([]byte(taskJSON))
	if rerr != nil {
		return rerr
	}

	taskID := resp.Id
	fmt.Println(taskID)

	if wait {
		return cli.WaitForTask(taskID)
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}
