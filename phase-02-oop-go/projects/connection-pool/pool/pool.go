// pool/pool.go
//
// The heart of the connection pool.
//
// Design notes:
//   - The pool maintains an `idle` channel of available conns.
//     Acquire pulls from it; Release pushes back into it.
//     A buffered channel gives us a FIFO queue + blocking semantics for free.
//   - A `total` count tracks how many conns exist (idle + in-use).
//   - When idle is empty AND total < MaxConns, Acquire creates a new conn.
//   - When idle is empty AND total >= MaxConns, Acquire WAITS.
//   - A background goroutine (the "reaper") periodically:
//       - Runs HealthCheck on idle conns
//       - Evicts conns over MaxLifetime
//       - Closes idle conns that exceed MinConns and have been idle too long
//
// Why a channel and not a slice + mutex?
// Channels give us "wait until something is available" out of the box
// via the receive operation. Implementing wait queues with mutexes + condvars
// is correct but more code. Channels are idiomatic Go.

package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Pool is the public interface. Code that uses the pool depends on this,
// not on *pool. (Dependency Inversion — easier to mock in tests.)
type Pool interface {
	Acquire(ctx context.Context) (*Conn, error)
	Stats() Stats
	Close() error
}

// Errors returned by the pool.
var (
	ErrPoolClosed     = errors.New("pool: closed")
	ErrAcquireTimeout = errors.New("pool: acquire timeout")
)

// pool is the concrete implementation. lowercase = unexported.
type pool struct {
	cfg Config

	// idle is a buffered channel of size MaxConns holding currently-idle conns.
	idle chan *Conn

	// mu protects total, closed, and the inUse map.
	mu      sync.Mutex
	total   int             // count of EXISTING conns (idle + in-use)
	inUse   map[uint64]bool // for quick membership check + Stats
	closed  bool

	nextID atomic.Uint64

	// counters for stats
	c counters

	// reaperDone is closed when the reaper goroutine exits.
	reaperDone chan struct{}
	stopReaper chan struct{}
}

// NewPool creates and warms up a new pool. Returns the Pool interface.
func NewPool(cfg Config) (Pool, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("pool: invalid config: %w", err)
	}

	p := &pool{
		cfg:        cfg,
		idle:       make(chan *Conn, cfg.MaxConns),
		inUse:      make(map[uint64]bool),
		stopReaper: make(chan struct{}),
		reaperDone: make(chan struct{}),
	}

	// Warm up: open MinConns conns up front.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < cfg.MinConns; i++ {
		c, err := p.newConn(ctx)
		if err != nil {
			// Close everything we already opened, then fail.
			_ = p.Close()
			return nil, fmt.Errorf("pool: warmup failed at #%d: %w", i, err)
		}
		// Put it directly in idle — it's already counted in total.
		p.idle <- c
	}

	// Start the background reaper.
	go p.reaper()

	return p, nil
}

// Acquire returns a conn. Caller MUST call conn.Release() (typically via defer).
//
// Selection logic:
//   1. Try to grab an idle conn immediately.
//   2. If none and total < MaxConns, create a new one.
//   3. Otherwise wait on the idle channel until a conn is released
//      OR the context is canceled OR AcquireTimeout elapses.
func (p *pool) Acquire(ctx context.Context) (*Conn, error) {
	// Check closed before doing any work.
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	p.mu.Unlock()

	// Step 1: non-blocking try on the idle channel.
	select {
	case c := <-p.idle:
		// Got one. But validate it before handing out.
		if !p.isStillUsable(ctx, c) {
			p.destroy(c, "expired or unhealthy")
			// Fall through to creating a new conn below.
		} else {
			return p.checkOut(c), nil
		}
	default:
		// idle empty.
	}

	// Step 2: try to create a new conn if under MaxConns.
	if c := p.tryCreateNew(ctx); c != nil {
		return p.checkOut(c), nil
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Step 3: block until a conn is available, with timeout.
	startWait := time.Now()
	p.c.waitCount.Add(1)

	timeoutCtx, cancel := context.WithTimeout(ctx, p.cfg.AcquireTimeout)
	defer cancel()

	for {
		select {
		case c := <-p.idle:
			p.c.waitDuration.Add(uint64(time.Since(startWait)))
			if !p.isStillUsable(ctx, c) {
				p.destroy(c, "expired during wait")
				continue
			}
			return p.checkOut(c), nil

		case <-timeoutCtx.Done():
			p.c.timeoutCount.Add(1)
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) &&
				ctx.Err() == nil {
				return nil, ErrAcquireTimeout
			}
			return nil, ctx.Err()
		}
	}
}

