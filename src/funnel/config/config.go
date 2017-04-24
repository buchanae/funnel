package config

import (
	os_servers "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"time"
)

// Config describes configuration for Funnel.
type Config struct {
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
		UpdateRate time.Duration
		// How often the worker sends log updates to the server.
		LogUpdateRate time.Duration
		// The tails of stdout and stderr logs are streamed back from
		// the worker to the server. This the max. size stored for
		// each task, in bytes.
		LogTailSize size.Bytes
		LogPath     string
		LogLevel    string
		// TODO override or minimum?
		Resources struct {
			Disk size.Bytes
		}
		// How long to wait for RPC calls before timing out.
		RPCTimeout time.Duration
		// TODO used?
		Metadata map[string]string
	}

	Storage struct {
		// Local filesystem.
		Local []struct {
			// Path to local directory
			Path string
		}
		// S3 (Amazon-specific?
		S3 []struct {
			// S3 server endpoint
			Endpoint string
			// S3 authentication key
			Key string
			// S3 authentication secret
			Secret string
			// Use SSL?
			SSL bool
		}

		// Google Cloud Storage
		GCS []struct {
			// Path to auth. credentials file.
			CredentialsFile string
			// Look for auth. credentials in the environment?
			// TODO probably should always be true.
			FromEnv bool
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
					// Workers with a quicker startup time will be preferred.
					// Workers which are already running have a startup time of 0 (quickest).
					PreferQuickerStartupTime float32
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
		MaxJobLogSize int
		// How often to schedule a chunk of tasks.
		Rate time.Duration
		// How many tasks to schedule with each iteration.
		ChunkSize int
		// How long between pings before a worker is considered dead.
		WorkerPingTimeout time.Duration
		// How long to wait for a worker to start before it's considered dead.
		WorkerInitTimeout time.Duration
	}
}
