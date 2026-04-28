// pool/config.go
package pool

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Connector abstracts how a new connection is created.
// In production, you'd pass a function that calls sql.Open.
// In tests, you'd pass a fake that returns a stub. (Dependency Inversion!)
type Connector func(ctx context.Context) (*sql.DB, error)

// HealthCheck is the Strategy pattern: pluggable connection-validation.
// Default implementation pings the DB. You can replace it with anything
// that returns nil (healthy) or an error (unhealthy → conn gets recycled).
type HealthCheck func(ctx context.Context, db *sql.DB) error

// Observers — Observer pattern. Optional callbacks the user can hook into.
type Observers struct {
	OnAcquire     func(c *Conn)
	OnRelease     func(c *Conn)
	OnCreate      func(c *Conn)
	OnDestroy     func(c *Conn, reason string)
	OnHealthFail  func(c *Conn, err error)
}

// Config controls Pool behavior. All fields have sensible defaults via DefaultConfig.
type Config struct {
	// Connector creates a new underlying *sql.DB. Required.
	Connector Connector

	// MinConns: pool keeps at least this many idle conns warm.
	MinConns int
	// MaxConns: pool will never have more than this many total conns.
	MaxConns int

	// AcquireTimeout: how long Acquire blocks before giving up.
	AcquireTimeout time.Duration

	// MaxLifetime: a conn older than this gets recycled on next release.
	// Helps with rotating credentials and cleaning up server-side state.
	MaxLifetime time.Duration

	// MaxIdleTime: an idle conn that hasn't been used for this long gets closed.
	// Frees server resources during quiet periods.
	MaxIdleTime time.Duration

	// HealthInterval: how often the background reaper checks for dead conns.
	HealthInterval time.Duration

	// HealthCheck: function that validates a single conn. Defaults to Ping.
	HealthCheck HealthCheck

	// Observers: optional event hooks (Observer pattern).
	Observers Observers
}

// DefaultConfig returns sane defaults. Caller MUST set Connector.
func DefaultConfig() Config {
	return Config{
		MinConns:       2,
		MaxConns:       10,
		AcquireTimeout: 5 * time.Second,
		MaxLifetime:    30 * time.Minute,
		MaxIdleTime:    5 * time.Minute,
		HealthInterval: 30 * time.Second,
		HealthCheck:    PingHealthCheck,
	}
}

// PingHealthCheck is the default health check. Calls db.PingContext.
func PingHealthCheck(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

// validate checks the config. Encapsulates input validation in one place.
func (c *Config) validate() error {
	if c.Connector == nil {
		return errors.New("Connector is required")
	}
	if c.MinConns < 0 {
		return errors.New("MinConns must be >= 0")
	}
	if c.MaxConns <= 0 {
		return errors.New("MaxConns must be > 0")
	}
	if c.MinConns > c.MaxConns {
		return errors.New("MinConns cannot exceed MaxConns")
	}
	if c.AcquireTimeout < 0 {
		return errors.New("AcquireTimeout cannot be negative")
	}
	if c.HealthCheck == nil {
		c.HealthCheck = PingHealthCheck
	}
	return nil
}
