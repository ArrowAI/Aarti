package storage

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/opencontainers/go-digest"
)

var _ Artifact = (*File)(nil)

type File struct {
	name string
	data []byte
	r    io.Reader
}

func NewFile(name string, data []byte) *File {
	return &File{
		name: name,
		data: data,
		r:    bytes.NewReader(data),
	}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}

func (f *File) Name() string {
	return filepath.Base(f.name)
}

func (f *File) Arch() string {
	return ""
}

func (f *File) Version() string {
	return ""
}

func (f *File) Path() string {
	return f.name
}

func (f *File) Size() int64 {
	return int64(len(f.data))
}

func (f *File) Digest() digest.Digest {
	return digest.FromBytes(f.data)
}

func (f *File) Close() error {
	return nil
}
