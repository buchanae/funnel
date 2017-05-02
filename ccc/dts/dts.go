package dts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var log = logger.New("CCC DTS Client")

// NewClient returns a new HTTP client for accessing
// Create/List/Get/Cancel Task endpoints. "address" is the address
// of the CCC Central Function server.
func NewClient(address string) (*Client, error) {
	// Strip trailing slash. A quick and dirty fix.
	address = strings.TrimSuffix(address, "/")
	u, err := url.Parse(address)
	if err != nil {
		log.Error("Error parsing URL", err)
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		err := fmt.Errorf("Invalid URL scheme: { %s }", u.Scheme)
		log.Error("Error parsing URL", err)
		return nil, err
	}
	if u.Host != "" {
		err := fmt.Errorf("Invalid host: { %s }", u.Host)
		log.Error("Error parsing URL", err)
		return nil, err
	}
	if u.Path != "" {
		err := fmt.Errorf("Invalid path: { %s }", u.Path)
		log.Error("Error parsing URL", err)
		return nil, err
	}
	c := &Client{
		address: u.String(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return c, nil
}

// Client represents the HTTP Task client.
type Client struct {
	address string
	client  *http.Client
}

// Record represents a DTS record
type Record struct {
	ID       string
	Name     string
	Size     int64
	Location []Location
}

// Location represents a DTS location
type Location struct {
	Site             string
	Path             string
	TimestampUpdated time.Time
	User             struct {
		Name string
	}
}

// GetFile returns the raw bytes from GET /api/v1/dts/file/<id>
func (c *Client) GetFile(id string) (*Record, error) {
	// Send request
	u := c.address + "/api/v1/dts/file/" + id
	body, err := CheckHTTPResponse(c.client.Get(u))
	if err != nil {
		return nil, err
	}
	// Parse response
	resp := &Record{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PostFile returns the raw bytes from POST /api/v1/dts/file
func (c *Client) PostFile(msg []byte) error {
	err := isRecord(msg)
	if err != nil {
		return fmt.Errorf("Not a valid DTS Record message: %v", err)
	}

	// Send request
	r := bytes.NewReader(msg)
	u := c.address + "/api/v1/dts/file"
	_, err = CheckHTTPResponse(c.client.Post(u, "application/json", r))
	if err != nil {
		return err
	}
	return nil
}

// GenerateRecord is a helper method for creating Records
func GenerateRecord(path string, site string) (*Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := &Record{
		ID:   path,
		Name: fi.Name(),
		Size: fi.Size(),
		Location: []Location{
			{
				Site:             site,
				Path:             filepath.Dir(path),
				TimestampUpdated: fi.ModTime(),
			},
		},
	}
	return r, nil
}

// CheckHTTPResponse does some basic error handling
// and reads the response body into a byte array
func CheckHTTPResponse(resp *http.Response, err error) ([]byte, error) {
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

// TODO replace with proper message validation
func isRecord(b []byte) error {
	var js Record
	return json.Unmarshal(b, &js)
}
