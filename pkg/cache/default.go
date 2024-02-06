package cache

import (
	"time"
)

var (
	// DefaultTTL is the default TTL for cache items.
	// It matches the default TTL of docker.io registry.
	DefaultTTL = 6 * time.Hour
	d          = New()
)

func Set(key string, value any, opts ...Option) {
	d.Set(key, value, opts...)
}
func Get(key string) (any, bool) {
	return d.Get(key)
}
