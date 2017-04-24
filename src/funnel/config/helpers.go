package config

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

func InheritWorkerConfig(c Config) Config {
	if c.Worker.ServerAddress == "" {
		c.Worker.ServerAddress = c.RPCAddress()
	}
	return c
}
