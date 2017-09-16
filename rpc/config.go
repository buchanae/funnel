package rpc

import (
  "time"
)

type Config struct {
  Host, Port, Password string
	// Timeout duration for gRPC calls
  Timeout time.Duration
}
func (c Config) Address() string {
	if c.Host != "" && c.Port != "" {
		return c.Host + ":" + c.Port
	}
	return ""
}
