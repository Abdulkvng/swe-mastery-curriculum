# Phase 02 — Modern CPU Architecture

## The core idea

A modern CPU does not execute C++ one source line at a time. The compiler emits instructions, and the processor overlaps, speculates, reorders, and parallelizes their execution while preserving architectural behavior.

Performance engineering means learning where that machine cannot make progress.

## 1. Pipeline

Simplified pipeline:

1. fetch instructions
2. predict control flow
3. decode
4. rename / allocate
5. schedule
6. execute
7. access memory
8. retire

Real CPUs are wider and more complex, but this model is enough to reason about many interview questions.

### Throughput vs latency

- **latency**: cycles from an operation starting until its result is available
- **throughput**: how frequently independent operations can be completed

A multiply can have multi-cycle latency yet high throughput if the CPU has pipelined execution units and independent work.

## 2. Superscalar execution

Modern cores can issue multiple instructions per cycle. This requires instruction-level parallelism (ILP).

```cpp
for (...) {
    a += x[i];
    b += y[i];
}
```

The independent accumulators may expose more ILP than one dependency chain.

## 3. Out-of-order execution

The CPU can execute ready instructions before older stalled instructions, then retire results in program order.

This hides latency when independent work exists.

What defeats it?

- long dependency chains
- repeated cache misses with limited memory-level parallelism
- unpredictable branches
- serialization instructions

## 4. Branch prediction

The front end predicts which path execution will take so it can keep the pipeline full.

A misprediction discards speculative work and redirects fetch. Cost depends on microarchitecture and context; think "many cycles," not one universal number.

### Predictable vs unpredictable branch

```cpp
if (x < 128) { ... }
```

If data is almost always below 128, prediction can be excellent. Random 50/50 behavior is harder.

### Branchless is not always faster

Branchless code can perform unnecessary work or create longer dependency chains. Measure.

## 5. Front end / instruction cache

Large generated code, heavy templates, and aggressive unrolling can create instruction-cache and decode pressure.

Optimization can regress because code size matters.

## 6. Execution ports

Different operations use different execution resources. A loop may be limited by load ports, stores, integer ALUs, floating-point/vector units, or dependency latency.

At interview level, be able to say:

> "If the loop is already saturating the load/store subsystem, adding arithmetic micro-optimizations will not help much."

## 7. Data dependencies

```cpp
sum += a[i];
```

Each iteration depends on the previous `sum` value.

Unrolling into multiple partial sums can shorten the effective dependency bottleneck and improve ILP, though compilers may already do this.

## 8. Speculation

The CPU may execute instructions speculatively before knowing they are required. Correct architectural state is committed only when safe.

This is why branch prediction affects performance so strongly.

## 9. SIMD

Single Instruction Multiple Data applies one vector instruction to multiple data lanes.

Example conceptually:

- scalar: add one float at a time
- SIMD: add 4/8/16 floats per instruction depending on ISA/vector width

Vectorization likes:

- contiguous arrays
- simple loops
- independent iterations
- known alignment / aliasing properties

It dislikes:

- pointer chasing
- unpredictable control flow
- dependencies across iterations

## 10. CPU interview questions

### Why can sorted data make a branchy loop faster?
Because branch outcomes may become more predictable, reducing pipeline flushes.

### Why can an `O(n)` algorithm lose to another `O(n)` algorithm?
Different constants, locality, branches, vectorization, memory traffic, instruction count, and parallelism.

### Why can unrolling help?
Less loop-control overhead and more ILP; possibly better vectorization. It can also hurt instruction-cache behavior.

### Why does pointer chasing hurt?
Each next address may depend on a previous load, limiting memory-level parallelism and causing cache-miss latency to become exposed.

## Exercises

1. Benchmark a predictable branch vs random branch.
2. Replace a branch with arithmetic/select logic and compare.
3. Sum an array with one accumulator vs 4 independent accumulators.
4. Inspect compiler assembly at `-O0`, `-O2`, `-O3`.
5. Benchmark a tiny hot loop after aggressive manual unrolling and see if it helps.

## Mastery gate

Without notes, draw the simplified pipeline and explain how a cache miss, branch misprediction, and dependency chain each reduce useful work per cycle.
