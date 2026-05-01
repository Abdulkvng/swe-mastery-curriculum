# Problem 2 — Best Time to Buy and Sell Stock

> **Pattern:** Scan once, track running minimum
> **Difficulty:** Easy
> **Why it's #2:** Many array problems aren't about a special DS — they're about realizing you only need to remember ONE thing as you scan. This is the simplest example of that.

---

## Problem

You're given an array `prices` where `prices[i]` is the price of a given stock on day `i`.

You want to maximize profit by choosing a single day to buy and a different day in the future to sell.

Return the maximum profit you can achieve. If you cannot achieve any profit, return `0`.

### Examples

```
prices = [7, 1, 5, 3, 6, 4]   →  5    (buy day 1 at 1, sell day 4 at 6)
prices = [7, 6, 4, 3, 1]      →  0    (only decreasing — no profit possible)
prices = [2, 4, 1]            →  2    (buy day 0 at 2, sell day 1 at 4)
```

### Constraints

- `1 <= len(prices) <= 10^5`
- `0 <= prices[i] <= 10^4`

---

## Think before coding

1. What's the brute force? (For each day i, look at every later day j and compute profit. What's the complexity?)
2. As you walk through the array left-to-right, what's the **only thing you need to remember** about everything you've seen so far?
3. Why do you only need ONE variable to track that thing — not a list of all past prices?

---

## Hints

<details>
<summary>Hint 1 — what's the question really asking?</summary>

The question is: **for each day, what's the best profit if I sold today?** And the answer is `today's price - cheapest price seen so far`.

The total answer is the max of those daily-best-profits.

</details>

<details>
<summary>Hint 2 — collapse the state</summary>

You don't need to remember every past price. You only need the **minimum** price seen so far. Why? Because if you're selling today, you should have bought at the cheapest possible past price.

</details>

<details>
<summary>Hint 3 — the loop structure</summary>

```
min_so_far = prices[0]
best_profit = 0

for each price after the first:
    update best_profit with (price - min_so_far)
    update min_so_far with min(min_so_far, price)
```

Order of those two updates matters. Think about why.

</details>

<details>
<summary>Hint 4 — the order trap</summary>

If you update `min_so_far` *before* computing profit, you might accidentally "buy and sell on the same day," which the problem forbids. Compute profit first, then update min.

</details>

---

## Your solution

Write your solution in `02-best-time-stock.py`.

1. State time and space complexity at the top.
2. Trace through `[7, 1, 5, 3, 6, 4]` step-by-step on paper before checking solutions.

---

## Follow-up questions

1. What if you could complete as many transactions as you like (buy, sell, buy, sell...)?
2. What if you could only complete at most TWO transactions?
3. What if there's a transaction fee per trade?

(These are real LeetCode follow-ups — same problem family, increasing in difficulty. We'll get to the multi-transaction versions in the DP module.)
