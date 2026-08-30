# Phase 11 — Production Correctness, Risk, and Reliability

Low latency without correctness is useless.

## Invariants

Define invariants explicitly:

- book cannot have negative quantity
- best bid must not exceed best ask in a normal uncrossed snapshot (subject to feed semantics)
- internal order state transitions are valid
- risk limits cannot be bypassed
- sequence handling is monotonic according to protocol

Assert aggressively in tests and use safe production monitoring appropriate to latency constraints.

## State machines

Order lifecycle is naturally modeled as a state machine:

`new -> pending_ack -> live -> partially_filled -> filled/cancelled/rejected`

Not every transition is legal. Explicit state machines prevent "boolean soup."

## Kill switches

A system needs a reliable mechanism to stop order generation/cancel exposure under dangerous conditions. Design this independently from model logic.

## Observability

Measure:

- end-to-end latency
- stage latency
- message rates
- sequence gaps
- rejects
- queue depths
- positions/risk utilization
- heartbeat/session health

Hot-path observability should minimize allocation, locks, and synchronous I/O.

## Testing pyramid

- unit tests for math/state transitions
- property tests for invariants
- golden-data parity tests
- deterministic replay
- load tests
- failure injection
- benchmark regression tests

## Performance regressions

Track benchmark baselines in CI where hardware stability permits. Beware noisy shared runners; use thresholds and dedicated hardware for nanosecond-sensitive gates.

## Ownership

The JD emphasizes end-to-end ownership. Interview stories should show:

problem → hypothesis/design → implementation → measurement → deployment → monitoring → iteration.

Be ready for "What did you personally own?" and "What broke?"
