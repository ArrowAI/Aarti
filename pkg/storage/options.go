package storage

import (
	"context"

	"aarti/pkg/registry"
)

type optionsKey struct{}

func WithOptions(ctx context.Context, opts ...Option) context.Context {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return context.WithValue(ctx, optionsKey{}, o)
}

func Options(ctx context.Context) options {
	o, _ := ctx.Value(optionsKey{}).(options)
	return o
}

type options struct {
	host         string
	key          []byte
	repo         string
	artifactTags bool
	ropts        []registry.Option
}

func (o options) Host() string {
	return o.host
}

func (o options) Repo() string {
	return o.repo
}

func (o options) Key() []byte {
	return o.key
}

func (o options) NewRegistry(ctx context.Context) (registry.Registry, error) {
	return registry.NewRegistry(ctx, o.host, o.ropts...)
}

func (o options) NewRepository(ctx context.Context, name string) (registry.Repository, error) {
	return registry.NewRepository(ctx, name, o.ropts...)
}

type Option func(o *options)

func WithHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

func WithKey(key []byte) Option {
	return func(o *options) {
		o.key = key
	}
}

func WithRepo(repo string) Option {
	return func(o *options) {
		o.repo = repo
	}
}

func WithArtifactTags() Option {
	return func(o *options) {
		o.artifactTags = true
	}
}

func WithRegistryOptions(opts ...registry.Option) Option {
	return func(o *options) {
		o.ropts = opts
	}
}
