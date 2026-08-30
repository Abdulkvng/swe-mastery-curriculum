# Phase 10 — HPC, SIMD, Parallel Programming, GPUs

## CPU vectorization

SIMD works best when one operation applies across independent contiguous elements.

Compiler-friendly loop:

```cpp
for (std::size_t i = 0; i < n; ++i)
    out[i] = a[i] * b[i] + c[i];
```

Potential blockers:

- aliasing uncertainty
- loop-carried dependency
- irregular branches
- scattered memory access
- function calls not inlined

Use compiler vectorization reports and assembly to verify; do not assume.

## Amdahl's law

If fraction `P` of a workload parallelizes across N workers:

`speedup = 1 / ((1-P) + P/N)`

A tiny serial fraction caps scaling.

## Roofline intuition

Performance can be limited by compute throughput or memory bandwidth. Arithmetic intensity = operations per byte moved. Low-intensity kernels often remain memory-bound even on powerful compute hardware.

## GPUs

Strengths:

- massive data parallelism
- high memory bandwidth
- throughput-oriented kernels

Costs:

- transfer/launch overhead
- different memory hierarchy
- synchronization
- kernel divergence
- operational complexity

A GPU can be excellent for large batch research/simulation and wrong for a tiny per-message latency path.

## Heterogeneous design

Ask:

- batch size?
- transfer frequency?
- deadline?
- arithmetic intensity?
- parallelism?
- precision?
- determinism?

## Parallel reductions

Reductions need tree-like combining. Floating-point result may differ because addition is not associative. Reproducibility requirements can constrain optimization.

## Interview questions

1. Why doesn't adding 2× cores guarantee 2× speedup?
2. CPU SIMD vs GPU parallelism?
3. Why might GPU acceleration increase end-to-end latency?
4. What is arithmetic intensity?
5. How do memory bandwidth and cache behavior affect vectorization wins?
6. Why can parallel floating-point reduction produce a different result?
