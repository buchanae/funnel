package testutils

import (
	"bytes"
	"context"
	"fmt"
	dockerTypes "github.com/docker/docker/api/types"
	dockerFilters "github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	runlib "github.com/ohsu-comp-bio/funnel/cmd/run"
	"github.com/ohsu-comp-bio/funnel/cmd/server"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"google.golang.org/grpc"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"text/template"
	"time"
)

func init() {
	// nanoseconds are important because the tests run faster than a millisecond
	// which can cause port conflicts
	rand.Seed(time.Now().UTC().UnixNano())
	logger.ForceColors()
}

// RandomPort returns a random port string between 10000 and 20000.
func RandomPort() string {
	min := 10000
	max := 20000
	n := rand.Intn(max-min) + min
	return fmt.Sprintf("%d", n)
}

// RandomPortConfig returns a modified config with random HTTP and RPC ports.
func RandomPortConfig(conf config.Config) config.Config {
	conf.RPCPort = RandomPort()
	conf.HTTPPort = RandomPort()
	return conf
}

// TempDirConfig returns a modified config with workdir and db path set to a temp. directory.
func TempDirConfig(conf config.Config) config.Config {
	os.Mkdir("./test_tmp", os.ModePerm)
	f, _ := ioutil.TempDir("./test_tmp", "funnel-test-")
	conf.WorkDir = f
	conf.DBPath = path.Join(f, "funnel.db")
	return conf
}

type Funnel struct {
	// Clients
	RPC    tes.TaskServiceClient
	HTTP   *client.Client
	Docker *docker.Client

	// Config
	Conf       config.Config
	StorageDir string

	// Internal
	startTime string
	rate      time.Duration
}

func NewFunnel() *Funnel {
	var rate = time.Millisecond * 1000
	conf := config.DefaultConfig()
	conf = TempDirConfig(conf)
	conf = RandomPortConfig(conf)
	conf.LogLevel = "debug"
	conf.Worker.LogUpdateRate = rate
	conf.Worker.UpdateRate = rate
	conf.ScheduleRate = rate

	storageDir, _ := ioutil.TempDir("./test_tmp", "funnel-test-storage-")
	wd, _ := os.Getwd()

	// TODO need to fix the storage config so that you can't accidentally
	//      configure both S3 and Local on the same StorageConfig object,
	//      which is not valid.
	conf.Storage = append(conf.Storage,
		&config.StorageConfig{
			Local: config.LocalStorage{
				AllowedDirs: []string{storageDir, wd},
			},
		},
		&config.StorageConfig{
			S3: config.S3Storage{
				Endpoint: "localhost:9999",
				Key:      "AKIAIOSFODNN7EXAMPLE",
				Secret:   "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
		},
	)

	go server.Run(conf)
	time.Sleep(time.Second)

	conn, err := grpc.Dial(conf.RPCAddress(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	var derr error
	dcli, derr := util.NewDockerClient()
	if derr != nil {
		panic(derr)
	}

	return &Funnel{
		RPC:        tes.NewTaskServiceClient(conn),
		HTTP:       client.NewClient("http://localhost:" + conf.HTTPPort),
		Docker:     dcli,
		Conf:       conf,
		StorageDir: storageDir,
		startTime:  fmt.Sprintf("%d", time.Now().Unix()),
		rate:       rate,
	}
}

// wait for a "destroy" event from docker for the given container ID
// TODO probably could use docker.ContainerWait()
// https://godoc.org/github.com/moby/moby/client#Client.ContainerWait
func (f *Funnel) WaitForDockerDestroy(id string) {
	fil := dockerFilters.NewArgs()
	fil.Add("type", "container")
	fil.Add("container", id)
	fil.Add("event", "destroy")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s, err := f.Docker.Events(ctx, dockerTypes.EventsOptions{
		Since:   string(f.startTime),
		Filters: fil,
	})
	for {
		select {
		case e := <-err:
			panic(e)
		case <-s:
			return
		}
	}
}

// cancel a task by ID
func (f *Funnel) Cancel(id string) error {
	_, err := f.RPC.CancelTask(context.Background(), &tes.CancelTaskRequest{
		Id: id,
	})
	return err
}

func (f *Funnel) ListView(view tes.TaskView) []*tes.Task {
	t, err := f.RPC.ListTasks(context.Background(), &tes.ListTasksRequest{
		View: view,
	})
	if err != nil {
		panic(err)
	}
	return t.Tasks
}

func (f *Funnel) GetView(id string, view tes.TaskView) *tes.Task {
	t, err := f.RPC.GetTask(context.Background(), &tes.GetTaskRequest{
		Id:   id,
		View: view,
	})
	if err != nil {
		panic(err)
	}
	return t
}

// get a task by ID
func (f *Funnel) Get(id string) *tes.Task {
	return f.GetView(id, tes.TaskView_FULL)
}

// run a task and return it's ID
func (f *Funnel) Run(s string) string {
	id, err := f.RunE(s)
	if err != nil {
		panic(err)
	}
	return id
}

func (f *Funnel) RunE(s string) (string, error) {
	// Process the string as a template to allow a few helpers
	tpl := template.Must(template.New("run").Parse(s))
	var by bytes.Buffer
	data := map[string]string{
		"storage": "./" + f.StorageDir,
	}
	if eerr := tpl.Execute(&by, data); eerr != nil {
		return "", eerr
	}
	s = by.String()

	tasks, err := runlib.ParseString(s)
	if err != nil {
		return "", err
	}
	if len(tasks) > 1 {
		return "", fmt.Errorf("Funnel run only handles a single task (no scatter)")
	}
	return f.RunTask(tasks[0])
}

func (f *Funnel) RunTask(t *tes.Task) (string, error) {
	resp, cerr := f.RPC.CreateTask(context.Background(), t)
	if cerr != nil {
		return "", cerr
	}
	return resp.Id, nil
}

// wait for a task to complete
func (f *Funnel) Wait(id string) *tes.Task {
	for range time.NewTicker(f.rate).C {
		t := f.Get(id)
		if t.State != tes.State_QUEUED && t.State != tes.State_INITIALIZING &&
			t.State != tes.State_RUNNING {
			return t
		}
	}
	return nil
}

// wait for a task to be in the RUNNING state
func (f *Funnel) WaitForRunning(id string) {
	for range time.NewTicker(f.rate).C {
		t := f.Get(id)
		if t.State == tes.State_RUNNING {
			return
		}
	}
}

// wait for a task to reach the given executor index.
// 1 is the first executor.
func (f *Funnel) WaitForExec(id string, i int) {
	for range time.NewTicker(f.rate).C {
		t := f.Get(id)
		if len(t.Logs[0].Logs) >= i {
			return
		}
	}
}

// write a file to local storage
func (f *Funnel) WriteFile(name string, content string) {
	err := ioutil.WriteFile(f.StorageDir+"/"+name, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// read a file from local storage
func (f *Funnel) ReadFile(name string) string {
	b, err := ioutil.ReadFile(f.StorageDir + "/" + name)
	if err != nil {
		panic(err)
	}
	return string(b)
}
