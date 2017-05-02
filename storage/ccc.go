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
	site        string
	dtsURL      string
}

// NewCCCBackend returns a CCCBackend instance, configured to limit
// file system access to the given allowed directories.
func NewCCCBackend(conf config.CCCStorage) (*CCCBackend, error) {
	b := &CCCBackend{
		allowedDirs: conf.AllowedDirs,
		site:        conf.Site,
		dtsURL:      conf.DTSUrl,
	}
	return b, nil
}

// Get copies a file from storage into the given hostPath.
func (ccc *CCCBackend) Get(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting download", "url", url, "hostPath", hostPath)

	path := strings.TrimPrefix(url, CCCProtocol)
	path, rerr := resolveCCCID(path, ccc.site, ccc.dtsURL)
	log.Info("Resolved DTS url", "url", url, "sitePath", path)
	if rerr != nil {
		return rerr
	}
	if !isAllowed(path, ccc.allowedDirs) {
		return fmt.Errorf("Can't access file, path is not in allowed directories:  %s", path)
	}

	var err error
	if class == File {
		err = linkFile(path, hostPath)
	} else if class == Directory {
		err = copyDir(path, hostPath)
	} else {
		err = fmt.Errorf("Unknown file class: %s", class)
	}

	if err == nil {
		log.Info("Finished download", "url", url, "hostPath", hostPath)
	}
	return err
}

// Put copies a file from the hostPath into storage.
func (ccc *CCCBackend) Put(ctx context.Context, url string, hostPath string, class tes.FileType) error {
	log.Info("Starting upload", "url", url, "hostPath", hostPath)

	id := strings.TrimPrefix(url, CCCProtocol)
	path, rerr := resolveCCCID(id, ccc.site, ccc.dtsURL)
	if rerr == nil {
		return fmt.Errorf("CCCID %s conflicts with an existing record: %v", id, path)
	}
	if !isAllowed(path, ccc.allowedDirs) {
		return fmt.Errorf("Can't access file, path is not in allowed directories:  %s", url)
	}

	var err error
	if class == File {
		err = copyFile(hostPath, path)
	} else if class == Directory {
		err = copyDir(hostPath, path)
	} else {
		err = fmt.Errorf("Unknown file class: %s", class)
	}
	if err != nil {
		return fmt.Errorf("Failed to upload %s: %v", url, err)
	}

	r, cerr := createDTSRecord(path, ccc.site, ccc.dtsURL)
	if cerr != nil {
		return fmt.Errorf("Failed to create DTS Record for output %s. %v", url, cerr)
	}

	log.Debug("Created DTS Record.", "record", r)
	log.Info("Finished upload", "url", url, "hostPath", hostPath)

	return nil
}

// Supports indicates whether this backend supports the given storage request.
// For the CCCBackend, the url must start with "ccc://"
func (ccc *CCCBackend) Supports(url string, hostPath string, class tes.FileType) bool {
	return strings.HasPrefix(url, CCCProtocol)
}

func resolveCCCID(path string, site string, dtsURL string) (string, error) {
	cli, err := dts.NewClient(dtsURL)
	if err != nil {
		return "", err
	}
	entry, err := cli.GetFile(path)
	if err != nil {
		return "", err
	}
	log.Debug("DTS Record", "record", entry)
	for _, location := range entry.Location {
		if site == location.Site {
			return filepath.Join(location.Path, entry.Name), nil
		}
	}
	err = fmt.Errorf("No DTS record for %s on site %s", path, site)
	return "", err
}

func createDTSRecord(path string, site string, dtsURL string) (*dts.Record, error) {
	r, err := dts.GenerateRecord(path, site)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate DTS Record for output %s. %v", path, err)
	}
	cli, err := dts.NewClient(dtsURL)
	if err != nil {
		return nil, err
	}
	msg, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	err = cli.PostFile(msg)
	if err != nil {
		return nil, err
	}
	return r, err
}
