# CPU / Performance Quiz

## Questions

1. Cache latency vs memory bandwidth: difference?
2. What is a cache line and why does it matter?
3. Define spatial and temporal locality.
4. Why is pointer chasing hard for CPUs?
5. What happens conceptually on a branch misprediction?
6. Why is branchless code not always faster?
7. What is out-of-order execution trying to achieve?
8. What is a dependency chain?
9. Define IPC. Is higher IPC always better application performance?
10. What is a TLB?
11. What is false sharing?
12. How would you prove false sharing is causing a regression?
13. Why benchmark multiple working-set sizes?
14. Why can an optimized benchmark be invalid because of dead-code elimination?
15. Why report percentiles?
16. What makes p99.9 statistically tricky?
17. What is CPU affinity useful for?
18. What is NUMA?
19. When can batching hurt?
20. What would you inspect after seeing high branch-miss counters?
21. What would you inspect after high LLC misses?
22. Why can `-O3` sometimes lose to `-O2`?
23. What is SIMD?
24. What workload prevents easy SIMD vectorization?
25. Give a rigorous optimization workflow.

## Expected themes

A strong answer connects mechanisms to measurement. For Q25: define target → representative workload → baseline → profile/counters → hypothesis → one controlled change → rerun distribution → verify correctness → document tradeoff.
