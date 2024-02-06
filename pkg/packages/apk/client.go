package apk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	hclient "aarti/pkg/http/client"
	"aarti/pkg/packages"
)

var _ packages.Client = (*client)(nil)

func NewClient(registry, repository, branch, repo string, opts ...hclient.Option) (packages.Client, error) {
	if registry == "" {
		return nil, fmt.Errorf("registry, is required")
	}
	var base string
	if strings.HasPrefix(registry, Name+".") {
		base = fmt.Sprintf("%s/%s", registry, repository)
	} else {
		base = fmt.Sprintf("%s/%s/%s", registry, Name, repository)
	}
	return &client{
		c:          hclient.New(opts...),
		registry:   registry,
		base:       strings.TrimSuffix(base, "/"),
		repository: repository,
		branch:     branch,
		repo:       repo,
	}, nil
}

type client struct {
	c          hclient.Client
	base       string
	registry   string
	repository string
	branch     string
	repo       string
}

func (c *client) Key(ctx context.Context) (string, error) {
	if c.branch == "" || c.repo == "" {
		return "", fmt.Errorf("branch and repo, are required")
	}
	res, err := c.c.Get(ctx, c.path(c.branch, c.repo, "key"))
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
	if c.branch == "" || c.repo == "" {
		return "", fmt.Errorf("branch and repo, are required")
	}
	res, err := c.c.Get(ctx, c.path(c.branch, c.repo, "setup"))
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
	if c.branch == "" || c.repo == "" {
		return fmt.Errorf("branch and repo, are required")
	}
	_, err := c.c.Put(ctx, c.path(c.branch, c.repo, "push"), r)
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
