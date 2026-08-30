# Phase 00 — Decode the Role

## What the title really means

A Research Engineer in electronic trading sits between quantitative research and production systems engineering. You are expected to understand a model deeply enough to preserve its mathematical intent while rewriting it so it can run predictably, safely, and extremely quickly in production.

The interview is therefore likely to probe the intersection of:

- C++ language depth
- low-level systems knowledge
- performance reasoning
- algorithms / data structures
- concurrency
- quantitative thinking
- trading / market-microstructure intuition
- communication with researchers

## Job-description-to-skill map

### "Translate mathematical models into production-grade, ultra-low latency implementations"

You should be ready to discuss:

- turning vectorized/research code into deterministic production code
- numerical precision and reproducibility
- precomputation
- choosing data layouts and types
- avoiding dynamic allocation in hot paths
- minimizing copies
- branch reduction
- batching vs latency tradeoffs
- testing model parity between research and production

### "Improve trading algorithms quantitatively and for performance"

Two axes:

**Quantitative improvement**: better estimator, feature, signal, calibration, model assumptions, execution logic, or risk tradeoff.

**Computational improvement**: same mathematical result with lower latency, lower variance, less CPU, fewer cache misses, fewer branches, better vectorization, or better parallelism.

The strongest answers keep those axes separate and then explain how they interact.

### "Market microstructure"

Know:

- limit order book
- bids, asks, spread, midpoint
- price-time priority
- market vs limit orders
- queue position
- depth / liquidity
- maker/taker
- slippage and market impact
- adverse selection
- inventory risk
- latency and stale quotes

### "High-performance C++"

Expect depth beyond syntax:

- value categories and move semantics
- object lifetime
- RAII
- templates / compile-time polymorphism
- virtual dispatch tradeoffs
- memory allocation
- `std::vector` growth and locality
- `std::unordered_map` tradeoffs
- custom allocators / arenas
- alignment / padding
- atomics and memory ordering
- undefined behavior
- compiler optimizations
- reading generated assembly at a basic level

### "Modern CPU architecture"

Know the causal chain:

source code → compiler → instructions → front end → branch prediction → decode → scheduling → execution ports → loads/stores → caches/TLB/DRAM → retirement.

You do not need to be a CPU designer. You do need to reason about why a loop stalls.

### "Profile, measure, reason rigorously"

This wording is critical. They do not merely want optimization folklore.

A good loop is:

1. define the metric
2. create representative workload
3. establish baseline
4. profile
5. form hypothesis
6. change one thing
7. rerun
8. check statistical significance / noise
9. inspect tail latency
10. validate correctness

## How to answer technical performance questions

Use **H-M-V-T**:

- **Hypothesis** — what you think is happening.
- **Mechanism** — CPU/memory/compiler reason.
- **Verification** — benchmark/counter/trace you would use.
- **Tradeoff** — what could get worse.

Example:

> "A flat vector may beat a tree for this small working set. My hypothesis is better spatial locality and fewer pointer-chasing cache misses. I would compare representative workloads, inspect cache-miss and branch-miss counters, and check whether insertion cost becomes unacceptable at larger sizes."

## Interview dimensions

Score yourself 0–3 on each:

| Dimension | 0 | 1 | 2 | 3 |
|---|---|---|---|---|
| C++ | unfamiliar | writes code | understands internals | performance-level depth |
| CPU/cache | vocabulary only | basic explanation | predicts behavior | profiles and optimizes |
| concurrency | threads only | mutexes | atomics/orderings | lock-free reasoning |
| performance | guesses | benchmarks | profiles | counter-driven methodology |
| microstructure | none | definitions | order-book intuition | system/performance implications |
| quant math | rusty | formulas | derives | applies under uncertainty |
| DSA | slow | solves mediums | clean C++ | performance-aware tradeoffs |

Your goal is mostly 2s and 3s before interview day.

## Questions you should be able to answer by the end

1. Why might `std::vector` outperform `std::list` even when the list has theoretically cheaper insertion?
2. What exactly happens after an L1 cache miss?
3. What causes branch misprediction cost?
4. Why can false sharing destroy multicore scaling?
5. What is the difference between `memory_order_relaxed`, acquire/release, and sequential consistency?
6. How would you benchmark two order-book data structures fairly?
7. What metrics matter besides average latency?
8. What happens when a market maker's quote becomes stale?
9. What is adverse selection?
10. How would you preserve parity between a research model and production implementation?
