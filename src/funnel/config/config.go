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
type Weights map[string]float32

// StorageConfig describes configuration for all storage types
type StorageConfig struct {
	Local LocalStorage
	S3    S3Storage
	GS    GSStorage
}

// LocalStorage describes the directories Funnel can read from and write to
type LocalStorage struct {
	AllowedDirs []string
}

// GSStorage describes configuration for the Google Cloud storage backend.
type GSStorage struct {
	AccountFile string
	FromEnv     bool
}

// Valid validates the GSStorage configuration.
func (g GSStorage) Valid() bool {
	return g.FromEnv || g.AccountFile != ""
}

// Valid validates the LocalStorage configuration
func (l LocalStorage) Valid() bool {
	return len(l.AllowedDirs) > 0
}

// S3Storage describes the directories Funnel can read from and write to
type S3Storage struct {
	Endpoint string
	Key      string
	Secret   string
}

// Valid validates the LocalStorage configuration
func (l S3Storage) Valid() bool {
	return l.Endpoint != "" && l.Key != "" && l.Secret != ""
}

// LocalSchedulerBackend describes configuration for the local scheduler.
type LocalSchedulerBackend struct {
	Weights Weights
}

// OpenStackSchedulerBackend describes configuration for the openstack scheduler.
type OpenStackSchedulerBackend struct {
	KeyPair    string
	ConfigPath string
	Server     os_servers.CreateOpts
	Weights    Weights
}

// GCESchedulerBackend describes configuration for the Google Cloud scheduler.
type GCESchedulerBackend struct {
	AccountFile string
	Project     string
	Zone        string
	Weights     Weights
	CacheTTL    time.Duration
}

// SchedulerBackends describes configuration for all schedulers.
type SchedulerBackends struct {
	Local     LocalSchedulerBackend
	Condor    LocalSchedulerBackend
	OpenStack OpenStackSchedulerBackend
	GCE       GCESchedulerBackend
	LogPath       string
}

type Server struct {
	HTTPPort      string
	RPCPort       string
	DisableHTTPCache  bool
	LogPath       string
}

type Scheduler struct {
  Backend string
	Backends      SchedulerBackends
	MaxJobLogSize int
	ScheduleRate  time.Duration
	ScheduleChunk int
	// How long to wait for a worker ping before marking it as dead
	WorkerPingTimeout time.Duration
	// How long to wait for worker initialization before marking it dead
	WorkerInitTimeout time.Duration
	LogPath       string
}

type Database struct {
  Path string
	LogPath       string
}

// Config describes configuration for Funnel.
type Config struct {
	worker        Worker
  Server        Server
  Scheduler     Scheduler
  Database      Database
	WorkDir       string
}

func (c Config) Worker() Worker {
  if c.worker.ServerAddress == "" {
    c.worker.ServerAddress = c.Server.HostName + ":" + c.Server.RPCPort
  }
  return c.worker
}

// HTTPAddress returns the HTTP address based on HostName and HTTPPort
func (c Server) HTTPAddress() string {
	return "http://" + c.HostName + ":" + c.HTTPPort
}

// RPCAddress returns the RPC address based on HostName and RPCPort
func (c Server) RPCAddress() string {
	return c.HostName + ":" + c.RPCPort
}

// Worker contains worker configuration.
type Worker struct {
	ID string
	// Address of the scheduler, e.g. "1.2.3.4:9090"
  ServerAddress string
	// Directory to write task files to
	WorkDir string
	// How long (seconds) to wait before tearing down an inactive worker
	// Default, -1, indicates to tear down the worker immediately after completing
	// its task
	Timeout time.Duration
	// How often the worker sends update requests to the server
	UpdateRate time.Duration
	// How often the worker sends task log updates
	LogUpdateRate time.Duration
	TrackerRate   time.Duration
	LogTailSize   int64
	Storage       []*StorageConfig
	LogPath       string
	LogLevel      string
	Resources     *pbf.Resources
	// Timeout duration for UpdateWorker() and UpdateTaskLogs() RPC calls
	UpdateTimeout time.Duration
	Metadata      map[string]string
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
