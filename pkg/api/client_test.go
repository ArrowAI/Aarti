package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hclient "aarti/pkg/http/client"
)

var ErrSkip = errors.New("skip")

func TestClientURL(t *testing.T) {
	type test struct {
		name       string
		registry   string
		repository string
		fn         func(ctx context.Context, c *client) error
		url        string
		wantErr    bool
	}
	tests := []test{
		{
			name:       "invalid registry",
			repository: "my-repo",
			wantErr:    true,
		},
		{
			name:       "with repo",
			registry:   "example.org",
			repository: "my-repo",
		},
		{
			name:     "without repo (subpath)",
			registry: "example.org",
			url:      "https://example.org/_packages/apk",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
		{
			name:     "without repo (subdomain)",
			registry: "apk.example.org",
			url:      "https://apk.example.org/_packages",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
		{
			name:     "without repo (subdomain other type)",
			registry: "deb.example.org",
			url:      "https://deb.example.org/_packages/apk",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
		{
			name:       "with repo (subpath)",
			registry:   "example.org",
			repository: "my-repo",
			url:        "https://example.org/_packages/apk/my-repo",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
		{
			name:       "with repo (subdomain)",
			registry:   "apk.example.org",
			repository: "my-repo",
			url:        "https://apk.example.org/_packages/my-repo",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
		{
			name:       "with repo (subdomain other type)",
			registry:   "deb.example.org",
			repository: "my-repo",
			url:        "https://deb.example.org/_packages/apk/my-repo",
			fn: func(ctx context.Context, c *client) error {
				_, err := c.Packages(ctx, "apk")
				return err
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			c, err := NewClient(v.registry, v.repository)
			if v.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if v.fn == nil {
				return
			}
			c.(*client).c = hclient.New(hclient.WithTransport(hclient.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				assert.Equal(t, v.url, r.URL.String())
				return nil, ErrSkip
			})))
			err = v.fn(ctx, c.(*client))
			assert.ErrorIs(t, err, ErrSkip)
		})
	}
}
