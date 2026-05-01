# Module 1 — Arrays & Strings: Interview Simulation

> **Setup:** Imagine you're 25 minutes into a 45-minute coding round. The interviewer pastes the problem below into a shared editor.
>
> **Rules of engagement:**
> 1. Don't start coding immediately. Spend at least 60 seconds restating and asking clarifying questions.
> 2. Talk through your approach BEFORE coding. State complexity.
> 3. Code while narrating.
> 4. Test with at least one example by tracing through.
> 5. Then read the "interviewer follow-ups" section and answer each.

---

## The problem (as the interviewer would state it)

> "Given a string `s`, find the length of the longest substring that contains at most two distinct characters."

That's it. That's all the interviewer says. They don't define "substring," they don't give examples, they don't say what to return for the empty string.

**Stop. Before reading further, list at least 4 clarifying questions you would ask.**

---

## Clarifying questions you should ask

<details>
<summary>Reveal the questions a strong candidate would ask</summary>

1. "By **substring**, do you mean contiguous? Or are you using it loosely to mean subsequence?" → contiguous.
2. "What characters can appear? ASCII? Unicode? Just lowercase letters?" → assume ASCII; the answer doesn't depend on the alphabet size for correctness, but mentioning it shows awareness.
3. "What should I return for the empty string?" → 0.
4. "What if the entire string has only one or two distinct characters? I should just return its length, right?" → yes.
5. "Do I need to return the substring itself, or just its length?" → just the length.
6. "Are there any size constraints? Could `s` be enormous? Streaming?" → assume up to 10⁵ characters, fits in memory.

A candidate who jumps straight to coding without asking ANY of these is signaling they aren't careful. Asking 2–3 of these is plenty.

</details>

---

## Your turn

Now solve it. Take 15–20 minutes. Treat this like a real round: write your solution in `interview-sim-solution.py` (you'll create this file), narrate as you go (you can mutter to yourself or write a comment-narration if working solo), then trace an example.

Examples to test against:

```
"eceba"        →  3   ("ece")
"ccaabbb"      →  5   ("aabbb")
"a"            →  1
""             →  0
"abcabcabc"    →  2   (any two adjacent)
```

---

## Interviewer follow-ups (don't read until you've solved it)

<details>
<summary>Follow-up 1 — generalize to K</summary>

"Now generalize. What if I asked for at most K distinct characters?"

Your answer should be: "Same algorithm, just replace `2` with `K`. The hash map tracks at most K+1 entries at any time. Time stays O(n), space becomes O(K)."

If you wrote the original cleanly, this is a 30-second change.

</details>

<details>
<summary>Follow-up 2 — what's the complexity of building the answer string?</summary>

"What if I asked for the actual substring, not just the length?"

Your answer: "I'd track `best_start` and `best_length` as I go. At the end, return `s[best_start:best_start + best_length]`. Time is still O(n), space adds an O(K) substring at the end."

</details>

<details>
<summary>Follow-up 3 — streaming version</summary>

"What if `s` is a stream — characters arriving one at a time over hours, and at any moment I need to answer 'what's the longest valid suffix ending at the current position?'"

Your answer: "The current window I'm maintaining IS that answer. So as long as I keep `left`, `right`, and the hash map updated as each character arrives, I can answer in O(1) at any time. Total ongoing work is O(n) over the entire stream."

This is a great moment to mention: "This is essentially how monitoring tools handle 'last N seconds of activity' — sliding window is everywhere in observability." (Bonus points at Datadog.)

</details>

<details>
<summary>Follow-up 4 — what if K is large?</summary>

"What if K could be up to a million?"

Your answer: "Hash map size grows linearly with K, so memory becomes O(K). Time per character is still O(1) amortized — we're just inserting/deleting/looking up in a hash map. The algorithm doesn't degrade."

If they push: "Could the inner `while` loop ever do a lot of work for a single outer step?" — the answer is yes per-step, but each character can only leave the window once across the entire scan. Amortized O(1) per character. This is the classic sliding window argument.

</details>

<details>
<summary>Follow-up 5 — adversarial input</summary>

"Construct an input where your algorithm does the maximum amount of total work."

Something like `"abababab..."` — for every position, we add a char then immediately have to shrink. But notice: the total work is still O(n) because each character is added and removed at most once. There's no input that breaks the linear bound.

</details>

---

## Self-evaluation rubric

Score yourself honestly:

| Criterion | Did you...? | Yes/No |
|---|---|---|
| Restated the problem in your own words | | |
| Asked at least 2 clarifying questions before coding | | |
| Stated complexity BEFORE writing code | | |
| Wrote clean, named-variable code (not `i, j, k, x`) | | |
| Tested with at least one example by tracing | | |
| Caught your own bug (if any) before claiming "done" | | |
| Answered follow-up 1 (generalize to K) without re-coding | | |
| Mentioned the amortized argument | | |

5+ yes: solid. 7+ yes: this is a hire. Anything below 5: re-read INTUITION.md and redo.
