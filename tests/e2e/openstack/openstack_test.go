package openstack

import (
	"flag"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/tests/e2e"
	"os"
	"testing"
)

var fun *e2e.Funnel
var confPath = flag.String("openstack-e2e-config", "", "OpenStack end-to-end test config file")
var log = logger.New("e2e-openstack")

func TestMain(m *testing.M) {
	log.Configure(logger.DebugConfig())
	flag.Parse()

	if *confPath == "" {
		log.Info("Skipping openstack e2e tests, no config")
		os.Exit(0)
	}

	conf := e2e.DefaultConfig()
	if err := config.ParseFile(*confPath, &conf); err != nil {
		panic(err)
	}

	fun = e2e.NewFunnel(conf)
	fun.WithLocalBackend()
	fun.StartServer()

	os.Exit(m.Run())
}

func TestSwiftStorage(t *testing.T) {
	id := fun.Run(`
    --cmd "sh -c 'md5sum $in'"
    -i in=swift://buchanan-scratch/funnel
  `)
	task := fun.Wait(id)

	expect := "da385a552397a4ac86ee6444a8f9ae3e  /opt/funnel/inputs/buchanan-scratch/funnel\n"

	if task.Logs[0].Logs[0].Stdout != expect {
		t.Fatal("Missing stdout")
	}
}
