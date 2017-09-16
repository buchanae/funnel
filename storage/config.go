package storage

// Config describes configuration for all storage types
type Config struct {
	Local LocalConfig
	S3    []S3Config
	GS    []GSConfig
}

// LocalConfig describes the directories Funnel can read from and write to
type LocalConfig struct {
	AllowedDirs []string
}

// Valid validates the LocalConfig configuration
func (l LocalConfig) Valid() bool {
	return len(l.AllowedDirs) > 0
}

// GSConfig describes configuration for the Google Cloud storage backend.
type GSConfig struct {
	AccountFile string
	FromEnv     bool
}

// Valid validates the GSConfig configuration.
func (g GSConfig) Valid() bool {
	return g.FromEnv || g.AccountFile != ""
}

// S3Config describes the directories Funnel can read from and write to
type S3Config struct {
	Endpoint string
	Key      string
	Secret   string
}

// Valid validates the S3Config configuration
func (l S3Config) Valid() bool {
	return l.Endpoint != "" && l.Key != "" && l.Secret != ""
}
