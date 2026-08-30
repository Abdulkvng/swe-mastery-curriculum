# Phase 14 — Quizzes, Oral Drills, and Mock Interviews

Do these **closed-book**. Grade only after committing to an answer.

## Scoring

- 0: incorrect / cannot explain
- 1: recognizes concept
- 2: correct explanation
- 3: correct + mechanism + tradeoff + measurement method

Target average: **2.5+**, with no 0s in C++, CPU/cache, concurrency, profiling, or market microstructure.

## Mock interview structure

### Mock A — C++ + CPU (60 min)
- 10 min C++ language depth
- 15 min coding
- 15 min CPU/cache reasoning
- 10 min performance-debug scenario
- 10 min follow-ups

### Mock B — Trading systems (60 min)
- 10 min market microstructure
- 20 min design a market-data/order path
- 15 min order-book coding/data structure
- 10 min risk/reliability
- 5 min behavioral ownership

### Mock C — Quant + systems (60 min)
- 15 min probability/statistics
- 15 min concurrency
- 15 min optimization case
- 15 min research-to-production discussion

## Oral drill rule

For every answer, use:

**Definition → mechanism → example → tradeoff → measurement**.

Example: "What is false sharing?" Do not stop at the definition. Explain cache-line coherence, give a two-thread counter example, mention padding/partitioning tradeoffs, and say how you would benchmark or detect it.
