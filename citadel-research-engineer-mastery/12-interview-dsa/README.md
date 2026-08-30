# Phase 12 — Coding / DSA for a Performance Engineer

You still need clean interview coding. Use C++ for this track unless your interviewer explicitly allows another language and you have chosen it.

## Priority patterns

### Arrays / prefix / two pointers
- Two Sum
- Product Except Self
- Subarray Sum Equals K
- Maximum Subarray
- Container With Most Water

### Hashing
- Longest Consecutive Sequence
- Group Anagrams
- Top K Frequent Elements

### Heaps
- Kth Largest
- Merge K Sorted Lists
- Median from Data Stream

### Intervals
- Merge Intervals
- Insert Interval
- Meeting Rooms II

### Binary search
- Search Rotated Array
- Koko Eating Bananas
- lower_bound / upper_bound usage

### Trees / graphs
- BFS/DFS
- Lowest Common Ancestor
- Course Schedule
- Dijkstra

### Dynamic programming
- House Robber
- Coin Change
- LIS
- Weighted Interval Scheduling

## Performance twist

For every solution also answer:

1. time complexity?
2. auxiliary memory?
3. allocations?
4. cache locality?
5. can input bounds permit a flat array instead of hash map?
6. worst-case behavior?
7. what changes if this is called 10 million times?

## C++ fluency checklist

You should type without hesitation:

```cpp
std::vector<int>
std::array<T, N>
std::unordered_map<K,V>
std::map<K,V>
std::priority_queue<T>
std::queue<T>
std::deque<T>
std::sort
std::lower_bound
std::upper_bound
```

Know lambda syntax, custom comparators, references, `const`, range loops, and integer overflow concerns.

## Timed format

For each problem:

- 3 min clarify + examples
- 5 min derive
- 20 min code
- 5 min test
- 5 min optimize/discuss hardware implications

Do not prematurely optimize before correctness.
