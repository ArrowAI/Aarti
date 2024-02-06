package packages

import (
	"context"
	"io"
)

type Client interface {
	Signer
	Setuper
	Puller
	Pusher
	Deleter
}

type Puller interface {
	Pull(ctx context.Context, path string) (io.ReadCloser, int64, error)
}

type Pusher interface {
	Push(ctx context.Context, r io.Reader) error
}

type Deleter interface {
	Delete(ctx context.Context, path string) error
}

type Signer interface {
	Key(ctx context.Context) (string, error)
}

type Setuper interface {
	SetupScript(ctx context.Context) (string, error)
	SetupLocal(ctx context.Context, force bool) error
}
