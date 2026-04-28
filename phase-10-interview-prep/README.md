# Phase 10 — Interview Prep

> By now you have the technical skills. This phase is about conversion — turning what you know into offers. Apple's bar is precise and culture-fit-heavy. Datadog's day-to-day asks specific things. We'll cover Apple's behavioral framework (different bar than IBM's FDE), Datadog's typical scenarios, and the soft skills that separate "competent IC" from "team-shaping engineer."

**Time:** 3–5 days, then ongoing.

**You'll know you're done when:** you have ~20 STAR stories ready, you've practiced 3+ system designs out loud, and you've role-played coding rounds where you talked the entire time.

---

## Table of contents

1. [Apple — what they actually look for](#apple)
2. [Apple behavioral STAR — your stories](#stories)
3. [Apple coding round playbook](#apple-coding)
4. [Apple system design playbook](#apple-system-design)
5. [Datadog day-to-day simulation](#datadog-daily)
6. [Code review etiquette](#code-review)
7. [On-call & runbooks](#oncall)
8. [Blameless post-mortems](#postmortems)
9. [Communication patterns](#communication)
10. [The 90-day plan](#90-day)
11. [Common pitfalls & recovery](#pitfalls)

---

<a name="apple"></a>
## Apple — what they actually look for

Apple's interview process for new grads:
1. **Recruiter screen** — culture / motivation / role fit.
2. **Hiring manager screen** — technical sketch + project deep dive.
3. **Onsite (4–6 rounds):**
   - 2-3 coding rounds (algorithms, problem solving)
   - 1 system design round (often lighter for new grads, but real)
   - 1-2 behavioral / project deep-dive rounds
   - Sometimes a domain-specific round (kernel, ML, SoC depending on team)

What Apple values, in order:
1. **Substance over style.** Quiet competence is preferred over loud ambition.
2. **Ownership.** "I caused/fixed/shipped/owned this" — direct stories, not "we" stories.
3. **Detail orientation.** Apple cares about the corners of a problem.
4. **Curiosity.** Why something works, not just that it does.
5. **Quality.** Pragmatic but with high taste.
6. **Collaboration.** Fierce direct conversation; respectful disagreement.

Apple's interviewers tend to be:
- Concrete (concrete questions, concrete examples)
- A little understated
- Allergic to bullshit and buzzwords
- Generous when you say "I don't know" and reason out loud

What Apple is NOT looking for at new grad:
- Trivia. They don't ask "what year was Swift open-sourced?"
- Buzzword soup. Don't say "leverage synergies." Don't talk like a LinkedIn post.
- Hyperscale flexing. They're not Google. Don't pretend you've designed for 1B users when you haven't.

> 🎯 **Apple vs IBM bar:** IBM FDE prep emphasized client conversation, executive presence, and a specific deployment narrative. Apple is more *engineering-internal*: how you think when handed a hard technical problem, how you push back on a teammate, how you decompose ambiguity. Same STAR framework, different stories.

---

<a name="stories"></a>
## Apple behavioral STAR — your stories

The framework:
- **S** — Situation (1 sentence — context)
- **T** — Task (1 sentence — what was your responsibility?)
- **A** — Action (the meat — 60% of the story; specific, technical, in first person)
- **R** — Result (measurable when possible; what did you learn?)

### The 12 stories you should have ready

Map your real experiences (StackSense, PwC ML, Capital One, ColorStack, GenVote, Browser Use hackathon) to these prompts. Write them out longhand. Practice telling each in 2-3 minutes.

1. **A time you took ownership of an ambiguous problem.**
2. **A time you disagreed with a manager / senior engineer and how you handled it.**
3. **A time you shipped something with quality you're proud of.**
4. **A technical mistake you made and how you recovered.**
5. **A time you helped a teammate.**
6. **A time you had to learn something completely new, fast.**
7. **A time you pushed back on scope or requirements.**
8. **A time you had to make a decision with incomplete information.**
9. **A time you had to advocate for a technical direction.**
10. **A time you debugged something nasty.**
11. **A time you made code/process better than you found it.**
12. **A time you led without authority.**

### Story-writing template

For each, fill out:
```
PROMPT:
SITUATION (1-2 sentences):
TASK / your specific stake (1 sentence):
ACTION:
  - I diagnosed by ___
  - I tried ___, which failed because ___
  - I changed approach to ___
  - The decision came down to ___ vs ___; I chose ___ because ___
  - I built/wrote/tested ___
RESULT:
  - measurable: ___
  - what I'd do differently: ___
  - what I learned: ___
```

### Example: "A time you debugged something nasty"

> *(Adapt to your actual experience. Sample skeleton:)*
>
> **S:** "On StackSense, the AI gateway was returning 502s about 1% of the time, but only in production."
>
> **T:** "I owned the gateway service end-to-end and was responsible for the SLO."
>
> **A:** "I started by checking logs — they showed normal completion times. I added structured request IDs and noticed the failures all had longer-than-normal latency just before failing. I suspected timeouts. I then traced through the load balancer config and found the LB had a 30-second idle timeout, but our LLM upstream sometimes exceeded that for long completions. The fix was bumping the LB idle timeout AND implementing a 'streaming heartbeat' that flushed bytes during long runs so the connection stayed live. I also added a span for 'time-to-first-byte' to our OpenTelemetry traces so future regressions would be easier to spot."
>
> **R:** "Error rate dropped from 1% to <0.05% over a week. The bigger lesson for me was that 'normal logs' don't mean normal — I should have looked at latency distributions earlier."

That's the shape. Substance, detail, ownership, learning.

### Apple-specific tips

- **First person, always.** "I" not "we." Apple wants to know what *you* did.
- **Acknowledge other people generously** when they helped, but never hide behind them.
- **Resist self-deprecation.** Confident reflection > "oh, it was nothing."
- **Volunteer the failure.** Stories with no struggle ring false. Apple knows real engineering has failures.

---

<a name="apple-coding"></a>
## Apple coding round playbook

### Before they ask the question

- Have a clean coding setup ready. Default to Python or Go for whiteboard rounds. Apple often uses CoderPad / similar.
- Know your IDE shortcuts. Spending 30 seconds finding `delete line` looks bad.
- Test your audio/screen-share an hour before. Tech failures are interviewer-side too but you don't get the time back.

### When they ask the question (the 4-step framework, applied)

1. **Restate.** "So I have an array of integers, possibly with duplicates, and I want to return all unique pairs that sum to a target — making sure I don't return the same pair twice. Do I need to return them in any particular order?"

2. **Walk through examples.** Pick one yourself if they don't give you one. Trace through what the answer should be. Catches misunderstandings before you've coded.

3. **Discuss approaches.**
   - Brute force first: "the obvious O(n²) is to nest two loops..."
   - Then optimize: "we can use a hash set seen so far for O(n)..."
   - State complexity: "O(n) time, O(n) space; the trade-off vs sorting + two pointers is..."
   - Ask: "Should I code this approach?"

4. **Code.** Talk through what each line does as you write. Use clear names. Don't be silent.

5. **Test.** Walk through your example. Then edge cases:
   - Empty input
   - Single element
   - All duplicates
   - Negatives, zero
   - The target itself appearing in the array

6. **Reflect.** "If we needed to handle a streaming version of this..." "If memory were tight..."

### What Apple interviewers really evaluate

- **Communication.** Were you clear about what you were doing?
- **Problem solving.** Did you decompose? Did you handle stuck-ness?
- **Code quality.** Sensible names, no spaghetti, basic correctness on the first pass.
- **Curiosity.** Did you ask good questions? Did you push the problem yourself?
- **Edge cases.** Did you consider them without prompting?

You'll get partial credit. A clear brute force you can defend > a brilliant optimal solution you can't explain.

---

<a name="apple-system-design"></a>
## Apple system design playbook

For new grads, often "lighter" — a 30-45 min sketch instead of 60+ min. But take it seriously.

### Common new-grad prompts at Apple

- Design a notification service.
- Design a URL shortener.
- Design a key-value store.
- Design a photo upload service.
- Design a rate limiter.
- Design an autocomplete service.

### Apple-specific flavor

- **Be opinionated.** Apple respects taste. Pick a database and defend, don't list 5.
- **Concretize.** "It's roughly 1B writes/day, 100KB per object, retained 30 days = ~3 PB" beats "it's a lot of data."
- **Care about quality.** Apple asks more "how do you ensure correctness" questions than "how do you scale to a trillion."
- **Failure modes.** Always volunteer 2-3 failure modes and how you'd handle them.

Use the 7-step framework from Phase 7. Just go faster.

---

<a name="datadog-daily"></a>
## Datadog day-to-day simulation

What does ADP-Notebooks day-to-day actually look like? Predicting:

### Morning
- Slack: any overnight alerts? Check the on-call channel.
- Standup (10-15 min). Async or sync.
- Triage: any incoming bug reports? Customer-reported issues?
- Plan the day in your head: 2-3 PRs to write/review, 1-2 deeper tasks.

### Mid-day
- IDE work. Most likely:
  - TypeScript/React for UI work
  - Go for backend services
  - YAML for k8s manifests
  - Terraform for infra (maybe)
- Code review for teammates.
- Periodic context-switch to Slack threads.

### Afternoon
- Pair-debug something with a coworker.
- Architecture discussion in a small meeting.
- Investigation: dig into a customer's slow notebook, file under "performance regression."
- Push a PR, wait for CI (~10 min), iterate.

### Mental model

You're not "writing features all day." You're a thinking-and-debugging engine that occasionally produces code. The debugging and design work is what you're paid for; code is the byproduct.

### What gets junior engineers stuck

- **Not asking for help fast enough.** If you're stuck for 30+ min, ask. Cost of asking < cost of staying stuck.
- **Spending too long on "perfect" PRs.** Push small incremental changes; don't drop a 2000-line PR on a Friday afternoon.
- **Avoiding hard meetings.** When in doubt, attend the design review. Listen.
- **Not reading code.** Senior engineers spend half their time reading code. So should you.

---

<a name="code-review"></a>
## Code review etiquette

You'll do this every day at Datadog.

### As reviewer

**Be specific.** "This loop is O(n²); could it be O(n) with a map?" beats "this is slow."

**Ask, don't decree.** "What if the input is empty here?" beats "you didn't handle empty input."

**Distinguish blocking from non-blocking.** Prefix with **nit:** for non-blocking nits. Save "blocking" for actual blockers.

**Approve when you'd ship it.** Don't gate on aesthetic preferences.

**Prioritize.** A 200-line PR review with 50 nits and 1 critical bug — make sure the critical bug is unmistakable.

**Tone matters.** Written communication has no tone, so default to friendly. Use "Could we..." instead of "Why did you...". Add a 🎉 when something is great.

### As author

**Small PRs.** Aim for <300 lines changed when possible. Reviewers' attention is finite.

**Title and description.** Title says *what*. Description says *why*. Include screenshots/snippets for UI changes.

**Self-review first.** Open your own PR; read every line. You'll catch half your bugs before any human reviewer sees them.

**Respond, don't just push.** If a reviewer asks a question, answer it (in the thread) AND push a fix. Don't just push a fix and leave the question hanging.

**Push back politely.** If you disagree, say why with reasoning. "I considered X but went with Y because Z. Open to other approaches if Z isn't compelling."

**Mark threads resolved as you address them.** Helps the reviewer see what's left.

---

<a name="oncall"></a>
## On-call & runbooks

Most Datadog engineers are on a primary on-call rotation. ADP is no exception.

### What you do on-call

- **Acknowledge alerts** — within ~5 min if it's a paging alert.
- **Triage:** is this real, severity, scope.
- **Mitigate** — restore service first, root-cause later.
- **Communicate** — start an incident channel, say what you know, update at intervals.
- **Hand off cleanly** — when your shift ends, brief the next person.

### Runbooks

When you build a service, write a runbook for it. A runbook is "if alert X fires, here's how to investigate and fix":

```markdown
# Notebook API service runbook

## Alert: notebook_api_p99_latency_high

### What it means
p99 of `/api/notebooks/*` exceeded 500ms for 5+ minutes.

### Triage
1. Check `notebook_api` deployment in <env>
2. Check Postgres connection saturation: <link to dashboard>
3. Check Redis cache hit rate: <link to dashboard>

### Common causes
- Postgres slow query (check pg_stat_statements)
- Cache cold (recent deploy, restart)
- Downstream kernel orchestrator lagging

### Mitigation
- Bump replica count: `kubectl scale deployment/notebook-api --replicas=N`
- If DB-related: check long-running queries: `SELECT * FROM pg_stat_activity WHERE state != 'idle' ORDER BY query_start;`

### Escalation
- @notebook-team in Slack
- If after-hours and severity 2+: page secondary on-call
```

The point is: future-you (or future-someone-else) at 3am should not have to rediscover what to do.

---

<a name="postmortems"></a>
## Blameless post-mortems

After every significant incident, write a post-mortem.

### Structure

```markdown
# Incident: <title>

## Summary
1-2 paragraphs. What happened, who was affected, when, duration.

## Timeline
- 14:32 UTC — alert fires
- 14:34 UTC — on-call ack'd
- 14:41 UTC — root cause hypothesized
- 14:48 UTC — fix deployed
- 14:52 UTC — service recovered

## Root cause
What was the actual underlying issue? Be technical and specific.

## What went well
- Alerts fired promptly
- On-call response < 2 min
- Mitigation was clear from runbook

## What went wrong
- Post-deploy canary did not catch the regression because [...]
- The dashboard linked from the alert was missing the relevant query

## Action items
- [ ] @kvng: add canary check for X
- [ ] @teammate: update runbook
- [ ] @sre: alert tuning
```

### Blameless

The "blameless" part: focus on **systems and processes**, not individuals.

- Bad: "Alice deployed a bug."
- Good: "The deploy pipeline allowed a change to ship that didn't pass the integration suite because of a misconfigured CI flag."

Engineers do their best work when mistakes don't lead to fear. The Datadog post-mortem culture is genuinely good — read [Datadog incident reports](https://www.datadoghq.com/blog/) for examples.

---

<a name="communication"></a>
## Communication patterns

### Status updates

Once a week, write a brief async update:
- What I shipped
- What I'm working on
- What's blocking me
- Where I need help

### Saying "I don't know"

Strong engineers say it constantly. Weak ones bluff. The strong-engineer way:
> "I don't know off the top of my head. My guess is X based on Y — let me verify."

### Disagreeing

Three steps:
1. Restate their position to confirm you understand.
2. State your concern with reasoning.
3. Ask: "Am I missing something?"

If after that you still disagree, you can:
- Defer (their call, you're not the decider) — say so explicitly.
- Escalate to a higher-up — only when the stakes warrant.
- "Disagree and commit" — you've voiced it, decision is made, you support execution.

### Asking for help

Bad: "It's not working."
Good: "I'm trying to do X. I expected Y. I'm getting Z. I've tried A and B without luck. I'm stuck on the question of whether C is even possible. Any pointers?"

The good version means whoever helps spends 5 min, not 30.

---

<a name="90-day"></a>
## The 90-day plan

For Datadog Summer 2026 — what you should aim for.

### Days 1–30 — onboarding

- Get the dev environment running end-to-end.
- Read the team's main service code top-to-bottom (yes, all of it).
- Make 2–3 small PRs (typo fixes, doc improvements, small features).
- Understand who owns what. Map the team and adjacent teams.
- Set up your own observability for things you care about.
- Pair-program with each teammate at least once.

### Days 31–60 — contributing

- Own a small project end-to-end. Spec it, design it, ship it, instrument it.
- Be on-call (shadow first, then primary).
- Review others' code with substance.
- Start surfacing your opinions in design reviews.
- Write at least one substantial doc (RFC, runbook, post-mortem).

### Days 61–90 — finishing strong

- Land your project. Measure its impact.
- Have a clear mid-term update with your manager.
- Identify one thing you'd love to push further if it became a return offer.
- Leave the codebase better than you found it.

### Return offer signals

You convert at Datadog by being:
- Reliable (ship what you say, on time, with quality)
- Self-directed (don't need constant hand-holding)
- A good teammate (helpful in code review, supportive in Slack)
- Thoughtful (you have opinions and can defend them)
- Visible (your work is known, but not in a self-promotional way)

Show up. Do good work. Tell your manager what you're doing. The offer follows.

---

<a name="pitfalls"></a>
## Common pitfalls & recovery

### "I bombed a coding round."

It happens. Recovery:
- Do NOT bring it up unprompted in later rounds.
- If asked "how did the previous round go?" — be honest, brief: "I struggled with the optimization on problem 2 but I think I had a clean brute force."
- Apple often recovers people from one weak round. Hiring is a vote.

### "I went blank on a behavioral question."

- Take a breath. "Let me think for a second."
- Pick a story that's even *partially* relevant.
- "It's not a perfect match, but here's a similar situation I had..."

### "I disagreed with the interviewer."

If you're confident: "I see your point — I was thinking about it differently because of X. Could you walk me through why Y matters more here?"

If you're unsure: defer gracefully. "Fair. Let me update my approach."

Never argue from authority ("but the textbook says..."). Argue from reasoning.

### "I made a math/code error and noticed mid-explanation."

Catch it explicitly. "Wait, I just realized — my analysis above is wrong because [...]. Let me redo this part." Interviewers love this. Self-correction = the strongest signal of seniority.

### Final piece of advice

Your interviews are not a test of who you are. They're a snapshot of how you happen to perform in 4-6 hours of artificial conversation. You'll do well in some, worse in others. Both Apple and Datadog calibrate for that. **Just keep showing up as a thoughtful, curious engineer who's trying to solve the actual problem in front of you.** That works.

---

**Next:** [Phase 11 — ML/AI Engineering](../phase-11-ml-ai/README.md)