// tryCreateNew attempts to create a conn if the pool isn't full.
// Returns nil if we're at capacity (caller must wait instead).
func (p *pool) tryCreateNew(ctx context.Context) *Conn {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	if p.total >= p.cfg.MaxConns {
		p.mu.Unlock()
		return nil
	}
	// Optimistically reserve a slot before doing the actual dial.
	p.total++
	p.mu.Unlock()

	c, err := p.newConnUncounted(ctx)
	if err != nil {
		// Roll back the reservation.
		p.mu.Lock()
		p.total--
		p.mu.Unlock()
		return nil
	}
	return c
}

// checkOut marks a conn as in-use, fires observers, returns it.
func (p *pool) checkOut(c *Conn) *Conn {
	p.mu.Lock()
	p.inUse[c.id] = true
	p.mu.Unlock()

	c.touch()
	p.c.acquireCount.Add(1)

	if p.cfg.Observers.OnAcquire != nil {
		p.cfg.Observers.OnAcquire(c)
	}
	return c
}

// release puts a conn back. Called by Conn.Release().
func (p *pool) release(c *Conn) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		_ = c.close()
		return nil
	}
	if !p.inUse[c.id] {
		// Already released. Idempotent no-op.
		p.mu.Unlock()
		return nil
	}
	delete(p.inUse, c.id)
	p.mu.Unlock()

	c.touch()

	if p.cfg.Observers.OnRelease != nil {
		p.cfg.Observers.OnRelease(c)
	}

	// Recycle if the conn is past its lifetime.
	if p.cfg.MaxLifetime > 0 && c.age() > p.cfg.MaxLifetime {
		p.destroy(c, "exceeded MaxLifetime")
		return nil
	}

	// Try to put back in idle. This shouldn't block because idle has
	// capacity == MaxConns and we just removed it from inUse.
	select {
	case p.idle <- c:
	default:
		// Truly should never happen. If it does, destroy the conn
		// rather than leak it.
		p.destroy(c, "idle channel full")
	}
	return nil
}

// isStillUsable runs a quick check: not too old, passes health check.
func (p *pool) isStillUsable(ctx context.Context, c *Conn) bool {
	if p.cfg.MaxLifetime > 0 && c.age() > p.cfg.MaxLifetime {
		return false
	}
	// Quick health check with a tight timeout.
	hcCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if !c.healthOK(hcCtx, p.cfg.HealthCheck) {
		p.c.healthFailures.Add(1)
		if p.cfg.Observers.OnHealthFail != nil {
			p.cfg.Observers.OnHealthFail(c, errors.New("health check failed"))
		}
		return false
	}
	return true
}

// newConn creates a conn AND increments total. Used during warmup.
func (p *pool) newConn(ctx context.Context) (*Conn, error) {
	p.mu.Lock()
	p.total++
	p.mu.Unlock()
	c, err := p.newConnUncounted(ctx)
	if err != nil {
		p.mu.Lock()
		p.total--
		p.mu.Unlock()
	}
	return c, err
}

