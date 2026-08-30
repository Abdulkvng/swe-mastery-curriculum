# Phase 03 — Memory Hierarchy, Caches, TLBs, and Layout

## Why this matters

For many low-latency workloads, the expensive operation is not arithmetic; it is waiting for data.

Think in a hierarchy:

registers → L1 → L2 → shared last-level cache → DRAM → storage/network.

Exact sizes and latency vary by CPU. The invariant is that farther storage is generally larger and slower.

## 1. Cache lines

CPUs transfer memory in cache-line-sized chunks (commonly 64 bytes on modern x86 systems, but reason from actual hardware when needed).

If you access one byte, nearby bytes often arrive too. This creates **spatial locality**.

Reusing recently accessed data creates **temporal locality**.

## 2. Why contiguous data is powerful

```cpp
for (const auto& x : vec) consume(x);
```

Contiguous elements enable:

- fewer cache lines
- hardware prefetching
- fewer pointer loads
- easier SIMD
- less allocator metadata

This is a major reason `vector` often beats linked structures in real workloads.

## 3. Working set

The working set is the data actively needed over some execution window. Performance changes sharply when the working set stops fitting in a cache level.

Benchmark across sizes: 4 KB, 32 KB, 256 KB, 2 MB, 16 MB, etc. You can often see latency/bandwidth regime changes.

## 4. Associativity and conflicts

Caches are set-associative, not fully associative. Certain address patterns can map repeatedly to the same sets and evict one another even when total capacity seems sufficient.

You rarely need to calculate set indices in an interview, but know **conflict misses** exist.

## 5. TLB

The Translation Lookaside Buffer caches virtual-to-physical address translations.

Large/random memory footprints can cause TLB misses, adding page-walk overhead.

Helpful ideas:

- locality
- fewer pages touched
- huge pages in specialized environments
- avoiding giant sparse structures

## 6. Array of Structures vs Structure of Arrays

AoS:

```cpp
struct Quote { double bid, ask, size; };
std::vector<Quote> quotes;
```

SoA:

```cpp
std::vector<double> bids, asks, sizes;
```

If a hot loop only uses bids, SoA loads less irrelevant data and may vectorize better. If code usually consumes all fields together, AoS may be convenient and locality-friendly.

There is no universal winner.

## 7. False sharing

Two threads update different variables that happen to share the same cache line. Coherence still moves ownership of the whole line between cores.

Symptoms:

- scaling collapses with multiple threads
- CPU usage high
- no obvious mutex contention

Mitigation:

- partition writes
- align/pad per-thread counters
- reduce shared writable state

```cpp
struct alignas(64) Counter {
    std::atomic<std::uint64_t> value{0};
};
```

Do not assume 64 blindly for every platform; it is illustrative.

## 8. NUMA

On multi-socket/NUMA systems, memory access time depends on which CPU/node owns the memory.

Key ideas:

- first-touch placement
- thread affinity
- local vs remote memory
- avoiding unnecessary cross-node sharing

## 9. Prefetching

Hardware prefetchers detect regular access patterns. Sequential arrays are friendly; random linked structures are not.

Manual prefetch can help in specialized cases but can also:

- fetch useless lines
- evict useful data
- consume bandwidth
- arrive too early/late

Measure.

## 10. Memory bandwidth vs latency

A streaming workload may be bandwidth-bound: many outstanding accesses keep memory busy.

A pointer-chasing workload may be latency-bound: the next address is unknown until the current load returns.

Same DRAM, very different bottleneck.

## 11. Interview drills

1. Why can a 24-byte struct be worse than a 16-byte struct in a tight scan?
2. Why can padding improve multicore performance while using more memory?
3. Why might binary search on a vector lose to linear scan for tiny `n`?
4. Why do random accesses degrade as the array becomes larger than LLC?
5. What is a TLB miss?
6. AoS vs SoA for a pricing kernel: what determines the choice?

## Exercises

- stride benchmark: access every 1, 2, 4, 8, 16... elements
- random vs sequential access benchmark
- AoS vs SoA benchmark
- false-sharing benchmark with 2 threads
- working-set sweep and plot ns/access

Always record CPU, compiler flags, data size, iteration count, and whether data was warm or cold.
