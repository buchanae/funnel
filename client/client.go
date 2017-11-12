package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
  "io"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Host string
	Port int
	// Use HTTP instead of HTTPS.
	Insecure bool
	// Path to TLS certification. If not provided, the system's certificates will be used.
	Cert string
	// User and Password for basic auth.
	User     string
	Password string
	// Request timeout.
	Timeout time.Duration
}
func DefaultConfig() Config {
  return Config{
    Host: "localhost",
    Port: 8000,
		Timeout:        time.Second * 5,
	}
}

func (c *Config) Address() string {
	proto := "https"
	if c.Insecure {
		proto = "http"
	}
	// Clean trailing slashes.
  host := strings.TrimRight(c.Host, "/") + "/"
	addr := proto + "://" + host

	// Only add the port if necessary.
	if c.Port != 0 && !(proto == "https" && c.Port == 443) && !(proto == "http" && c.Port == 80) {
		addr += fmt.Sprintf(":%d", c.Port)
	}
	return addr
}

func (c *Config) Validate() error {
	// TODO something needs to validate the config, to ensure
	return nil
}

// NewClient returns a new HTTP client for accessing
// Create/List/Get/Cancel Task endpoints. "address" is the address
// of the TES server.
func NewClient(conf Config) (*Client, error) {
	// Validate the config
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	// Load TLS certs from system.
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	// If the config provided a cert file, add it to the pool.
	if conf.Cert != "" {
		b, err := ioutil.ReadFile(conf.Cert)
		if err != nil {
			return nil, err
		}

		ok := pool.AppendCertsFromPEM(b)
		if !ok {
			return nil, fmt.Errorf("failed to parse PEM data from cert: %s", conf.Cert)
		}
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
	}
	tlsConfig.BuildNameToCertificate()

	return &Client{
		client: &http.Client{
			Timeout: conf.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		conf: conf,
	}, nil
}

// Client represents the HTTP Task client.
type Client struct {
	client    *http.Client
	conf      Config
}

func (c *Client) do(ctx context.Context, method, url string, in, out proto.Message) error {
	url = c.conf.Address() + url

	var b io.Reader
	if in != nil {
    by := &bytes.Buffer{}
    b = by
		m := &jsonpb.Marshaler{}
		err := m.Marshal(by, in)
		if err != nil {
			return fmt.Errorf("error marshaling task message: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}

	if in != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if c.conf.Password != "" {
		req.SetBasicAuth(c.conf.User, c.conf.Password)
	}
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if (resp.StatusCode / 100) != 2 {
    body, err := ioutil.ReadAll(resp.Body)
    var b string
    if err != nil {
      b = err.Error()
    } else {
      b = string(body)
    }
		return fmt.Errorf("[STATUS CODE - %d]\t%s", resp.StatusCode, b)
	}

	err = jsonpb.Unmarshal(resp.Body, out)
	if err != nil {
		return err
	}

	return nil
}

// GetTask returns the raw bytes from GET /v1/tasks/{id}
func (c *Client) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
	v := url.Values{}
	addString(v, "view", req.GetView().String())
	task := &tes.Task{}
	err := c.do(ctx, "GET", "/v1/tasks/"+req.Id+"?"+v.Encode(), nil, task)
	return task, err
}

// ListTasks returns the result of GET /v1/tasks
func (c *Client) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {
	// Build url query parameters
	v := url.Values{}
	addString(v, "name_prefix", req.GetNamePrefix())
	addUInt32(v, "page_size", req.GetPageSize())
	addString(v, "page_token", req.GetPageToken())
	addString(v, "view", req.GetView().String())

	list := &tes.ListTasksResponse{}
	err := c.do(ctx, "GET", "/v1/tasks?"+v.Encode(), nil, list)
	return list, err
}

// CreateTask POSTs a Task message to /v1/tasks
func (c *Client) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
	verr := tes.Validate(task)
	if verr != nil {
		return nil, fmt.Errorf("invalid task message: %v", verr)
	}

	resp := &tes.CreateTaskResponse{}
	err := c.do(ctx, "POST", "/v1/tasks", task, resp)
	return resp, err
}

// CancelTask POSTs to /v1/tasks/{id}:cancel
func (c *Client) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {
	resp := &tes.CancelTaskResponse{}
	err := c.do(ctx, "POST", "/v1/tasks/"+req.Id+":cancel", nil, resp)
	return resp, err
}

// GetServiceInfo returns result of GET /v1/tasks/service-info
func (c *Client) GetServiceInfo(ctx context.Context, req *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	resp := &tes.ServiceInfo{}
	err := c.do(ctx, "GET", "/v1/tasks/service-info", nil, resp)
	return resp, err
}

// WaitForTask polls /v1/tasks/{id} for each Id provided and returns
// once all tasks are in a terminal state.
func (c *Client) WaitForTask(ctx context.Context, taskIDs ...string) error {
	for range time.NewTicker(time.Second * 2).C {
		done := false
		for _, id := range taskIDs {
			r, err := c.GetTask(ctx, &tes.GetTaskRequest{
				Id:   id,
				View: tes.TaskView_MINIMAL,
			})
			if err != nil {
				return err
			}
			switch r.State {
			case tes.State_COMPLETE:
				done = true
			case tes.State_EXECUTOR_ERROR, tes.State_SYSTEM_ERROR, tes.State_CANCELED:
				errMsg := fmt.Sprintf("Task %s exited with state %s", id, r.State.String())
				return errors.New(errMsg)
			default:
				done = false
			}
		}
		if done {
			return nil
		}
	}
	return nil
}
