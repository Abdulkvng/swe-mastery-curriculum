# Phase 01 — Expert C++ for Low-Latency Systems

## 1. The language model you must own

### Storage duration vs lifetime

Do not conflate where bytes live with whether an object is alive.

- automatic storage: usually stack-backed local variables
- dynamic storage: usually obtained through allocator / heap
- static storage: exists for program lifetime
- thread storage: one instance per thread

Object lifetime begins and ends according to construction/destruction rules. Undefined behavior can occur even if the bytes still physically exist.

### RAII

Resource Acquisition Is Initialization means a resource is tied to an object's lifetime.

Use it for memory, locks, file descriptors, sockets, mappings, timers, and temporary state.

Why interviewers care: production C++ must remain correct under exceptions, early returns, and refactors.

### Rule of zero / five

Prefer types whose members manage their own resources (rule of zero). If you manually manage ownership, understand destructor, copy ctor, copy assignment, move ctor, move assignment.

### Value categories

Know `lvalue`, `xvalue`, `prvalue`, and why `std::move` is a cast enabling move overload resolution — it does not itself move anything.

```cpp
std::vector<int> a(1'000'000);
std::vector<int> b = std::move(a); // b can steal a's buffer
```

After move, `a` is valid but its value is unspecified unless the type documents more.

## 2. Copies, moves, and allocations

A common low-latency mistake is optimizing arithmetic while silently copying or allocating.

Questions to ask:

- Is this parameter passed by value unnecessarily?
- Does this function return cause a copy? (Often elided.)
- Is a vector growing and reallocating?
- Is a string creating heap allocations?
- Is an abstraction hiding ownership work?

### `std::vector`

Usually the default high-performance sequence container because elements are contiguous.

Important concepts:

- `size()` vs `capacity()`
- geometric growth
- `reserve()`
- iterator/reference invalidation
- relocation cost
- cache locality

```cpp
std::vector<Order> orders;
orders.reserve(expected_orders);
```

Do not call `reserve()` blindly; estimate and measure.

## 3. Polymorphism

### Runtime polymorphism

Virtual functions usually involve an indirect load/call via a vtable. Cost depends on predictability, inlining loss, instruction cache, and surrounding work.

### Compile-time polymorphism

Templates allow specialization and inlining but may increase compile time and binary size.

Interview point: "virtual is slow" is too simplistic. Explain mechanism and workload.

## 4. Templates and zero-cost abstractions

Be fluent with:

- function/class templates
- concepts or `enable_if` at a basic level
- `constexpr`
- compile-time dispatch
- generic algorithms

A zero-cost abstraction means you do not pay unnecessary runtime overhead relative to a hand-written equivalent — not that abstraction is literally free in every dimension.

## 5. Memory ownership

Know when to use:

- stack object
- `std::unique_ptr`
- `std::shared_ptr`
- raw pointer/reference as non-owning view
- `std::span`
- arena / monotonic allocator

`shared_ptr` introduces reference-count traffic and often atomic operations. It is useful when ownership is truly shared, but it is rarely the first choice for a latency-sensitive hot path.

## 6. Allocators / arenas

General-purpose allocation can have variable latency and synchronization. Hot systems often preallocate or use arenas/pools.

A monotonic arena is useful when many objects have the same bulk lifetime.

Tradeoffs:

- fast allocation
- predictable layout
- reduced fragmentation
- but memory may not be individually reclaimable
- capacity planning matters

## 7. Alignment, padding, and object layout

Understand:

```cpp
struct A {
    char c;
    std::uint64_t x;
};
```

The struct may contain padding so `x` is aligned. Reordering fields can reduce size.

But do not pack objects recklessly: misaligned access can cost performance or be unsupported on some architectures.

For concurrency, alignment can intentionally separate frequently-written fields to avoid false sharing.

## 8. Undefined behavior

High-performance C++ often relies heavily on compiler optimization, so UB is especially dangerous.

Know examples:

- use-after-free
- out-of-bounds access
- signed integer overflow
- data race
- invalid shift
- dereferencing invalid pointer
- violating object lifetime / aliasing rules

"It worked in debug" is meaningless if the program has UB.

## 9. Exceptions, RTTI, and latency

Do not repeat myths. Exceptions can have near-zero cost on the non-throwing path in many implementations, but thrown exceptions are expensive and unpredictable. Many low-latency codebases avoid throwing on the hot path for control-flow and predictability reasons.

## 10. Containers under performance pressure

### `vector`
Great locality, cheap iteration, random access.

### `deque`
Segmented storage; stable-ish growth characteristics, not fully contiguous.

### `map`
Tree, ordered, pointer-heavy, `O(log n)` operations.

### `unordered_map`
Average `O(1)` lookup, but hashing, buckets, pointer chasing, rehashing, poor worst-case behavior.

### flat / sorted vector
For small-to-medium mostly-read datasets, binary search over contiguous memory can beat node-based structures despite `O(log n)` being shared with trees.

Always connect asymptotics to constants and memory hierarchy.

## 11. Compiler optimization basics

Know what `-O2` / `-O3`, LTO, PGO, and `-march=native` conceptually do.

Inspect assembly with Compiler Explorer or local compiler output:

```bash
g++ -O3 -march=native -S -masm=intel example.cpp
```

Look for:

- vector instructions
- function calls that failed to inline
- branches in hot loops
- loads/stores
- redundant work

## 12. C++ interview drills

Be able to explain:

1. move constructor vs copy constructor
2. why `std::move` can still copy
3. dangling references
4. vector reallocation invalidation
5. `unique_ptr` vs `shared_ptr`
6. virtual dispatch mechanism
7. why contiguous memory often wins
8. struct padding
9. `constexpr` vs runtime
10. UB from data races
11. placement new / arenas conceptually
12. why passing `std::string_view` can avoid allocation/copy

## Exercises

1. Write a move-only `Buffer` type using RAII.
2. Benchmark push-back with and without `reserve`.
3. Compare iterating 1M integers in `vector` vs `list`.
4. Compare virtual dispatch vs templated dispatch for a tiny operation.
5. Create two structs with identical fields but different ordering; print `sizeof`.
6. Implement a fixed-capacity object pool.

For each exercise, predict the result **before** running it.
