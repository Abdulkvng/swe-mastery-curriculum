# btree-from-scratch

> An in-memory B-tree in Go. The data structure underneath every database index you'll touch.

## Why build this

You won't write a B-tree at work. But if you understand it deeply, you understand:
- Why Postgres uses 8KB pages
- Why your composite index `(state, age)` doesn't help `WHERE age = 21`
- Why every modern database has a B-tree (or LSM) at the bottom
- Why CLRS chapter 18 is hard the first time you read it

## Run

```bash
go test ./...
go test -bench=. ./...
```

## Concepts demonstrated

- **Generics** in Go (`BTree[K, V]`)
- **User-supplied comparator** (`less func(a, b K) bool`) — works for any K
- **Splitting full nodes top-down** during insert (CLRS strategy)
- **Range queries** in O(log n + r)

## What's NOT included

- A proper `Delete` (the slow rebuild is here to keep tests passing). Implementing the real CLRS delete is your stretch goal — it's the longest part of the algorithm.
- Concurrency. Postgres' B-tree (Lehman-Yao) is concurrent. We're single-threaded.
- Disk persistence. We're in-memory. Translating to a page-based on-disk layout is its own multi-month project.

## Suggested next reading

- [Postgres' nbtree README](https://github.com/postgres/postgres/blob/master/src/backend/access/nbtree/README) — the real thing, in C, with concurrency.
- "Modern B-Tree Techniques" by Goetz Graefe — the survey.
