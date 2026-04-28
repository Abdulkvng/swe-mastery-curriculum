// pool/concurrency_test.go
//
// Stress tests: many goroutines hammering the pool simultaneously.
// Run with `go test -race ./pool` to catch data races.

package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStressManyGoroutines(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	const goroutines = 50
	const iterations = 100

	p, _ := makeTestPool(t, func(c *Config) {
		c.MaxConns = 8
		c.AcquireTimeout = 2 * time.Second
	})

	var wg sync.WaitGroup
	var done atomic.Uint64
	var failures atomic.Uint64

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				c, err := p.Acquire(context.Background())
				if err != nil {
					failures.Add(1)
					continue
				}
				// Simulate some work.
				time.Sleep(time.Microsecond * 100)
				_ = c.Release()
				done.Add(1)
			}
		}()
	}

	wg.Wait()

	if failures.Load() > 0 {
		t.Errorf("had %d acquire failures (some are OK under contention, but suspicious)", failures.Load())
	}

	stats := p.Stats()
	if stats.InUse != 0 {
		t.Errorf("after work done, InUse = %d, want 0", stats.InUse)
	}
	if stats.AcquireCount < uint64(goroutines*iterations-int(failures.Load())) {
		t.Errorf("AcquireCount = %d, expected ~ %d",
			stats.AcquireCount, goroutines*iterations)
	}

	t.Logf("stats: %+v", stats)
}
