package storage

import (
	"context"
	"runtime"
	"sync"
	"time"

	"aarti/pkg/utils/logger"

	"github.com/dustin/go-humanize"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
)

func copts(name string) oras.CopyOptions {
	var times sync.Map
	return oras.CopyOptions{
		CopyGraphOptions: oras.CopyGraphOptions{
			Concurrency: runtime.NumCPU(),
			PreCopy: func(ctx context.Context, desc ocispec.Descriptor) error {
				times.Store(desc.Digest.String(), time.Now())
				logger.C(ctx).WithFields(
					"digest", desc.Digest.String(),
					"size", humanize.Bytes(uint64(desc.Size)),
					"ref", name,
				).Infof("uploading")
				return nil
			},
			OnCopySkipped: func(ctx context.Context, desc ocispec.Descriptor) error {
				logger.C(ctx).WithFields(
					"digest", desc.Digest.String(),
					"size", humanize.Bytes(uint64(desc.Size)),
					"ref", name,
				).Infof("already exists")
				return nil
			},
			PostCopy: func(ctx context.Context, desc ocispec.Descriptor) error {
				var dur time.Duration
				if v, ok := times.Load(desc.Digest.String()); ok {
					dur = time.Since(v.(time.Time))
				}
				logger.C(ctx).WithFields(
					"digest", desc.Digest.String(),
					"size", humanize.Bytes(uint64(desc.Size)),
					"ref", name,
					"duration", dur,
				).Infof("uploaded")
				return nil
			},
		},
	}
}
