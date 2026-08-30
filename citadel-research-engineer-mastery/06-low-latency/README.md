# Phase 06 — Low-Latency Systems Patterns

## Latency is a distribution

A low-latency system is not just one with a low mean. You care about jitter and tails because rare slow events can arrive exactly when markets are moving fastest.

Track p50, p95, p99, p99.9, max, and outliers over time.

## Hot path vs cold path

Keep the latency-critical path narrow.

**Hot path examples:** parse normalized market update, update book, compute feature, run decision, risk check, serialize/send order.

**Cold path examples:** logging formatting, configuration reload, analytics, persistence, dashboards.

Move work off the hot path when correctness allows.

## Allocation avoidance

Possible strategies:

- pre-size vectors
- object pools / arenas
- fixed-capacity containers
- reuse buffers
- stack storage for bounded objects

The goal is predictability and reduced allocator work, not a dogma of "never allocate."

## Syscalls and kernel crossings

Syscalls, scheduling, network stack traversal, and blocking can add latency and variance. Specialized systems may use busy polling or kernel-bypass networking, but interviews usually reward knowing the mechanism and tradeoffs rather than memorizing vendor APIs.

## Batching

Batching amortizes overhead and increases throughput, but waiting for a batch adds latency. This is a central throughput-vs-latency tradeoff.

## CPU affinity and isolation

Pin critical threads to cores to reduce migration and cache disruption. Production systems may also isolate cores from unrelated work. Tradeoffs include wasted capacity and operational complexity.

## NUMA placement

Place threads and hot data together. Remote-node access adds cost and variability.

## Precomputation

If a result depends on slowly-changing state, compute outside the critical loop. Examples include lookup tables, coefficients, encoded headers, and risk limits.

## Logging

Never assume logging is free. Formatting, locks, allocation, timestamps, and I/O can dominate tiny hot paths. Consider ring-buffered/asynchronous logging and fixed-format event records where appropriate.

## Determinism

Predictable performance often matters more than peak benchmark speed. Avoid hidden work such as lazy initialization or unpredictable allocation in critical sections.

## Interview drill

You have a handler with p50 = 2 µs and p99.9 = 700 µs. What do you investigate?

Discuss:

- scheduler/context switches
- page faults
- allocation
- locks/contention
- logging
- cache/TLB misses
- IRQs/network behavior
- frequency/power changes
- occasional slow input path
- GC only if another runtime is involved

Then explain how you would collect evidence rather than guess.
