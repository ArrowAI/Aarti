package registry

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	auth2 "aarti/pkg/auth"
	"aarti/pkg/cache"
)

const clientID = "lk-artifact-registry"

var clientCache = cache.New()

type Option func(*options)

func makeOptions(host string, opts ...Option) options {
	o := options{
		clientID: clientID,
		host:     host,
		proxy: &options{
			clientID: clientID,
			creds:    &creds{},
		},
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

type options struct {
	host      string
	plainHTTP bool
	insecure  bool
	clientID  string
	clientCA  *x509.CertPool

	debug bool

	// creds are valid only for proxy
	creds *creds

	proxy *options
}

type creds struct {
	user, password string
}

func (o options) apply(ctx context.Context, r *remote.Repository) {
	var u, p string
	if o.creds == nil {
		if a := auth2.FromContext(ctx); a != nil {
			u, p, _ = a.BasicAuth()
		}
	} else {
		u, p = o.creds.user, o.creds.password
	}
	h := sha256.New()
	h.Write([]byte(u))
	h.Write([]byte(p))
	h.Write([]byte(o.host))
	key := fmt.Sprintf("%x", h.Sum(nil))
	if v, ok := clientCache.Get(key); ok {
		clientCache.Set(key, v)
		r.Client = v.(remote.Client)
		r.PlainHTTP = o.plainHTTP
		return
	}
	t := http.RoundTripper(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: o.insecure,
			ClientCAs:          o.clientCA,
		},
	})
	if o.debug {
		t = DebugTransport(t)
	}
	c := &auth.Client{
		ClientID: o.clientID,
		Client: &http.Client{
			Transport: t,
		},
		// expectedHostAddress is of form ipaddr:port
		Credential: auth.StaticCredential(o.host, auth.Credential{
			Username: u,
			Password: p,
		}),
		// Cache caches credentials for accessing the remote registry.
		Cache: auth.NewCache(),
	}
	clientCache.Set(key, c)
	r.Client = c
	r.PlainHTTP = o.plainHTTP
}

func WithClientID(id string) Option {
	return func(o *options) {
		o.clientID = id
	}
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

func WithClientCA(clientCA *x509.CertPool) Option {
	return func(o *options) {
		o.clientCA = clientCA
	}
}

func WithProxy(host string) Option {
	return func(o *options) {
		o.proxy.host = host
	}
}

func WithProxyPlainHTTP() Option {
	return func(o *options) {
		o.proxy.plainHTTP = true
	}
}

func WithProxyInsecure() Option {
	return func(o *options) {
		o.proxy.insecure = true
	}
}

func WithProxyClientCA(clientCA *x509.CertPool) Option {
	return func(o *options) {
		o.proxy.clientCA = clientCA
	}
}

func WithProxyUser(user string) Option {
	return func(o *options) {
		o.proxy.creds.user = user
	}
}

func WithProxyPassword(password string) Option {
	return func(o *options) {
		o.proxy.creds.password = password
	}
}

func WithDebug() Option {
	return func(o *options) {
		o.debug = true
	}
}
