package auth

import (
	"context"
)

type Basic interface {
	BasicAuth() (username, password string, ok bool)
}

type key struct{}

func Context(ctx context.Context, a Basic) context.Context {
	return context.WithValue(ctx, key{}, a)
}

func FromContext(ctx context.Context) Basic {
	a, _ := ctx.Value(key{}).(Basic)
	return a
}
