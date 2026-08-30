# 100 Interview Questions

## C++
1. Implement a move-only aligned buffer.
2. Explain why `std::move` isn't a move.
3. Design an object pool.
4. What does vector growth cost?
5. When would you avoid a hash table?
6. Explain vtables and devirtualization.
7. What UB have you seen in real systems?
8. Explain strict ownership with `unique_ptr`.
9. `string` vs `string_view` on a parser hot path?
10. How would you inspect compiler output?

## CPU / memory
11. Walk a load from L1 miss to DRAM conceptually.
12. Why does cache-line size matter?
13. Give a false-sharing example.
14. Why can SoA vectorize better than AoS?
15. Why can linear scan beat binary search for tiny arrays?
16. What is branch prediction learning?
17. Why can sorted input change runtime dramatically?
18. What is the TLB?
19. What is NUMA?
20. How does pointer chasing expose memory latency?

## Concurrency
21. Mutex vs spinlock?
22. What is a data race?
23. Atomicity vs visibility/order?
24. Explain release/acquire.
25. When is relaxed ordering useful?
26. Build an SPSC ring buffer.
27. Why separate head and tail cache lines?
28. What is ABA?
29. Lock-free vs wait-free?
30. Why can more threads make code slower?

## Performance
31. How do you benchmark a 20ns function?
32. How do you stop dead-code elimination?
33. What counters would you collect first?
34. Mean improved but p99 regressed. Ship?
35. What causes scheduler jitter?
36. How do you find allocation on a hot path?
37. What does a flame graph tell you?
38. What does it not tell you?
39. What if profiling changes timing?
40. How do you make performance tests reproducible?

## Microstructure
41. Draw an order book.
42. Define spread/mid.
43. What is depth?
44. Why does queue priority matter?
45. Explain adverse selection.
46. What is inventory risk?
47. What is market impact?
48. What is a stale quote?
49. Why can latency have economic value?
50. What events update a book?

## Trading systems
51. Design a market-data handler.
52. How detect feed gaps?
53. Design normalized market events.
54. What is hot vs cold path?
55. Where do risk checks live?
56. Design kill switch behavior.
57. How would you replay an incident?
58. How test research/production parity?
59. How handle order rejects?
60. How instrument without perturbing latency?

## Algorithms/data structures
61. Implement fixed-size ring buffer.
62. Design top-K stream.
63. Implement median stream.
64. Design price-level map.
65. Hash map vs sorted vector?
66. Heap vs ordered tree?
67. Solve interval scheduling.
68. Solve shortest path.
69. Build LRU cache.
70. Optimize a slow O(n²) matching loop.

## Quant
71. Expected value puzzle.
72. Conditional probability puzzle.
73. Bayes update puzzle.
74. Variance of independent sum.
75. Explain covariance.
76. Explain regression coefficient.
77. What is leakage?
78. Why time-series split?
79. Floating-point reproducibility?
80. How validate a signal?

## HPC
81. Explain SIMD.
82. Why didn't compiler vectorize?
83. Amdahl's law.
84. Memory-bound vs compute-bound?
85. Why GPU for simulation but not necessarily hot-path decisions?
86. What is arithmetic intensity?
87. Parallel reduction issues?
88. Data transfer overhead?
89. What is kernel divergence?
90. How benchmark GPU end-to-end fairly?

## Ownership / production
91. Tell me about a performance bottleneck you owned.
92. How did you prove the bottleneck?
93. What tradeoff did you reject?
94. Tell me about a production failure.
95. How did monitoring help?
96. How did you validate correctness after optimization?
97. How do you disagree with a researcher constructively?
98. How do you prioritize speed vs maintainability?
99. What would you learn first entering electronic trading?
100. Why Research Engineering rather than general SWE?
