package storage

import (
	"context"
)

type storageKey struct{}

func Context(ctx context.Context, r Storage) context.Context {
	return context.WithValue(ctx, storageKey{}, r)
}

func FromContext(ctx context.Context) Storage {
	s, ok := ctx.Value(storageKey{}).(Storage)
	if !ok {
		panic("missing storage in context")
	}
	return s
}
