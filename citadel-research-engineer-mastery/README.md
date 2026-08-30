# Citadel Securities Research Engineer — Mastery Curriculum

A rigorous, hands-on course for the Research Engineer / low-latency trading engineering profile described in the job description.

> Goal: be able to **explain, implement, benchmark, optimize, and defend** every major concept in the role — not just recognize vocabulary.

## What this course maps to

The role asks you to translate quantitative models into production-grade ultra-low-latency implementations, improve trading algorithms quantitatively and computationally, reason about market microstructure, write expert-level high-performance C++, understand CPU pipelines/caches/memory models/parallel execution, profile systems under extreme constraints, and optimize across abstraction layers from algorithms to hardware.

This curriculum is built directly around those requirements.

## Mastery standard

For every topic, you should reach four levels:

1. **Explain** — define it simply and accurately.
2. **Trace** — predict what the machine or trading system does step-by-step.
3. **Implement** — write a correct version from scratch.
4. **Optimize & defend** — benchmark alternatives and justify tradeoffs with data.

If you can only explain a concept but cannot implement or measure it, you do not yet own it.

## Curriculum

| Phase | Topic | Interview outcome |
|---|---|---|
| 00 | Role map & strategy | Know exactly what the interviewer is testing |
| 01 | Modern / expert C++ | Write idiomatic, allocation-aware, performance-conscious C++ |
| 02 | CPU architecture | Reason about pipelines, branch prediction, OoO execution, ILP |
| 03 | Memory & caches | Predict cache behavior, layout data, avoid false sharing/TLB pain |
| 04 | Concurrency & C++ memory model | Atomics, ordering, lock-free reasoning, race analysis |
| 05 | Performance engineering | Benchmark correctly; use perf/flame graphs; reason from counters |
| 06 | Low-latency patterns | Tail latency, allocation avoidance, batching, affinity, NUMA |
| 07 | Market microstructure | Order books, spread, queue priority, maker/taker, adverse selection |
| 08 | Trading systems | Market data → strategy → risk → order gateway architecture |
| 09 | Quant math | Probability, stats, expectation, regression, time series intuition |
| 10 | SIMD/GPU/HPC | Vectorization, parallel kernels, accelerator tradeoffs |
| 11 | Production & risk | Correctness, determinism, kill switches, observability, failure modes |
| 12 | Interview DSA | Performance-oriented coding patterns and problem set |
| 13 | Labs & capstones | Build and optimize real components |
| 14 | Quizzes & mocks | Retrieval practice + interview simulation |

## Recommended order

### If the interview is soon (7–10 days)

Prioritize:

1. `00-role-map.md`
2. `01-cpp-mastery/README.md`
3. `02-cpu-architecture/README.md`
4. `03-memory-caches/README.md`
5. `04-concurrency-memory-model/README.md`
6. `05-performance-engineering/README.md`
7. `07-market-microstructure/README.md`
8. `08-trading-systems/README.md`
9. quizzes + mock interviews

### If you have 3–6 weeks

Do every phase in order and complete all labs.

## Daily study loop

1. **Learn for 45–60 min.** Read one section and make sure you can explain it without notes.
2. **Code for 60–90 min.** Implement or modify a lab.
3. **Measure for 30 min.** Run benchmarks, inspect assembly or counters, form a hypothesis.
4. **Quiz for 20 min.** Answer closed-book.
5. **Interview drill for 20 min.** Speak the answer out loud with a 2-minute limit.

## The performance-engineering mental model

Always ask, in this order:

1. Is the algorithm asymptotically appropriate?
2. What is the actual hot path?
3. Is work compute-bound, memory-bound, synchronization-bound, or I/O-bound?
4. What is the access pattern and working-set size?
5. What is happening to branch prediction?
6. Are allocations or syscalls occurring on the hot path?
7. Are threads contending or bouncing cache lines?
8. Does the benchmark measure the real workload?
9. What did hardware counters / traces show?
10. Did the optimization improve p50 **and** p99/p999 without breaking correctness?

## Interview rule

Never answer a performance question with "X is faster" alone. Prefer:

> "I would expect X to be faster because ___. I would verify it using ___, controlling for ___, and I would inspect ___ if the result disagreed with my hypothesis."

That is the mindset this role rewards.

## Build the labs

```bash
cd citadel-research-engineer-mastery/labs
cmake -S . -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build -j
```

Use Linux when you want `perf` hardware counters. macOS is fine for learning C++, architecture concepts, and microbenchmarks, but several low-level profiling tools differ.

## Mastery gates

Before calling yourself interview-ready, you should be able to:

- Explain stack vs heap, object lifetime, RAII, move semantics, templates, virtual dispatch cost, allocators, alignment, and undefined behavior.
- Explain cache lines, associativity, locality, TLBs, branch predictors, pipelines, superscalar/out-of-order execution, SIMD, and prefetching.
- Explain acquire/release/seq_cst and correctly build an SPSC queue.
- Design a benchmark that avoids dead-code elimination, warmup errors, timer noise, and unrealistic inputs.
- Read a simplified `perf stat` report and suggest the next experiment.
- Implement a limit order book and discuss data-structure tradeoffs.
- Explain bid/ask, spread, midpoint, depth, queue priority, maker/taker, slippage, impact, adverse selection, and inventory risk.
- Sketch an electronic trading system and identify hot-path vs cold-path components.
- Solve performance-oriented coding questions cleanly in C++ under time pressure.
- Discuss why a seemingly lower-complexity algorithm can still lose on real hardware.

## Important note

This is an educational simulator/prep course. The trading labs use toy data and intentionally omit exchange-specific protocols and production trading logic.
