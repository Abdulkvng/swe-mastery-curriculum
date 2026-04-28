// pool/pool_test.go
//
// Tests use a fake Connector that doesn't need a real database.
// We open an actual *sql.DB but with a no-op driver via "sqlite3 :memory:"
// equivalent — easier: we use a fake that returns a *sql.DB pointing at /dev/null.
//
// For true isolation we'd build a fake `driver.Driver`. Here we use the
// stdlib's sql package with a registered "fakedriver" that lets us avoid
// real I/O. To keep this file dependency-free we use database/sql's
// Open("sqlite3", ":memory:") IF you have github.com/mattn/go-sqlite3 — but
// to keep the curriculum tidy, we mock with a simpler approach:
// we override Connector to return a *sql.DB that's never actually queried,
// and override HealthCheck to return nil.

package pool

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeConnector returns a fresh *sql.DB-shaped value each time.
// We can't easily fabricate a *sql.DB, so we open one against a driver
// that errors on real queries. Tests should NOT actually call db.Query.
//
// To keep this file zero-dependency, we do something subtle: we use
// (*sql.DB)(nil) as a marker. Our HealthCheck doesn't deref it.
// In real life you'd register a fake driver via sql.Register.

type counterConnector struct {
	creates atomic.Uint64
	failAt  uint64 // if > 0, fail the Nth create
	fail    error
}

func (c *counterConnector) connect(ctx context.Context) (*sql.DB, error) {
	n := c.creates.Add(1)
	if c.failAt > 0 && n == c.failAt {
		return nil, c.fail
	}
	// Return a pointer that's never deref'd. Tests don't call .DB().
	return new(sql.DB), nil
}

func noopHealth(ctx context.Context, _ *sql.DB) error { return nil }

func makeTestPool(t *testing.T, modify func(*Config)) (Pool, *counterConnector) {
	t.Helper()
	cc := &counterConnector{}
	cfg := DefaultConfig()
	cfg.Connector = cc.connect
	cfg.HealthCheck = noopHealth
	cfg.MinConns = 0
	cfg.MaxConns = 5
	cfg.AcquireTimeout = 200 * time.Millisecond
	cfg.HealthInterval = 0 // disable reaper for tests
	if modify != nil {
		modify(&cfg)
	}
	p, err := NewPool(cfg)
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })
	return p, cc
}

func TestAcquireRelease(t *testing.T) {
	p, _ := makeTestPool(t, nil)
	ctx := context.Background()

	c1, err := p.Acquire(ctx)
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	if got := p.Stats().InUse; got != 1 {
		t.Errorf("InUse = %d, want 1", got)
	}

	if err := c1.Release(); err != nil {
		t.Errorf("Release: %v", err)
	}
	if got := p.Stats().InUse; got != 0 {
		t.Errorf("InUse after release = %d, want 0", got)
	}
	if got := p.Stats().Idle; got != 1 {
		t.Errorf("Idle after release = %d, want 1", got)
	}
}

func TestReleaseIdempotent(t *testing.T) {
	p, _ := makeTestPool(t, nil)
	ctx := context.Background()

	c, _ := p.Acquire(ctx)
	if err := c.Release(); err != nil {
		t.Fatal(err)
	}
	// Second release should be a no-op (no panic, no error).
	if err := c.Release(); err != nil {
		t.Errorf("second Release returned %v, want nil", err)
	}
}

func TestMaxConnsRespected(t *testing.T) {
	p, cc := makeTestPool(t, func(c *Config) { c.MaxConns = 3 })
	ctx := context.Background()

	conns := make([]*Conn, 3)
	for i := range conns {
		c, err := p.Acquire(ctx)
		if err != nil {
			t.Fatalf("Acquire %d: %v", i, err)
		}
		conns[i] = c
	}

	// 4th acquire should time out.
	_, err := p.Acquire(ctx)
	if err == nil {
		t.Fatal("expected timeout, got nil")
	}

	if got := cc.creates.Load(); got != 3 {
		t.Errorf("creates = %d, want 3", got)
	}

	for _, c := range conns {
		_ = c.Release()
	}
}

func TestAcquireWaitsForRelease(t *testing.T) {
	p, _ := makeTestPool(t, func(c *Config) {
		c.MaxConns = 1
		c.AcquireTimeout = 2 * time.Second
	})
	ctx := context.Background()

	c1, err := p.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var acquireErr error
	var got *Conn
	start := time.Now()
	go func() {
		defer wg.Done()
		got, acquireErr = p.Acquire(ctx)
	}()

	// Release after a short delay — the waiter should pick up.
	time.Sleep(100 * time.Millisecond)
	_ = c1.Release()

	wg.Wait()
	if acquireErr != nil {
		t.Fatalf("waiter got error: %v", acquireErr)
	}
	if got == nil {
		t.Fatal("waiter got nil conn")
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("waiter blocked too long: %v", elapsed)
	}
	if p.Stats().WaitCount == 0 {
		t.Error("WaitCount should have incremented")
	}
	_ = got.Release()
}

func TestContextCancel(t *testing.T) {
	p, _ := makeTestPool(t, func(c *Config) { c.MaxConns = 1 })

	c1, _ := p.Acquire(context.Background())
	defer c1.Release()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := p.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Either context.Canceled or ErrAcquireTimeout is acceptable, depending
	// on which fires first.
}

func TestCloseRejectsAcquire(t *testing.T) {
	p, _ := makeTestPool(t, nil)
	_ = p.Close()
	_, err := p.Acquire(context.Background())
	if err != ErrPoolClosed {
		t.Errorf("got %v, want ErrPoolClosed", err)
	}
}

func TestObserversFire(t *testing.T) {
	var acquireFires, releaseFires, createFires int
	var omu sync.Mutex
	p, _ := makeTestPool(t, func(c *Config) {
		c.Observers = Observers{
			OnAcquire: func(_ *Conn) { omu.Lock(); acquireFires++; omu.Unlock() },
			OnRelease: func(_ *Conn) { omu.Lock(); releaseFires++; omu.Unlock() },
			OnCreate:  func(_ *Conn) { omu.Lock(); createFires++; omu.Unlock() },
		}
	})

	c, _ := p.Acquire(context.Background())
	_ = c.Release()

	if createFires != 1 || acquireFires != 1 || releaseFires != 1 {
		t.Errorf("creates=%d acquires=%d releases=%d", createFires, acquireFires, releaseFires)
	}
}

func TestSingletonOnce(t *testing.T) {
	// Reset (only safe in tests because the package's vars persist between calls).
	defaultPool = nil
	defaultErr = nil
	defaultOnce = sync.Once{}

	cc := &counterConnector{}
	SetDefaultConfig(Config{
		Connector:   cc.connect,
		MinConns:    0,
		MaxConns:    1,
		HealthCheck: noopHealth,
	})

	p1, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	p2, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	if p1 != p2 {
		t.Error("Default should return the same instance")
	}
	_ = p1.Close()
}
