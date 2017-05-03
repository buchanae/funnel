package dts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/util"
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
	if u.Host == "" {
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
	CCCID    string     `json:"cccId"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	Location []Location `json:"location"`
}

// Location represents a DTS location
type Location struct {
	Site             string `json:"site"`
	Path             string `json:"path"`
	TimestampUpdated int64  `json:"timestampUpdated"`
	User             struct {
		Name string `json:"name"`
	} `json:"user"`
}

// SitePath returns the absolute path of the file/directory at the specified site
func (r *Record) SitePath(site string) string {
	for _, location := range r.Location {
		if location.Site == site {
			return filepath.Join(location.Path, r.Name)
		}
	}
	return ""
}

// HasSiteLocation returns true if the record contains a Location with this Site
func (r *Record) HasSiteLocation(site string) bool {
	for _, location := range r.Location {
		if location.Site == site {
			return true
		}
	}
	return false
}

// GetFile returns the raw bytes from GET /api/v1/dts/file/<id>
func (c *Client) GetFile(id string) (*Record, error) {
	// convert ID to be URL safe
	cccID := url.PathEscape(id)
	// Send request
	u := c.address + "/api/v1/dts/file/" + cccID
	body, err := util.CheckHTTPResponse(c.client.Get(u))
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
	_, err = util.CheckHTTPResponse(c.client.Post(u, "application/json", r))
	if err != nil {
		return err
	}
	return nil
}

// PutFile returns the raw bytes from PUT /api/v1/dts/file
func (c *Client) PutFile(msg []byte) error {
	err := isRecord(msg)
	if err != nil {
		return fmt.Errorf("Not a valid DTS Record message: %v", err)
	}

	// Send request
	r := bytes.NewReader(msg)
	u := c.address + "/api/v1/dts/file"
	req, err := http.NewRequest(http.MethodPut, u, r)
	if err != nil {
		return err
	}
	_, err = util.CheckHTTPResponse(c.client.Do(req))
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
		CCCID: path,
		Name:  fi.Name(),
		Size:  fi.Size(),
		Location: []Location{
			{
				Site:             site,
				Path:             filepath.Dir(path),
				TimestampUpdated: fi.ModTime().Unix(),
				User: struct {
					Name string `json:"name"`
				}{
					Name: os.Getenv("USER"),
				},
			},
		},
	}
	return r, nil
}

// TODO replace with proper message validation
func isRecord(b []byte) error {
	var js Record
	return json.Unmarshal(b, &js)
}
