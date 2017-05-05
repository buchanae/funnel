package ccc

import (
  "fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"google.golang.org/grpc"
	"net/url"
  "encoding/json"
  "encoding/base64"
)

type giddata struct {
  Site string
  LocalID string
}

type SiteMapper interface {
	GlobalID(site, lid string) string
	LocalID(gid string) (string, error)
	Site(gid string) (string, error)
	Sites() []string
	Client(site string) (tes.TaskServiceClient, error)
}

type siteMapper struct {
	conf config.Config
	// Overrideable for testing
	getClient func(address string) (tes.TaskServiceClient, error)
}

func (s *siteMapper) Sites() []string {
  return s.conf.CCC.Sites
}

func (s *siteMapper) Site(gid string) (string, error) {
	site, _, err := parse(gid)
	return site, err
}

func (s *siteMapper) LocalID(gid string) (string, error) {
	_, lid, err := parse(gid)
	return lid, err
}

func (s *siteMapper) GlobalID(site, lid string) string {
  gid := giddata{site, lid}
  js, _ := json.Marshal(gid)
  return base64.StdEncoding.EncodeToString([]byte(js))
}

func (s *siteMapper) Client(site string) (tes.TaskServiceClient, error) {
	u := normalize(site, "")
  if u.Hostname() == "" {
    return nil, fmt.Errorf("No site hostname")
  }

	address := u.Hostname() + ":" + s.conf.RPCPort
	if s.getClient != nil {
		return s.getClient(address)
	}
	return getClient(address)
}

func normalize(site, path string) *url.URL {
	u, err := url.Parse(site)
	if err != nil || u.Host == "" {
		return &url.URL{Scheme: "http", Host: site, Path: path}
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if path != "" {
		u.Path = path
	}
	return u
}

func getClient(address string) (tes.TaskServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	return tes.NewTaskServiceClient(conn), err
}

func parse(raw string) (string, string, error) {
  js, _ := base64.StdEncoding.DecodeString(raw)
  gid := giddata{}
  json.Unmarshal(js, &gid)
  return gid.Site, gid.LocalID, nil
}
