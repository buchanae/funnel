package dts

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
  "github.com/ohsu-comp-bio/funnel/logger"
  "time"
  "encoding/json"
  "io/ioutil"
)

var log = logger.New("CCC DTS")

type Client interface {
  GetFile(id string) (*Entry, error)
}

// NewClient returns a new HTTP client for accessing
// Create/List/Get/Cancel Task endpoints. "address" is the address
// of the CCC Central Function server.
func NewClient(address string) (Client, error) {
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Can't parse DTS address", err)
		return nil, err
	}
	if u.Scheme != "http" || u.Scheme != "https" {
		errors.New("Invalid URL scheme.")
		log.Error("Invalid DTS URL scheme", err)
		return nil, err
	}
	c := &client{
		address: u,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return c, nil
}

// Client represents the HTTP Task client.
type client struct {
	address *url.URL
	client  *http.Client
}

// Entry represents a DTS entry
type Entry struct {
	ID       string
	Name     string
	Size     uint32
	Location []Location
}

// Location represents a DTS location
type Location struct {
	Site             string
	Path             string
	TimestampUpdated uint32
	User             struct {
		Name string
	}
}

// Get returns the raw bytes from GET /api/v1/dts/file/<id>
func (c *client) GetFile(id string) (*Entry, error) {
	// Send request
	body, err := check(c.client.Get("api/v1/dts/file/" + id))
	if err != nil {
		return nil, err
	}
	// Parse response
	resp := &Entry{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// check does some basic error handling
// and reads the response body into a byte array
func check(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if (resp.StatusCode / 100) != 2 {
		return nil, fmt.Errorf("[STATUS CODE - %d]\t%s", resp.StatusCode, body)
	}
	return body, nil
}
