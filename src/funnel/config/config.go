package config

import (
	"funnel/logger"
	pbf "funnel/proto/funnel"
	"github.com/ghodss/yaml"
	os_servers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

// Weights describes the scheduler score weights.
// All fields should be float32 type.
type Weight float32

// Config describes configuration for Funnel.
type Config struct {
	Root string

	Database struct {
		// Path to database file
		Path     string
		LogLevel string
		LogPath  string
	}

	Server struct {
		// TODO discoverable?
		HostName string
		// Port which HTTP traffic will be served on.
		HTTPPort string
		// Port which RPC traffic will be served on.
		RPCPort string
		// Normally "Cache-Control: no-store" header is set on HTTP responses,
		// which disables caching by proxies. This makes setup simpler, but
		// when you want to allow cached responses, you can disable that header.
		DisableCacheBuster bool
		EnableHTTPCache    bool
		LogLevel           string
		LogPath            string
	}

	Worker struct {
		ID string
		// Address of the server's RPC endpoint, e.g. "funnel.com:9090"
		ServerAddress string
		// Directory to write job files to
		WorkDir string
		// When the worker has been idle for this duration, it will shut down.
		// -1 means the worker will never shut down.
		Timeout time.Duration
		// How often the worker syncs with the server.
		UpdateRate time.Duration `advanced`
		// How often the worker sends log updates to the server.
		LogUpdateRate time.Duration `advanced`
		// The tails of stdout and stderr logs are streamed back from
		// the worker to the server. This the max. size stored for
		// each task, in bytes.
		LogTailSize size.Bytes `advanced`
		LogPath     string
		LogLevel    string
		// TODO override or minimum?
		Resources struct {
			Disk size.Bytes
		}
		// How long to wait for RPC calls before timing out.
		RPCTimeout time.Duration `advanced`
		// TODO used?
		Metadata map[string]string
	}

	Storage struct {
		Local struct {
			AllowedDirs []string
		}
		// TODO allow multiple S3
		S3 struct {
			Endpoint string
			Key      string
			Secret   string
		}
		// TODO allow multiple GS
		GS struct {
			CredentialsFile string
			FromEnv         bool
		}
	}

	Scheduler struct {
		// The name of the active scheduler backend.
		// Options: local, htcondor, gce, openstack
		Backend string

		// Scheduler backends provide support for scheduling tasks
		// on a variety of infrastructures.
		Backends struct {

			// The local scheduler backend runs tasks on the local computer
			// (in the same process as the server/scheduler)
			Local struct{}

			// HTCondor scheduler backend
			// https://research.cs.wisc.edu/htcondor/
			HTCondor struct {
				// The scheduler will stage some per-task files in this directory,
				// e.g. such as the worker config file, HTCondor submission script,
				// logs, and more.
				WorkDir string
				// It's usually preferred that HTCondor workers only handle a single
				// task then then immediately shut down, so that workers don't stay
				// in the HTCondor queue, taking up valuable resources.
				//
				// This can disable that behavior.
				ReuseWorkers bool
			}

			// The OpenStack scheduler backend
			// https://www.openstack.org/software/
			// TODO
			OpenStack struct {
				KeyPair    string
				ConfigPath string
				Server     os_servers.CreateOpts
			}

			// Google Compute Engine (GCE) scheduler backend
			// https://cloud.google.com/compute/
			GCE struct {
				// GCE project name.
				//
				// See the GCE docs for more information:
				// https://cloud.google.com/docs/authentication//projects_and_resources
				CredentialsFile string
				// Path to a GCE credentials file.
				//
				// This is only necessary when running Funnel outside of GCE.
				// When running on a GCE VM, the credentials can be automatially detected.
				//
				// See the GCE docs for more information:
				// https://cloud.google.com/docs/authentication//getting_credentials_for_server-centric_flow
				Project string
				// See the GCE docs for more information:
				// https://cloud.google.com/compute/docs/regions-zones/regions-zones
				Zone string
				// Weights provide tuning of the scheduling policies
				// by weighting the score for each possible worker.
				//
				// Weights range from 0.0 (no effect) to 1.0 (full effect).
				Weights struct {
					// Workers with a lower startup time will be preferred.
					// Workers which are already running have the lowest startup time.
					PreferLowStartupTime Weight
				}
				// How long before cached GCE metadata is expired.
				//
				// The GCE backend needs to query GCE APIs to look for instance templates,
				// machine types, and other metadata. This info is temporarily cached for efficiency.
				CacheTTL time.Duration
			}
		}
		LogLevel string
		LogPath  string
		// The tails of stdout and stderr logs are streamed back from
		// the worker to the server. This the max. size stored for
		// each task, in bytes.
		// TODO duplicated with Worker
		MaxJobLogSize int `advanced`
		// How often to schedule a chunk of tasks.
		Rate time.Duration `advanced`
		// How many tasks to schedule with each iteration.
		ChunkSize int `advanced`
		// How long between pings before a worker is considered dead.
		WorkerPingTimeout time.Duration `advanced`
		// How long to wait for a worker to start before it's considered dead.
		WorkerInitTimeout time.Duration `advanced`
	}
}

// HTTPAddress returns the HTTP address based on HostName and HTTPPort
func (c Config) HTTPAddress() string {
	return "http://" + c.Server.HostName + ":" + c.Server.HTTPPort
}

// RPCAddress returns the RPC address based on HostName and RPCPort
func (c Config) RPCAddress() string {
	return c.Server.HostName + ":" + c.Server.RPCPort
}

func WithDebug(c config.Config) Config {
	c.Server.LogLevel = logger.DebugLevel
	c.Database.LogLevel = logger.DebugLevel
	c.Scheduler.LogLevel = logger.DebugLevel
	c.Worker.LogLevel = logger.DebugLevel
	return c
}

// DefaultConfig returns configuration with simple defaults.
func WithDefaults(c config.Config) Config {
	srv := &c.Server
	db := &c.Database
	sched := &c.Scheduler
	store := &c.Storage
	w := &c.Worker

	srv.HostName = "localhost"
	srv.RPCPort = 9090
	srv.HTTPPort = 8000

	db.Path = "./funnel-work-dir/funnel.db"
	db.LogLevel = logger.InfoLevel

	sched.Backend = "local"

	gce := &sched.Backends.GCE
	gce.Weights.PreferLowStartupTime = 1.0
	gce.Zone = "us-west-1a"
	gce.CacheTTL = time.Minute

	sched.CacheTTL = time.Minute
	sched.LogLevel = logger.InfoLevel
	sched.MaxJobLogSize = size.KB * 10
	sched.Rate = time.Second
	sched.ChunkSize = 10
	sched.WorkerPingTimeout = time.Minute
	sched.WorkerInitTimeout = time.Minute * 5

	w.WorkDir = "./funnel-work-dir"
	w.Timeout = Never
	w.UpdateRate = time.Second * 5
	w.UpdateTimeout = time.Second
	w.LogUpdateRate = time.Second * 5
	w.LogLevel = logger.InfoLevel
	w.LogTailSize = size.KB * 10
	w.Resources.Disk = 100.0

	return c
}

func InheritWorkerConfig(c Config) Config {
	if c.Worker.ServerAddress == "" {
		c.Worker.ServerAddress = c.RPCAddress()
	}
	return c
}

// ToYaml formats the configuration into YAML and returns the bytes.
func (c Config) ToYaml() []byte {
	// TODO handle error
	yamlstr, _ := yaml.Marshal(c)
	return yamlstr
}

// Parse parses a YAML doc into the given Config instance.
func Parse(raw []byte, conf *Config) error {
	err := yaml.Unmarshal(raw, conf)
	if err != nil {
		return err
	}
	return nil
}

// ParseFile parses a Funnel config file, which is formatted in YAML,
// and returns a Config struct.
func ParseFile(relpath string, conf *Config) error {
	if relpath == "" {
		return nil
	}

	// Try to get absolute path. If it fails, fall back to relative path.
	path, abserr := filepath.Abs(relpath)
	if abserr != nil {
		path = relpath
	}

	// Read file
	source, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("Failure reading config", "path", path, "error", err)
		return err
	}

	// Parse file
	perr := Parse(source, conf)
	if perr != nil {
		logger.Error("Failure reading config", "path", path, "error", perr)
		return perr
	}
	return nil
}
