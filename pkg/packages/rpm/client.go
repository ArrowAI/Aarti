package rpm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	hclient "aarti/pkg/http/client"
	"aarti/pkg/packages"
)

var _ Client = (*client)(nil)

type Client interface {
	packages.Client
	Repo(ctx context.Context) (string, error)
}

func NewClient(registry, repository string, opts ...hclient.Option) (Client, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry is required")
	}
	var base string
	if strings.HasPrefix(registry, Name+".") {
		base = fmt.Sprintf("%s/%s", registry, repository)
	} else {
		base = fmt.Sprintf("%s/%s/%s", registry, Name, repository)
	}
	return &client{
		c:    hclient.New(opts...),
		base: strings.TrimSuffix(base, "/"),
		repo: repository,
	}, nil
}

type client struct {
	c    hclient.Client
	base string
	repo string
}

func (c *client) Key(ctx context.Context) (string, error) {
	res, err := c.c.Get(ctx, c.path(RepositoryPublicKey))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *client) Repo(ctx context.Context) (string, error) {
	u := fmt.Sprintf("%s.repo", c.base)
	if c.repo == "" {
		u = fmt.Sprintf("%s/.repo", c.base)
	}
	res, err := c.c.Get(ctx, u)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *client) SetupScript(ctx context.Context) (string, error) {
	res, err := c.c.Get(ctx, c.path("setup"))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *client) Push(ctx context.Context, r io.Reader) error {
	_, err := c.c.Put(ctx, c.path("push"), r)
	return err
}

func (c *client) Pull(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	res, err := c.c.Get(ctx, c.path(path))
	if err != nil {
		return nil, 0, err
	}
	return res.Body, res.ContentLength, nil
}

func (c *client) Delete(ctx context.Context, path string) error {
	_, err := c.c.Delete(ctx, c.path(path))
	return err
}

func (c *client) path(parts ...string) string {
	return fmt.Sprintf("%s/%s", c.base, strings.Join(parts, "/"))
}
