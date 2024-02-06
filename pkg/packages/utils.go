package packages

import (
	"errors"
	"net/http"
)

var ErrAlreadyConfigured = errors.New("repository already configured")

func Scheme(r *http.Request) string {
	if p := r.Header.Get("X-Forwarded-Proto"); p != "" {
		return p
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
