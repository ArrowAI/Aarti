package react

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

const (
	EndpointEnv = "REACT_ENDPOINT"
)

func NewHandler(dir fs.FS, subpath string) (http.Handler, error) {
	if e := os.Getenv(EndpointEnv); e != "" {
		return newProxy(e)
	}
	return newStatic(dir, subpath)
}

func newStatic(dir fs.FS, subpath string) (http.Handler, error) {
	s, err := fs.Sub(dir, subpath)
	if err != nil {
		return nil, err
	}
	fsrv := http.FileServer(http.FS(s))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Stat(s, strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r.URL.Path = "/"
		}
		fsrv.ServeHTTP(w, r)
	}), nil
}

func newProxy(endpoint string) (http.Handler, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	p := httputil.NewSingleHostReverseProxy(u)
	return p, nil
}

func DevEnv() bool {
	return os.Getenv(EndpointEnv) != ""
}
