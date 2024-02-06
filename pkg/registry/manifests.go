package registry

import (
	"context"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
)

type ManifestStore = registry.ManifestStore

type ManifestProxy interface {
	registry.ReferenceFetcher
	content.Fetcher
}

type manifests struct {
	ManifestStore
	p ManifestProxy
}

func (m *manifests) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
	return m.maybeProxy().Fetch(ctx, target)
}

func (m *manifests) FetchReference(ctx context.Context, reference string) (ocispec.Descriptor, io.ReadCloser, error) {
	return m.maybeProxy().FetchReference(ctx, reference)
}

func (m *manifests) maybeProxy() ManifestProxy {
	if m.p != nil {
		return m.p
	}
	return m.ManifestStore
}
