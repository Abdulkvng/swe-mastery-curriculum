// pool/conn.go
package pool

import (
	"context"
	"database/sql"
	"sync/atomic"
	"time"
)

// Conn wraps a *sql.DB with pool metadata.
//
// Why wrap? Because the pool needs to track when each conn was created
// and last used (for lifetime / idle-time eviction). The *sql.DB itself
// doesn't know any of that.
type Conn struct {
	db        *sql.DB
	createdAt time.Time
	lastUsed  atomic.Int64 // unix-nanos; atomic so no lock needed
	id        uint64       // for logging/debugging
	pool      *pool        // back-pointer for Release()
}

// DB exposes the underlying *sql.DB. The CALLER uses this to actually run queries.
//
// Notice: the field `db` is lowercase (private). Outside callers MUST go through
// this method. That's encapsulation — we control access to the underlying conn.
func (c *Conn) DB() *sql.DB {
	c.touch()
	return c.db
}

// ID returns the conn's stable ID for logging/metrics.
func (c *Conn) ID() uint64 {
	return c.id
}

// Release returns the conn to the pool. Callers should always defer this.
//
// Idempotent: releasing an already-released conn is a no-op (nil error).
func (c *Conn) Release() error {
	if c.pool == nil {
		return nil
	}
	return c.pool.release(c)
}

// touch updates the lastUsed timestamp. Called on every use.
func (c *Conn) touch() {
	c.lastUsed.Store(time.Now().UnixNano())
}

// age returns how long ago the conn was created.
func (c *Conn) age() time.Duration {
	return time.Since(c.createdAt)
}

// idleFor returns how long since the conn was last used.
func (c *Conn) idleFor() time.Duration {
	last := c.lastUsed.Load()
	if last == 0 {
		return 0
	}
	return time.Since(time.Unix(0, last))
}

// healthOK runs the configured HealthCheck and returns true if healthy.
func (c *Conn) healthOK(ctx context.Context, hc HealthCheck) bool {
	return hc(ctx, c.db) == nil
}

// close shuts down the underlying *sql.DB.
func (c *Conn) close() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	return err
}
