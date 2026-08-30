# Phase 04 — Concurrency, Atomics, and the C++ Memory Model

## 1. Start with correctness

A data race in C++ is undefined behavior unless synchronization rules make the accesses safe.

Do not begin lock-free programming by trying to remove locks. Begin by defining:

- which thread owns which data
- which state is shared
- which ordering guarantees are required
- how shutdown/lifetime works

## 2. Mutexes are not failure

A mutex is often the correct design. Lock-free code adds complexity and can perform worse under some workloads.

Use lock-free structures when you have a measured reason: tail-latency constraints, known topology, contention pattern, or progress requirements.

## 3. Atomics

`std::atomic<T>` provides indivisible operations and participates in the language memory model.

Common operations:

- load
- store
- exchange
- fetch_add
- compare_exchange_weak/strong

Atomic does **not** automatically make a multi-step algorithm correct.

## 4. Memory ordering intuition

### `memory_order_relaxed`

Atomicity without cross-variable synchronization ordering. Useful for counters or carefully-designed algorithms where ordering comes from elsewhere.

### Release / acquire

Producer writes ordinary data, then performs a release store. Consumer observes the corresponding value with an acquire load. Earlier producer writes become visible to the consumer after the acquire.

Conceptual pattern:

```cpp
// producer
payload = value;
ready.store(true, std::memory_order_release);

// consumer
if (ready.load(std::memory_order_acquire)) {
    use(payload);
}
```

### Sequential consistency

Strongest commonly-used ordering and easiest to reason about conceptually. It can constrain optimization more than weaker orderings.

Interview rule: correctness first. Do not choose relaxed ordering just because it sounds fast.

## 5. Happens-before

The memory model defines relationships that make writes visible and races well-defined.

You should be able to explain:

- sequenced-before within a thread
- synchronization between threads
- happens-before as the transitive correctness relation you care about

## 6. Compare-and-swap

CAS updates a value only if it equals an expected value.

It is the core primitive behind many lock-free structures.

Know:

- failed CAS loops
- contention
- ABA problem conceptually
- weak vs strong CAS (weak may spuriously fail and is often used in loops)

## 7. SPSC queue

A single-producer/single-consumer ring buffer is a perfect interview-prep structure because it combines:

- fixed-capacity array
- indices
- cache lines
- atomics
- acquire/release
- wraparound
- avoiding allocation

The lab implements one.

Important design detail: producer owns `head`, consumer owns `tail`, and synchronization is only needed when publishing/observing progress.

## 8. False sharing revisited

Even correct atomics can be slow if producer and consumer repeatedly write atomics on the same cache line.

Place independently-written hot fields on separate lines when justified.

## 9. Spin vs block

A spin wait burns CPU but can avoid scheduler/wakeup latency for extremely short waits. Blocking yields CPU but may involve kernel scheduling.

Latency-sensitive systems may spin on dedicated cores — but this is a system-level tradeoff, not a universal best practice.

## 10. Thread affinity

Pinning a latency-sensitive thread can reduce migration, improve cache warmth, and make timing more predictable. Costs include reduced scheduler flexibility and operational complexity.

## 11. Interview questions

1. What is a data race?
2. Atomicity vs ordering?
3. What does acquire/release guarantee?
4. When is relaxed ordering safe?
5. Why is a lock-free algorithm not automatically faster?
6. What is false sharing?
7. SPSC vs MPMC complexity?
8. What is ABA?
9. Why can a spinlock be disastrous under oversubscription?
10. Why pin threads?

## Exercises

- implement atomic counter with relaxed ordering
- build producer/consumer handoff with release/acquire
- intentionally create false sharing and benchmark it
- complete the SPSC lab
- change memory orderings, reason about correctness **before** benchmarking

Mastery means you can explain the ordering proof of your queue, not merely make tests pass.
