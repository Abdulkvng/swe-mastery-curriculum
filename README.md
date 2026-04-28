# SWE Mastery Curriculum

> A self-paced, project-driven curriculum to master everything needed for a Datadog ADP Notebooks internship and an Apple New Grad SWE interview. Built for Kvng.

## Who this is for

You. Specifically:
- USC junior, double major in CS + Business, AI minor
- Incoming Datadog SWE intern (Summer 2026), Advanced Data Platform - Notebooks team
- Targeting Apple New Grad SWE
- Already has ML internship experience (PwC), SWE internship experience (Capital One)
- Builder of StackSense (AI gateway/observability platform)

## How to use this

Read top-to-bottom. Each phase **builds on the last** — you can't skip ahead without missing the foundation. Every concept has runnable code. Every keyword you might not know is defined inline the first time it appears, like this:

> 📖 **Definition — Keyword:** Plain-English explanation here.

When you see a 🛠️ icon, that's a project. Build it. Don't just read it.

When you see a 🧠 icon, that's a "what does this even mean?" deep-dive on a concept that gets thrown around but rarely explained.

When you see a 🎯 icon, that's an interview-relevance flag — pay extra attention.

## The 13 Phases

| # | Phase | What you'll master | Languages |
|---|---|---|---|
| 0 | Foundations & Environment | Linux, bash, Git, GitHub, GitLab, CI/CD basics | Bash |
| 1 | Networking & Protocols | TCP/IP, HTTP, HTTPS, TLS handshake, encryption | Rust |
| 2 | OOP & Design Patterns in Go | OOP pillars, singleton, factory, **connection pool** | Go |
| 3 | Data Structures & Algorithms | Arrays → graphs → DP, ~80 problems mapped | Go, Rust |
| 4 | Databases Deep Dive | ACID, isolation, indexes, Postgres, Redis, pooling | Go, SQL |
| 5 | APIs & Backend | REST, gRPC, protobuf, JWT, rate limiting | TypeScript |
| 6 | Concurrency & OS | Threads, processes, parallelism, mutexes, embedded basics | Go, Rust |
| 7 | Distributed Systems & System Design | Framework, Kafka, sharding, consensus, CAP | Go, conceptual |
| 8 | Datadog Stack | Spark, Scala, Jupyter internals, Kubernetes, Docker | Scala, Go |
| 9 | Capstones | Mini-ADP-Notebooks + Mini-Datadog | Polyglot |
| 10 | Interview Prep | Apple behavioral, Datadog day-to-day, system design rubric | — |
| 11 | ML/AI Engineering | MLOps, model serving, vector DBs, LLM serving | Python |
| 12 | Frontend Depth | React internals, state mgmt, perf, a11y, build tools | TypeScript |

## Repository structure

```
swe-mastery-curriculum/
├── README.md                          # this file
├── phase-00-foundations/
│   ├── README.md                      # the chapter
│   ├── projects/                      # runnable projects
│   └── exercises/                     # practice problems
├── phase-01-networking/
│   ├── README.md
│   ├── projects/
│   │   └── tcp-to-tls-server/        # Rust project
│   └── exercises/
├── ... (one folder per phase)
└── assets/
    └── diagrams/                      # ASCII + SVG diagrams
```

## Conventions

- **Mac (Apple Silicon, M3 Pro)** is the assumed dev environment. Most commands work on Linux too.
- Every project has its own README with: setup, run, test, what you should learn.
- Every code file has comments explaining *why*, not just *what*.
- "Build upon" means: Phase N projects literally import or reference Phase N-1 code where possible.

## Daily ritual (suggested)

1. Read 30–60 min of curriculum text.
2. Code 60–90 min on the current project.
3. Solve 1–2 LeetCode problems (mapped per phase).
4. Skim Datadog/Apple engineering blog or one paper from `phase-07-distributed-systems/papers.md`.

## Glossary

A live glossary lives at `GLOSSARY.md`. Every defined term ends up there. When in doubt, ctrl-F.

## Status

- [x] Phase 0 — Foundations & Environment
- [x] Phase 1 — Networking & Protocols
- [x] Phase 2 — OOP & Design Patterns in Go
- [x] Phase 3 — Data Structures & Algorithms
- [x] Phase 4 — Databases
- [x] Phase 5 — APIs & Backend
- [x] Phase 6 — Concurrency & OS
- [x] Phase 7 — Distributed Systems
- [x] Phase 8 — Datadog Stack
- [x] Phase 9 — Capstones
- [x] Phase 10 — Interview Prep
- [x] Phase 11 — ML/AI Engineering
- [x] Phase 12 — Frontend Depth

Built across multiple Claude turns. Each turn adds 1–2 phases.

---

*"The best way to predict the future is to build it."* — let's go.
