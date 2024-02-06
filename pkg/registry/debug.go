package registry

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func DebugTransport(t http.RoundTripper) http.RoundTripper {
	if t == nil {
		t = http.DefaultTransport
	}
	return &debugTransport{t: t}
}

type debugTransport struct {
	t http.RoundTripper
}

func (s *debugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", bytes)
	resp, err := s.t.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	bytes, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", bytes)

	return resp, err
}
