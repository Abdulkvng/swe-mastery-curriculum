// pool/stats.go
package pool

import "sync/atomic"

// Stats is a snapshot of pool counters at a point in time.
//
// Important: snapshot, not live — callers shouldn't assume the numbers
// are coherent. (Reading multiple atomics non-atomically gives near-coherent
// values, which is fine for monitoring.)
type Stats struct {
	Total          int    // currently-existing conns
	Idle           int    // available for Acquire
	InUse          int    // checked out
	WaitCount      uint64 // # times Acquire had to wait
	WaitDuration   uint64 // total wait time, nanoseconds
	AcquireCount   uint64 // # successful Acquires
	TimeoutCount   uint64 // # Acquire timeouts
	HealthFailures uint64 // # health-check failures
}

// counters lives inside the pool. Atomic so we don't need a lock for stats.
type counters struct {
	waitCount      atomic.Uint64
	waitDuration   atomic.Uint64
	acquireCount   atomic.Uint64
	timeoutCount   atomic.Uint64
	healthFailures atomic.Uint64
}
