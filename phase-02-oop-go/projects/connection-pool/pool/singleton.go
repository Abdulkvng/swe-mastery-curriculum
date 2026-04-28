// pool/singleton.go
//
// Optional Default() pool. Demonstrates the Singleton pattern using sync.Once.
//
// USAGE:
//   pool.SetDefaultConfig(cfg)
//   p := pool.Default()
//   conn, _ := p.Acquire(ctx)
//
// You typically would NOT use this in a real codebase — pass the pool
// explicitly to functions instead. This is here to demonstrate the pattern.

package pool

import (
	"errors"
	"sync"
)

var (
	defaultPool   Pool
	defaultErr    error
	defaultOnce   sync.Once
	defaultConfig Config
	defaultMu     sync.Mutex // protects defaultConfig before Default() is called
)

// SetDefaultConfig sets the config used the first time Default() is called.
// Calling this AFTER Default() has been invoked is a no-op (Singleton — set in stone).
func SetDefaultConfig(cfg Config) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultConfig = cfg
}

// Default returns the process-wide default pool. Created on first call,
// returned thereafter. Safe to call from many goroutines simultaneously.
func Default() (Pool, error) {
	defaultOnce.Do(func() {
		defaultMu.Lock()
		cfg := defaultConfig
		defaultMu.Unlock()

		if cfg.Connector == nil {
			defaultErr = errors.New("pool: SetDefaultConfig was never called")
			return
		}
		defaultPool, defaultErr = NewPool(cfg)
	})
	return defaultPool, defaultErr
}

// MustDefault is like Default but panics on error. Useful in main().
func MustDefault() Pool {
	p, err := Default()
	if err != nil {
		panic(err)
	}
	return p
}
