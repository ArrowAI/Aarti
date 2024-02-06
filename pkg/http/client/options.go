package client

import (
	"crypto/x509"
	"net/http"
)

type Option func(*options)

type Options interface {
	Scheme() string
	Host() string
	PlainHTTP() bool
	Insecure() bool
	CA() *x509.CertPool
	BasicAuth() (username, password string, ok bool)
}

type options struct {
	plainHTTP bool
	insecure  bool
	caPool    *x509.CertPool

	user, pass  string
	ua          string
	host        string
	errorParser ErrorParser

	transport http.RoundTripper
}

func (o options) Scheme() string {
	if o.plainHTTP {
		return "http"
	}
	return "https"
}

func (o options) Host() string {
	return o.host
}

func (o options) PlainHTTP() bool {
	return o.plainHTTP
}

func (o options) Insecure() bool {
	return o.insecure
}

func (o options) CA() *x509.CertPool {
	return o.caPool
}

func (o options) BasicAuth() (username, password string, ok bool) {
	return o.user, o.pass, o.user != "" || o.pass != ""
}

func (o options) apply(opts ...Option) options {
	for _, v := range opts {
		v(&o)
	}
	return o
}

func WithPlainHTTP() Option {
	return func(o *options) {
		o.plainHTTP = true
	}
}

func WithInsecure() Option {
	return func(o *options) {
		o.insecure = true
	}
}

func WithClientCA(ca *x509.CertPool) Option {
	return func(o *options) {
		o.caPool = ca
	}
}

func WithBasicAuth(user, pass string) Option {
	return func(o *options) {
		o.user = user
		o.pass = pass
	}
}

func WithErrorParser(p ErrorParser) Option {
	return func(o *options) {
		o.errorParser = p
	}
}

func WithUserAgent(ua string) Option {
	return func(o *options) {
		o.ua = ua
	}
}

func WithHost(h string) Option {
	return func(o *options) {
		o.host = h
	}
}

func WithTransport(t http.RoundTripper) Option {
	return func(o *options) {
		o.transport = t
	}
}
