# C++ Mastery Quiz

## Questions

1. What does `std::move` actually do?
2. Why can `std::move(x)` still result in a copy?
3. What is RAII and why is it useful in production systems?
4. Difference between `size()` and `capacity()` for `std::vector`?
5. What operations can invalidate vector references/iterators?
6. Why might `std::vector` beat `std::list` for traversal and even some insertion workloads?
7. When would `std::string_view` be dangerous?
8. Explain object lifetime vs storage duration.
9. Why is signed integer overflow dangerous in C++?
10. What is a data race under the C++ model?
11. Runtime polymorphism vs compile-time polymorphism?
12. Why might `shared_ptr` be undesirable on a hot path?
13. What is alignment? What is padding?
14. Why can field order change `sizeof(struct)`?
15. What is copy elision?
16. When would an arena allocator help?
17. Why is `unordered_map` not simply "O(1), therefore fastest"?
18. Explain `constexpr`.
19. What is placement new conceptually?
20. What does `const` protect and what does it not protect?

## Answer key

1. It casts to an rvalue/xvalue expression to enable overloads that may move; it performs no transfer by itself.
2. The type may lack an accessible move operation, or the expression may be const so a move constructor requiring non-const rvalue cannot bind.
3. Tie resource lifetime to object lifetime so cleanup happens deterministically on scope exit, including early returns/exceptions.
4. Size is number of live elements; capacity is allocated element space before reallocation is needed.
5. Reallocation invalidates all; erase/insert may invalidate at/after operation depending on case; know container rules.
6. Contiguous locality, fewer allocations/pointers, hardware prefetch and lower memory overhead can overwhelm theoretical node-operation advantages.
7. It is non-owning; dangerous when referenced storage expires or mutates unexpectedly.
8. Storage is where/how long bytes are reserved; lifetime is when a valid object exists in those bytes.
9. Signed overflow is UB, allowing optimizer assumptions that can invalidate intuitive behavior.
10. Conflicting unsynchronized accesses from threads with at least one write, absent required happens-before, produce UB.
11. Virtual dispatch chooses at runtime; templates/static polymorphism choose at compile time and often enable inlining, with code-size/compile-time tradeoffs.
12. Shared reference counting, ownership complexity, allocation patterns, and often atomic refcount traffic.
13. Alignment constrains address multiples; padding is unused bytes inserted to satisfy alignment/layout.
14. Different ordering changes padding requirements.
15. Construction directly in destination rather than materializing/copying temporary, mandatory in some modern C++ cases.
16. Many allocations with common/bulk lifetime and predictable capacity.
17. Hashing, collisions, rehashing, bucket indirection, cache misses, memory overhead and worst-case behavior matter.
18. Enables values/functions to be evaluated at compile time when used with constant-expression inputs, subject to language rules.
19. Constructs an object into already-provided storage; lifetime management remains explicit.
20. Depends on placement: const object/reference prevents mutation through that interface; it does not automatically imply thread safety or deep immutability.
