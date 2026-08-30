# Phase 07 — Market Microstructure

This phase teaches the vocabulary and mechanics needed to reason about electronic markets. It is educational, not trading advice.

## 1. Limit order book

A limit order book contains resting buy and sell interest organized by price.

- **bid**: price buyers are willing to pay
- **ask/offer**: price sellers are willing to accept
- **best bid**: highest bid
- **best ask**: lowest ask
- **spread**: best ask − best bid
- **midpoint**: (best bid + best ask) / 2

Example:

```text
ASKS
101.02 : 400
101.01 : 200  <- best ask
----------------
100.99 : 350  <- best bid
100.98 : 900
BIDS
```

Spread = 0.02; midpoint = 101.00.

## 2. Market vs limit orders

A **limit order** specifies a price constraint and can rest on the book.

A **marketable order** immediately consumes available opposing liquidity up to its constraints.

## 3. Price-time priority

Many matching engines prioritize better prices first, then earlier orders at the same price. Exact exchange rules vary.

Queue position matters because being first at a price level can affect fill probability.

## 4. Depth and liquidity

Depth describes available quantity at one or more price levels. A market can have a narrow spread yet little depth.

## 5. Maker vs taker

A maker generally adds resting liquidity; a taker removes it. Fee/rebate schedules differ by venue and product.

## 6. Slippage and impact

**Slippage**: realized execution differs from a reference/expected price.

**Market impact**: your trading activity moves prices or consumes liquidity.

Large urgent orders may walk through multiple levels.

## 7. Adverse selection

You provide liquidity at 100.00 and are immediately filled just before the fair value falls to 99.90. The counterparty may have traded against a stale quote or superior short-term information.

Market makers therefore balance spread capture against adverse-selection and inventory risk.

## 8. Inventory risk

A market maker accumulating a large directional position becomes exposed to price moves. Pricing and quote sizes may be adjusted to manage inventory.

## 9. Stale quotes and latency

Suppose your model says fair value moved from 100 to 99.80, but your old bid at 99.98 is still live. Faster participants may trade with that stale bid before you cancel/update it.

This links systems latency directly to economics.

## 10. Order-book events

A simplified feed might contain:

- add order
- cancel order
- modify/replace
- trade/execution

Some feeds provide market-by-price rather than individual orders. Your data structures depend on the semantics.

## 11. Microstructure features

Educational examples:

### Midpoint
`mid = (best_bid + best_ask)/2`

### Spread
`spread = best_ask - best_bid`

### Top-of-book imbalance

`imbalance = (bid_size - ask_size) / (bid_size + ask_size)`

This lies roughly in [-1, 1] when sizes are nonnegative and denominator > 0.

A feature is not automatically predictive. Researchers test statistical behavior and robustness.

## 12. Data-structure implications

A book may need:

- fast top-of-book access
- price-level updates
- order lookup by ID
- deterministic memory use
- predictable iteration

Candidate structures include arrays indexed by tick, balanced trees, hash maps + ordered levels, flat maps, and custom intrusive structures. Choice depends on price range, update type, memory, and latency requirements.

## 13. Questions

1. Why is a narrow spread not enough to call a market liquid?
2. What is queue priority?
3. Why might a market maker widen quotes during volatility?
4. Explain adverse selection without jargon.
5. How can stale market data create loss?
6. What is the difference between signal quality and execution quality?
7. Why might an order-book implementation avoid heap allocation per message?

## Lab goal

Build `labs/order_book.cpp`, then benchmark alternative level containers. Defend the chosen design using workload assumptions, not ideology.
