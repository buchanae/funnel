package config

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
