package webdash

import (
	//"github.com/elazarl/go-bindata-assetfs"
	"net/http"
)

// FileServer provides access to the bundled web assets (HTML, CSS, etc)
// via an http.Handler
func FileServer() http.Handler {
	return http.FileServer(http.Dir("./react-dash/funnel/build"))
}

/*
func FileServer() http.Handler {
	fs := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
	}
	return http.FileServer(fs)
}
*/

// Handler handles static webdash files
func Handler() *http.ServeMux {
	// Static files are bundled into webdash
	fs := FileServer()
	// Set up URL path handlers
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	return mux
}
