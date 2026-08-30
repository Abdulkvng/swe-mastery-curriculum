# Phase 05 — Performance Engineering and Measurement

## The most important phase

The job description explicitly says to profile, measure, and reason rigorously. Treat benchmarking methodology as a first-class engineering skill.

## 1. Define the metric

Possible metrics:

- throughput (ops/s)
- average latency
- p50 / p95 / p99 / p99.9
- CPU cycles/op
- instructions/op
- cache misses
- branch misses
- memory bandwidth
- context switches
- allocation count

For trading hot paths, tail latency and jitter often matter as much as mean latency.

## 2. Benchmark traps

### Dead-code elimination

If the compiler proves a computation has no observable effect, it may remove it.

Consume outputs or use a microbenchmark framework designed to inhibit incorrect elimination.

### Constant folding

If inputs are compile-time constants, the compiler may precompute results.

### Warmup

Cold instruction/data caches, page faults, dynamic linking, and frequency transitions can distort early iterations.

### Timer overhead

If an operation takes nanoseconds, timing each individual operation with a heavy clock call can dominate the measurement. Batch operations.

### Unrepresentative inputs

A benchmark with perfectly sorted or tiny data may measure the wrong behavior.

### Frequency / thermal effects

Turbo, power management, and thermal throttling create variation.

### Background noise

Scheduler migrations, other processes, interrupts, and virtualization affect tails.

## 3. Scientific loop

1. write down hypothesis
2. choose representative workload
3. baseline
4. profile
5. change one variable
6. rerun multiple times
7. compare distributions
8. verify correctness
9. explain mechanism

## 4. Linux `perf`

Typical commands:

```bash
perf stat -e cycles,instructions,branches,branch-misses,cache-references,cache-misses ./program
perf record -g ./program
perf report
```

Useful derived metric:

- IPC = instructions / cycles

Interpret carefully. Low IPC can result from memory stalls, branches, dependencies, front-end bottlenecks, and more.

## 5. Flame graphs

A flame graph shows where sampled time accumulates across call stacks. It identifies hot functions but does not itself tell you *why* a function is slow.

Pair profiles with counters and code inspection.

## 6. Hardware-counter reasoning

### High branch-miss rate
Investigate unpredictable conditionals, indirect calls, input distribution.

### High LLC/cache misses
Investigate working set, random access, pointer chasing, data layout.

### High instructions/op but decent cache behavior
Investigate algorithmic work, abstraction overhead, missed inlining/vectorization.

### Scaling stops after adding threads
Investigate contention, false sharing, memory bandwidth, NUMA, oversubscription.

## 7. Average is dangerous

Example latencies in microseconds:

`1, 1, 1, 1, 1, 1, 1, 1, 1, 100`

Average hides a severe tail.

Always know how percentile is defined and how many samples support a p99.9 estimate.

## 8. Benchmark confidence

Record:

- CPU model
- OS/kernel
- compiler/version
- flags
- affinity settings
- input distribution
- dataset size
- warm/cold state
- run count
- median/percentiles

A benchmark result without environment context is weak evidence.

## 9. Optimization across abstraction layers

Example: slow market-data handler.

Work top-down:

1. algorithm: unnecessary `O(n)` scan?
2. data structure: pointer-heavy tree?
3. allocation: per-message heap allocation?
4. copy: repeated serialization/deserialization?
5. concurrency: lock contention?
6. layout: cache-line waste?
7. compiler: missed inlining/vectorization?
8. CPU: branch/cache/TLB stalls?
9. system: thread migration / NUMA / interrupts?

## 10. Interview answer template

When asked "How would you optimize this?" do **not** immediately propose code changes.

Say:

1. define latency/throughput target
2. reproduce representative workload
3. profile baseline
4. isolate hot path
5. use counters/traces to identify bottleneck class
6. form a hypothesis
7. modify one dimension
8. compare distributions
9. validate correctness

## Exercises

- benchmark vector vs list traversal
- benchmark branchy vs branchless loop
- benchmark allocation per iteration vs preallocation
- run `perf stat` on the labs
- create a table containing cycles/op, instructions/op, cache misses, branch misses
- intentionally write a misleading benchmark, then fix it
