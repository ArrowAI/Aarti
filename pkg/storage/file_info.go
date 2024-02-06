package storage

import (
	"path/filepath"

	"github.com/opencontainers/go-digest"
)

type info struct {
	version string
	path    string
	size    int64
	digest  digest.Digest
	meta    []byte
}

func (i *info) Name() string {
	return filepath.Base(i.path)
}

func (i *info) Arch() string {
	return ""
}

func (i *info) Version() string {
	return i.version
}

func (i *info) Path() string {
	return i.path
}

func (i *info) Size() int64 {
	return i.size
}

func (i *info) Digest() digest.Digest {
	return i.digest
}

func (i *info) Meta() []byte {
	return i.meta
}
