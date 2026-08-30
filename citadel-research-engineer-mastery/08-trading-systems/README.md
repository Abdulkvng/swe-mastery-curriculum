# Phase 08 — Electronic Trading Systems Architecture

## Simplified architecture

```text
Exchange Feed
    |
    v
Market Data Receiver -> Decoder/Normalizer -> Order Book
                                           |
                                           v
                                    Features / Model
                                           |
                                           v
                                      Strategy Logic
                                           |
                                           v
                                       Risk Checks
                                           |
                                           v
                                      Order Gateway
                                           |
                                           v
                                        Exchange

Cold-path: capture, metrics, logs, research replay, dashboards
```

## Market-data handler

Responsibilities may include:

- packet/message validation
- sequence tracking
- decoding binary protocol
- timestamping
- normalization
- gap detection/recovery
- updating book state

Latency-sensitive design avoids repeated parsing, allocation, and copying.

## Research-to-production model translation

A researcher may prototype in Python/NumPy. Production engineering must answer:

- exact feature definitions?
- floating-point precision?
- missing data behavior?
- initialization/warmup?
- update frequency?
- state reset?
- deterministic ordering?
- numerical parity tolerance?

Create golden datasets: feed identical historical events into research and production implementations and compare intermediate states and outputs.

## Strategy / decision layer

Keep business logic separated enough to test it deterministically. Measure whether abstractions introduce unacceptable hot-path cost before removing them.

## Pre-trade risk

Examples of controls:

- max order quantity
- max position / notional
- price collars
- duplicate/order-rate controls
- self-trade prevention depending on system/venue
- kill switch

Risk checks must be fast but cannot be skipped for speed.

## Order gateway

Responsibilities include:

- order serialization
- sequence/session state
- acknowledgments
- rejects
- cancel/replace flow
- reconnection/recovery

Correctness and state-machine design are critical.

## Replay and simulation

A strong production stack can replay recorded events deterministically for debugging, regression tests, and research validation.

Useful properties:

- event ordering preserved
- timestamps represented explicitly
- random components seeded
- outputs diffable

## Failure modes

Be able to discuss:

- market-data packet loss / sequence gap
- stale book
- disconnected order session
- duplicate message
- partial fill
- reject
- clock drift
- process crash
- risk-limit breach
- overloaded logging/telemetry

## Interview design exercise

Design a system processing millions of market-data updates per second while producing orders with strict tail latency.

Talk through:

1. requirements and latency budget
2. thread ownership
3. feed parsing
4. book data structures
5. model state
6. risk checks
7. gateway
8. queues between components
9. CPU affinity / NUMA
10. observability without poisoning hot path
11. replay/recovery
12. correctness tests

The strongest answer includes both **speed and operational safety**.
