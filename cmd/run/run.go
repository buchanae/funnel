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

var log = logger.New("run")

type taskvars struct {
	name        string
	workdir     string
	server      string
	image       string
	project     string
	description string
	stdin       string
	stdout      string
	stderr      string
	preemptible bool
	wait        bool
	inputs      []string
	inputDirs   []string
	outputs     []string
	outputDirs  []string
	envVars     []string
	tags        []string
	volumes     []string
	zones       []string
	cpu         int
	ram         float64
	disk        float64
}

var printTask bool

// Scattering and loading extra args is currently only allowed
// at the top level in order to avoid any issues with circular
// includes. If we want this to be per-task, it's possible,
// but more work.
var extra, extraFiles, scatterFiles []string
var vals taskvars

func init() {
	vals.workdir = "/opt/funnel"
	vals.server = "http://localhost:8000"
}

// TODO with input contents, script could be loaded from file

var example = `
    funnel run
    'bowtie2 -f $factor -x $other -p1 $pair1 -p2 $pair2 -o $alignments'
    --container opengenomics/bowtie2:latest
    --name 'Bowtie2 test'
    --description 'Testings an example of using 'funnel run' for a bowtie2 command'
    --in pair1=./pair1.fastq
    --in pair2=./pair2.fastq
    --out alignments=gs://bkt/bowtie/alignments.bam
    --env factor=5
    --vol /tmp 
    --vol /opt
    --cpu 8
    --ram 32
    --disk 100
`

// Cmd represents the run command
var Cmd = &cobra.Command{
	Use:     "run [flags] --container IMAGE CMD",
	Short:   "Run a task.",
	Long:    ``,
	Example: example,
	RunE:    run,
}

func init() {
	f := Cmd.Flags()

	// These flags are separate because they are not allowed
	// in scattered tasks.
	f.BoolVarP(&printTask, "print", "p", printTask, "Print the task, instead of running it")
	f.StringSliceVar(&scatterFiles, "scatter", scatterFiles, "Scatter")
	f.StringSliceVarP(&extra, "extra", "x", extra, "Extra")
	f.StringSliceVarP(&extraFiles, "extra-file", "f", extraFiles, "Extra file")

	// Add per-task flags.
	addTaskFlags(f, &vals)
}

func addTaskFlags(f *pflag.FlagSet, v *taskvars) {
	f.StringVarP(&v.name, "name", "n", v.name, "Task name")
	f.StringVar(&v.description, "description", v.description, "Task description")
	f.StringVar(&v.project, "project", v.project, "Project")
	f.StringVarP(&v.workdir, "workdir", "w", v.workdir, "Set the containter working directory")
	f.StringVarP(&v.image, "container", "c", v.image, "Specify the containter image")
	f.StringSliceVarP(&v.inputs, "in", "i", v.inputs, "A key-value map of input files")
	f.StringSliceVarP(&v.inputDirs, "in-dir", "I", v.inputDirs, "A key-value map of input directories")
	f.StringSliceVarP(&v.outputs, "out", "o", v.outputs, "A key-value map of output files")
	f.StringSliceVarP(&v.outputDirs, "out-dir", "O", v.outputDirs, "A key-value map of output directories")
	f.StringSliceVarP(&v.envVars, "env", "e", v.envVars, "A key-value map of enviromental variables")
	f.StringVar(&v.stdin, "stdin", v.stdin, "File to pass via stdin to the command")
	f.StringVar(&v.stdout, "stdout", v.stdout, "File to write the stdout of the command")
	f.StringVar(&v.stderr, "stderr", v.stderr, "File to write the stderr of the command")
	f.StringSliceVar(&v.volumes, "vol", v.volumes, "Volumes to be defined on the container")
	f.StringSliceVar(&v.tags, "tag", v.tags, "A key-value map of arbitrary tags")
	f.IntVar(&v.cpu, "cpu", v.cpu, "Number of CPUs requested")
	f.Float64Var(&v.ram, "ram", v.ram, "Amount of RAM requested (in GB)")
	f.Float64Var(&v.disk, "disk", v.disk, "Amount of disk space requested (in GB)")
	f.BoolVar(&v.preemptible, "preemptible", v.preemptible, "Allow task to be scheduled on preemptible workers")
	f.StringSliceVar(&v.zones, "zone", v.zones, "Require task be scheduled in certain zones")
	// TODO
	//f.StringVar(&cmdFile, "cmd-file", cmdFile, "Read cmd template from file")
	f.StringVarP(&v.server, "server", "S", v.server, "Address of Funnel server")
	f.BoolVar(&v.wait, "wait", v.wait, "Wait for the task to finish before exiting")
}

func valsToTask(cmd []string, vals taskvars) (*tes.Task, error) {

	if vals.image == "" {
		return nil, fmt.Errorf("You must specify an image to run your command in.")
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

	// check for key collisions
	err = compareKeys(inputFileMap, inputDirMap, outputFileMap, outputDirMap, envVarMap)
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
	environ, err := mergeVars(inputEnvVars, inputDirEnvVars, outputEnvVars, outputDirEnvVars, envVarMap)
	checkErr(err)

	// Build task input parameters
	inputs, err := createTaskParams(inputFileMap, "/opt/funnel/inputs/", tes.FileType_FILE)
	checkErr(err)
	inputDirs, err := createTaskParams(inputDirMap, "/opt/funnel/inputs/", tes.FileType_DIRECTORY)
	checkErr(err)
	inputs = append(inputs, inputDirs...)

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
				ImageName: vals.image,
				Cmd:       cmd,
				Environ:   environ,
				Workdir:   vals.workdir,
				Stdin:     vals.stdin,
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
		return fmt.Errorf("You must specify a command to run.")
	}

	if len(args) > 1 {
		return fmt.Errorf("--in, --out and --env args should have the form 'KEY=VALUE' not 'KEY VALUE'. Extra args: %s", args[1:])
	}

	for _, xf := range extraFiles {
		b, _ := ioutil.ReadFile(xf)
		extra = append(extra, string(b))
	}

	b, _ := ioutil.ReadAll(os.Stdin)
	if len(b) > 0 {
		extra = append(extra, string(b))
	}

	for _, ex := range extra {
		sp, _ := shellquote.Split(ex)
		cmd.ParseFlags(sp)
	}

	rawcmd := args[0]
	executorCmd := []string{"bash", "-c", rawcmd}

	task, err := valsToTask(executorCmd, vals)
	if err != nil {
		return err
	}

	cli := client.NewClient(vals.server)

	if len(scatterFiles) > 0 {
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
				runTask(t, cli)
			}
		}
		return nil
	}

	return runTask(task, cli)
}

func runTask(task *tes.Task, cli *client.Client) error {
	// Marshal message to JSON
	taskJSON, merr := cli.Marshaler.MarshalToString(task)
	if merr != nil {
		return merr
	}

	if printTask {
		fmt.Println(taskJSON)
		return nil
	}

	resp, rerr := cli.CreateTask([]byte(taskJSON))
	if rerr != nil {
		return rerr
	}

	taskID := resp.Id
	fmt.Println(taskID)

	if vals.wait {
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
