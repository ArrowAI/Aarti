package mutex

import (
	"context"
	"sync"

	"aarti/pkg/utils/logger"
)

// New returns a properly initialized local
func New() Mutex {
	return &local{
		store: make(map[string]*sync.RWMutex),
	}
}

// Mutex is a simple key/value store for arbitrary mutexes. It can be used to
// serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on.
//
// The initial use case is to let aws_security_group_rule resources serialize
// their access to individual security groups based on SG ID.
type Mutex interface {
	// Lock the mutex for the given key. Caller is responsible for calling Unlock
	// for the same key
	Lock(ctx context.Context, key string)
	// Unlock the mutex for the given key. Caller must have called Lock for the same key first
	Unlock(ctx context.Context, key string)
	// RLock the mutex for the given key. Caller is responsible for calling RUnlock
	RLock(ctx context.Context, key string)
	// RUnlock the mutex for the given key. Caller must have called RLock for the same key first
	RUnlock(ctx context.Context, key string)
}

type local struct {
	lock  sync.Mutex
	store map[string]*sync.RWMutex
}

// Lock the mutex for the given key. Caller is responsible for calling Unlock
// for the same key
func (m *local) Lock(ctx context.Context, key string) {
	logger.C(ctx).Debugf("Locking %q", key)
	m.get(ctx, key).Lock()
	logger.C(ctx).Debugf("Locked %q", key)
}

// Unlock the mutex for the given key. Caller must have called Lock for the same key first
func (m *local) Unlock(ctx context.Context, key string) {
	logger.C(ctx).Debugf("Unlocking %q", key)
	m.get(ctx, key).Unlock()
	logger.C(ctx).Debugf("Unlocked %q", key)
}

// RLock the mutex for the given key. Caller is responsible for calling RUnlock
// for the same key
func (m *local) RLock(ctx context.Context, key string) {
	logger.C(ctx).Debugf("RLocking %q", key)
	m.get(ctx, key).RLock()
	logger.C(ctx).Debugf("RLocked %q", key)
}

// RUnlock the mutex for the given key. Caller must have called RLock for the same key first
func (m *local) RUnlock(ctx context.Context, key string) {
	logger.C(ctx).Debugf("RUnlocking %q", key)
	m.get(ctx, key).RUnlock()
	logger.C(ctx).Debugf("RUnlocked %q", key)
}

// Returns a mutex for the given key, no guarantee of its lock status
func (m *local) get(_ context.Context, key string) *sync.RWMutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.RWMutex{}
		m.store[key] = mutex
	}
	return mutex
}
