package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/ccc/dts"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"path/filepath"
	"strings"
)

// CCCProtocol defines the expected prefix of URL matching this storage system.
// e.g. "file:///path/to/file" matches the CCC storage system.
const CCCProtocol = "ccc://"

// CCCBackend provides access to a ccc-disk storage system.
type CCCBackend struct {
	allowedDirs []string
	localSite   string
	remoteSites []string
	outputSite  string
	dtsURL      string
}

// NewCCCBackend returns a CCCBackend instance, configured to limit
// file system access to the given allowed directories.
func NewCCCBackend(conf config.CCCStorage) (*CCCBackend, error) {
	b := &CCCBackend{
		allowedDirs: conf.AllowedDirs,
		dtsURL:      conf.DTSUrl,
		localSite:   conf.SiteMap.Local,
	}

	local := conf.SiteMap.Local
	central := conf.SiteMap.Central
	switch conf.Strategy {
	case "fetch_file":
		b.remoteSites = []string{central}
		b.outputSite = local
	case "push_file":
		b.remoteSites = nil
		b.outputSite = central
	default: // AKA "routed_file"
		b.remoteSites = nil
		b.outputSite = local
	}

	return b, nil
}

// Get copies a file from storage into the given hostPath.
func (ccc *CCCBackend) Get(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting download", "url", url, "hostPath", hostPath)

	path := strings.TrimPrefix(url, CCCProtocol)
	path, remote, rerr := ccc.resolveCCCID(path)
	if rerr != nil {
		return rerr
	}
	log.Info("Resolved DTS url", "url", url, "sitePath", path)
	if !isAllowed(path, ccc.allowedDirs) {
		return fmt.Errorf("Can't access file, path is not in allowed directories:  %s", path)
	}

	var err error
	if remote {
		if class == File {
			cli := NewSCPClient(ccc.outputSite)
			err = cli.SCPRemoteToLocal(path, hostPath)
		} else if class == Directory {
			err = fmt.Errorf("SCP of directories not supported")
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
	} else {
		if class == File {
			err = linkFile(path, hostPath)
		} else if class == Directory {
			err = copyDir(path, hostPath)
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
	}

	if err == nil {
		log.Info("Finished download", "url", url, "hostPath", hostPath)
	}
	return err
}

// Put copies a file from the hostPath into storage.
func (ccc *CCCBackend) Put(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting upload", "url", url, "hostPath", hostPath)

	path := strings.TrimPrefix(url, CCCProtocol)
	record, remote, rerr := ccc.resolveCCCID(path)
	if rerr == nil {
		return fmt.Errorf("CCCID %s conflicts with an existing record: %+v", path, record)
	}
	if !isAllowed(path, ccc.allowedDirs) {
		return fmt.Errorf("Can't access file, path is not in allowed directories:  %s", url)
	}

	var err error
	if remote {
		if class == File {
			cli := NewSCPClient(ccc.outputSite)
			err = cli.SCPLocalToRemote(hostPath, path)
		} else if class == Directory {
			err = fmt.Errorf("SCP of directories not supported")
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
	} else {
		if class == File {
			err = copyFile(hostPath, path)
		} else if class == Directory {
			err = copyDir(hostPath, path)
		} else {
			err = fmt.Errorf("Unknown file class: %s", class)
		}
	}
	if err != nil {
		return fmt.Errorf("Failed to upload %s: %v", url, err)
	}

	r, cerr := ccc.createDTSRecord(path)
	if cerr != nil {
		return fmt.Errorf("Failed to create DTS Record for output %s. %v", url, cerr)
	}

	log.Debug("Created DTS Record", "record", fmt.Sprintf("%+v", r))
	log.Info("Finished upload", "url", url, "hostPath", hostPath)

	return nil
}

// Supports indicates whether this backend supports the given storage request.
// For the CCCBackend, the url must start with "ccc://"
func (ccc *CCCBackend) Supports(url string, hostPath string, class tes.FileType) bool {
	return strings.HasPrefix(url, CCCProtocol)
}

func (ccc *CCCBackend) resolveCCCID(path string) (string, bool, error) {
	var remote bool
	cli, err := dts.NewClient(ccc.dtsURL)
	if err != nil {
		return "", remote, err
	}
	entry, err := cli.GetFile(path)
	if err != nil {
		return "", remote, err
	}
	log.Debug("DTS Record", "record", fmt.Sprintf("%+v", entry))
	for _, location := range entry.Location {
		if ccc.localSite == location.Site {
			return filepath.Join(location.Path, entry.Name), remote, nil
		}
		for _, r := range ccc.remoteSites {
			if r == location.Site {
				remote = true
				return filepath.Join(location.Path, entry.Name), remote, nil
			}
		}
	}
	return "", remote, fmt.Errorf("%s is not located at an accessible site", path)
}

func (ccc *CCCBackend) createDTSRecord(path string) (*dts.Record, error) {
	r, err := dts.GenerateRecord(path, ccc.localSite)
	if ccc.outputSite != ccc.localSite {
		var l dts.Location
		l = r.Location[0]
		l.Site = ccc.outputSite
		r.Location = append(r.Location, l)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to generate DTS Record for output %s. %v", path, err)
	}
	cli, err := dts.NewClient(ccc.dtsURL)
	if err != nil {
		return nil, err
	}
	msg, err := json.Marshal(r)
	log.Debug("Created DTS message", "message", string(msg))
	if err != nil {
		return nil, err
	}
	err = cli.PostFile(msg)
	if err != nil {
		return nil, err
	}
	return r, err
}
