package registry

import (
	"context"

	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

type Registry = registry.Registry

func NewRegistry(ctx context.Context, name string, opts ...Option) (Registry, error) {
	r, err := remote.NewRegistry(name)
	if err != nil {
		return nil, err
	}
	o := makeOptions(r.RepositoryOptions.Reference.Host(), opts...)
	o.apply(ctx, (*remote.Repository)(&r.RepositoryOptions))
	var proxy Registry
	if o.proxy.host != "" {
		p, err := remote.NewRegistry(o.proxy.host)
		if err != nil {
			return nil, err
		}
		o.proxy.apply(ctx, (*remote.Repository)(&p.RepositoryOptions))
		proxy = p
	}
	return &reg{r: r, p: proxy}, nil
}

type reg struct {
	r Registry
	p Registry
}

func (r *reg) Repositories(ctx context.Context, last string, fn func(repos []string) error) error {
	return r.r.Repositories(ctx, last, fn)
}

func (r *reg) Repository(ctx context.Context, name string) (Repository, error) {
	rep, err := r.r.Repository(ctx, name)
	if err != nil {
		return nil, err
	}
	var proxy Repository
	if r.p != nil {
		proxy, err = r.p.Repository(ctx, name)
		if err != nil {
			return nil, err
		}
	}
	return &repo{Repository: rep, p: proxy}, nil
}
