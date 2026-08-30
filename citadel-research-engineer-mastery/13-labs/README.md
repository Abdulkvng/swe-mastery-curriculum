# Phase 13 — Labs and Capstones

These are the most important pieces of the course. Reading creates familiarity; building creates skill.

## Lab 1 — Cache locality benchmark

Compare sequential vs random access and sweep working-set sizes.

**Mastery:** explain regime changes using cache/TLB behavior.

## Lab 2 — Branch prediction

Compare a branch on predictable input vs shuffled/random input.

**Mastery:** verify branch-miss differences on Linux with `perf`.

## Lab 3 — False sharing

Two threads increment adjacent counters vs cache-line-separated counters.

**Mastery:** explain coherence traffic.

## Lab 4 — SPSC ring buffer

Implement fixed-capacity producer/consumer queue with acquire/release ordering.

**Mastery:** give a correctness argument.

## Lab 5 — Limit order book

Implement simplified price-level book with add/cancel and best-bid/best-ask.

**Mastery:** benchmark at least two data structures.

## Lab 6 — Market-data pipeline

Generate synthetic updates, parse, update book, compute midpoint/imbalance, risk-check a toy decision, send to a sink.

Instrument per-stage latency.

## Capstone A — 10× optimization challenge

Start with intentionally slow pipeline:

- per-event heap allocation
- strings for side/type
- node-heavy structures
- locking between every stage
- synchronous formatted logging

Profile and improve it systematically. Do not change multiple things without measuring.

Create `results.md` containing baseline, hypothesis, change, counters, result, and tradeoff for each optimization.

## Capstone B — Research-to-production parity

Write a simple reference model in Python and C++ production implementation. Feed identical event data and compare all outputs within a documented numeric tolerance.

## Capstone C — Interview defense

Give a 10-minute explanation of your optimized architecture. A reviewer should interrupt with:

- Why this data structure?
- Why not lock-free?
- What does p99 look like?
- How did you know caches were the bottleneck?
- What happens if volume doubles?
- What happens on a feed gap?
- How do you prevent bad orders?

If you cannot defend it, revisit the relevant phase.