// newConnUncounted creates a conn WITHOUT touching the total counter.
// Caller is responsible for the counter (see tryCreateNew, newConn).
func (p *pool) newConnUncounted(ctx context.Context) (*Conn, error) {
	db, err := p.cfg.Connector(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	c := &Conn{
		db:        db,
		createdAt: time.Now(),
		id:        p.nextID.Add(1),
		pool:      p,
	}
	c.touch()
	if p.cfg.Observers.OnCreate != nil {
		p.cfg.Observers.OnCreate(c)
	}
	return c, nil
}

// destroy closes a conn and decrements total. Used for eviction and on errors.
func (p *pool) destroy(c *Conn, reason string) {
	p.mu.Lock()
	delete(p.inUse, c.id)
	p.total--
	p.mu.Unlock()

	if p.cfg.Observers.OnDestroy != nil {
		p.cfg.Observers.OnDestroy(c, reason)
	}
	_ = c.close()
}

// reaper runs in the background and enforces MaxIdleTime + periodic health.
func (p *pool) reaper() {
	defer close(p.reaperDone)

	if p.cfg.HealthInterval <= 0 {
		<-p.stopReaper
		return
	}

	t := time.NewTicker(p.cfg.HealthInterval)
	defer t.Stop()

	for {
		select {
		case <-p.stopReaper:
			return
		case <-t.C:
			p.reapOnce()
		}
	}
}

// reapOnce: pull each idle conn, check it, put it back or destroy.
//
// This is best-effort. We don't try to grab every idle conn — we drain
// what's currently in the channel and put back only the healthy ones.
func (p *pool) reapOnce() {
	// Snapshot the current channel by draining it into a slice.
	var batch []*Conn
loop:
	for {
		select {
		case c := <-p.idle:
			batch = append(batch, c)
		default:
			break loop
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	keep := batch[:0] // reuse underlying array
	for _, c := range batch {
		// Evict if past lifetime.
		if p.cfg.MaxLifetime > 0 && c.age() > p.cfg.MaxLifetime {
			p.destroy(c, "MaxLifetime exceeded")
			continue
		}
		// Evict idle conns over MinConns that have been idle too long.
		if p.cfg.MaxIdleTime > 0 && c.idleFor() > p.cfg.MaxIdleTime {
			p.mu.Lock()
			haveExtra := p.total > p.cfg.MinConns
			p.mu.Unlock()
			if haveExtra {
				p.destroy(c, "MaxIdleTime exceeded")
				continue
			}
		}
		// Health check.
		if !c.healthOK(ctx, p.cfg.HealthCheck) {
			p.destroy(c, "health check failed")
			continue
		}
		keep = append(keep, c)
	}

	// Put healthy conns back.
	for _, c := range keep {
		select {
		case p.idle <- c:
		default:
			p.destroy(c, "idle channel full during reap")
		}
	}
}

// Stats returns a snapshot of pool counters.
func (p *pool) Stats() Stats {
	p.mu.Lock()
	total := p.total
	inUse := len(p.inUse)
	p.mu.Unlock()
	idle := len(p.idle)
	return Stats{
		Total:          total,
		Idle:           idle,
		InUse:          inUse,
		WaitCount:      p.c.waitCount.Load(),
		WaitDuration:   p.c.waitDuration.Load(),
		AcquireCount:   p.c.acquireCount.Load(),
		TimeoutCount:   p.c.timeoutCount.Load(),
		HealthFailures: p.c.healthFailures.Load(),
	}
}

// Close shuts down the pool. After Close, all Acquires return ErrPoolClosed.
//
// Important: Close drains in-flight conns (returns to user). It does NOT
// wait for in-use conns to be released — those will be closed when the
// caller eventually releases them.
func (p *pool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	// Stop the reaper.
	close(p.stopReaper)
	<-p.reaperDone

	// Drain idle and close every conn.
loop:
	for {
		select {
		case c := <-p.idle:
			_ = c.close()
		default:
			break loop
		}
	}

	return nil
}
