package rpc

import (
	"golang.org/x/net/context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

type Config struct {
	Host     string
	Port     int
	Insecure bool
	Cert     string
	User     string
	// Password for basic auth. with the server APIs.
	Password string
	// Timeout duration for gRPC calls
	Timeout time.Duration
}

func DefaultConfig() Config {
  return Config{
    Host: "localhost",
    Port: 9090,
		Timeout:        time.Second * 5,
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
func (c *Config) Validate() error {
	// TODO something needs to validate the config, to ensure
	return nil
}

func NewConn(ctx context.Context, conf Config) (*grpc.ClientConn, error) {
	// Validate the config
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}

	// If a certificate file was given, load it.
	if conf.Cert != "" {
		creds, err := credentials.NewClientTLSFromFile(conf.Cert, conf.Host)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// Force insecure transport.
	if conf.Insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	// Configure basic auth.
	if conf.Password != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(&loginCreds{conf.User, conf.Password}))
	}

	conn, err := grpc.DialContext(ctx, conf.Address(), opts...)
	if err != nil {
		return nil, fmt.Errorf("couldn't open RPC connection to %s: %s", conf.Address(), err)
	}
	return conn, nil
}

type loginCreds struct {
	user     string
	password string
}

func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	v := base64.StdEncoding.EncodeToString([]byte(c.user + ":" + c.password))

	return map[string]string{
		"Authorization": "Basic " + v,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}
