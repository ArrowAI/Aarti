package registry

import (
	"context"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
)

type BlobStore = registry.BlobStore

type BlobProxy interface {
	registry.ReferenceFetcher
	content.Fetcher
}

type blobs struct {
	BlobStore
	p BlobProxy
}

func (b *blobs) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
	return b.maybeProxy().Fetch(ctx, target)
}

func (b *blobs) FetchReference(ctx context.Context, reference string) (ocispec.Descriptor, io.ReadCloser, error) {
	return b.maybeProxy().FetchReference(ctx, reference)
}

func (b *blobs) maybeProxy() BlobProxy {
	if b.p != nil {
		return b.p
	}
	return b.BlobStore
}
