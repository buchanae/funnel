package util

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/rpc"
	"strings"
  "strconv"
)

// MergeConfigFileWithFlags is a util used by server commands that use flags to set
// Funnel config values. These commands can also take in the path to a Funnel config file.
// This function ensures that the config gets set up properly. Flag values override values in
// the provided config file.
func MergeConfigFileWithFlags(file string, flagConf config.Config) (config.Config, error) {
	// parse config file if it exists
	conf := config.DefaultConfig()
	err := config.ParseFile(file, &conf)
	if err != nil {
		return conf, err
	}

	// make sure server address and password is inherited by scheduler nodes and workers
	conf = config.EnsureServerProperties(conf)
	flagConf = config.EnsureServerProperties(flagConf)

	// file vals <- cli val
	err = mergo.MergeWithOverwrite(&conf, flagConf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

// ParseRPCAddress parses a gRPC address and sets the relevant
// fields in the config
func ParseRPCAddress(address string, conf rpc.Config) (rpc.Config, error) {
	if address != "" {
		parts := strings.Split(address, ":")
		if len(parts) != 2 {
			return conf, fmt.Errorf("error parsing server address")
		}
		conf.Host = parts[0]
    p, err := strconv.Atoi(parts[1])
    if err != nil {
			return conf, fmt.Errorf("error parsing server address")
    }
		conf.Port = p
	}
	return conf, nil
}
