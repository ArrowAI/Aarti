package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"aarti/pkg/utils/logger"
)

type Client interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	Delete(ctx context.Context, url string) (*http.Response, error)

	Options() Options
	Close()
}

func New(opts ...Option) Client {
	o := options{}.apply(opts...)
	if o.errorParser == nil {
		o.errorParser = defaultErrorParser
	}
	if o.transport == nil {
		o.transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: o.insecure,
				ClientCAs:          o.caPool,
			},
		}
	}
	return &client{
		o: o,
		client: &http.Client{
			Transport: o.transport,
		},
	}
}

type ErrorParser func(status string, reader io.Reader) error

var defaultErrorParser = func(status string, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return errors.New("failed to parse error")
	}
	return errors.New(status + ": " + string(b))
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

type client struct {
	client httpClient
	o      options
}

func (c *client) Options() Options {
	return c.o
}

func (c *client) Close() {
	c.client.CloseIdleConnections()
}

func (c *client) Get(ctx context.Context, url string) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, url, nil)
}

func (c *client) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, url, body)
}

func (c *client) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, url, body)
}

func (c *client) Delete(ctx context.Context, url string) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, url, nil)
}

func (c *client) do(ctx context.Context, method string, u string, body io.Reader) (*http.Response, error) {
	if c.o.plainHTTP {
		u = "http://" + u
	} else {
		u = "https://" + u
	}
	msg := fmt.Sprintf("%s %s", method, u)
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	if c.o.host != "" {
		req.Host = c.o.host
	}
	if c.o.user != "" || c.o.pass != "" {
		msg += " (auth)"
		req.SetBasicAuth(c.o.user, c.o.pass)
	}
	if c.o.ua != "" {
		req.Header.Set("User-Agent", c.o.ua)
	}
	logger.C(ctx).Debugf(msg)
	start := time.Now()
	res, err := c.client.Do(req)
	d := time.Since(start)
	if err != nil {
		return nil, err
	}
	logger.C(ctx).Debugf("%s %s", res.Status, d)
	if res.StatusCode >= 400 {
		return res, c.o.errorParser(res.Status, res.Body)
	}
	return res, nil
}
