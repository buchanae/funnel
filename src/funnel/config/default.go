package config

import (
	log "funnel/logger"
	pbf "funnel/proto/funnel"
	"github.com/ghodss/yaml"
	os_servers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

func WithDebugMode(c config.Config) config.Config {
	c.Server.LogLevel = logger.DebugLevel
	c.Database.LogLevel = logger.DebugLevel
	c.Scheduler.LogLevel = logger.DebugLevel
	c.worker.LogLevel = logger.DebugLevel
	return c
}

func Default() Config {
	return Config{
		HostName:  "localhost",
		DBPath:    "funnel-work-dir/funnel.db",
		HTTPPort:  "8000",
		RPCPort:   "9090",
		WorkDir:   "funnel-work-dir",
		LogLevel:  logger.InfoLevel,
		Scheduler: "local",
		Backends: SchedulerBackends{
			Local: LocalSchedulerBackend{},
			GCE: GCESchedulerBackend{
				Weights: Weights{
					"startup time": 1.0,
				},
				CacheTTL: time.Minute,
			},
		},
		MaxJobLogSize:     10000,
		ScheduleRate:      time.Second,
		ScheduleChunk:     10,
		WorkerPingTimeout: time.Minute,
		WorkerInitTimeout: time.Minute * 5,
		Worker: Worker{
			ServerAddress: "localhost:9090",
			WorkDir:       "funnel-work-dir",
			Timeout:       -1,
			UpdateRate:    time.Second * 5,
			LogUpdateRate: time.Second * 5,
			TrackerRate:   time.Second * 5,
			LogTailSize:   10000,
			LogLevel:      logger.InfoLevel,
			UpdateTimeout: time.Second,
			Resources: &pbf.Resources{
				Disk: 100.0,
			},
		},
		DisableHTTPCache: true,
	}
}
