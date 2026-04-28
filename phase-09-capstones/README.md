# Phase 9 — Capstones

> Two final projects that combine everything from phases 0–8. The first simulates your day-to-day at Datadog ADP-Notebooks. The second simulates the work behind the Datadog product itself. Build them and you'll have repo links to point to in interviews and conversations.

**Time:** 4–8 weeks. Work on these in parallel with later phases.

---

## Capstone A — Mini-ADP-Notebooks

A simplified version of what you'll build at Datadog. Stack:

- **Frontend:** React + TypeScript + Monaco editor (the VSCode editor, embeddable)
- **Backend:** Go server, exposes REST + WebSocket
- **Kernel pool:** PySpark kernels in Docker containers, orchestrated by k3d
- **Storage:** Postgres for notebook metadata + cells; S3-compatible (MinIO locally) for outputs
- **Auth:** JWT (from Phase 5)
- **Observability:** OTel traces + Prometheus metrics + structured logs
- **CI:** GitLab CI pipeline (because Datadog uses GitLab)

### Spec

#### User stories
1. A user logs in (JWT-issued).
2. They create a new notebook, given a title.
3. They write Python code in cells, click run; the cell output appears.
4. They can run multiple cells; previous cell variables persist (kernel state).
5. They can read/write a small CSV in S3 and process it via PySpark.
6. They can share a notebook (read-only) with another user via a link.
7. The kernel auto-shuts-down after 30 min idle.
8. Notebooks save automatically every 30 sec.

#### Architecture

```
[Browser]
   │ WebSocket (cell exec, output stream)
   │ HTTPS REST (CRUD on notebooks)
   ▼
[API Gateway / nginx]
   ▼
┌──────────┐         ┌──────────┐
│  API svc │────────►│  Postgres│
│  (Go)    │         └──────────┘
│          │         ┌──────────┐
│          │────────►│  Redis   │  (sessions, kernel routing)
│          │         └──────────┘
│          │         ┌──────────┐
│          │────────►│  MinIO   │  (S3-compatible)
└────┬─────┘         └──────────┘
     │ k8s API
     ▼
[Kernel Operator]
     │ creates pods
     ▼
[Kernel Pod 1] [Kernel Pod 2] ...
  (PySpark)     (PySpark)
```

#### Repo layout

```
mini-adp-notebooks/
├── README.md
├── docker-compose.yml         <- local stack
├── .gitlab-ci.yml             <- CI pipeline
├── frontend/                  <- React + Vite + TS
├── api/                       <- Go backend
├── kernel-image/              <- Dockerfile for PySpark kernel
├── operator/                  <- k8s controller (Go) for kernel lifecycle
├── manifests/                 <- k3d manifests
└── docs/
    ├── architecture.md
    └── runbook.md
```

### Why this is the perfect interview project

When an interviewer asks "tell me about a recent project," this is what you show. You can talk about:
- Architectural decisions (why Go over Node? why kernel-per-pod?)
- Trade-offs (kernel cold-start latency vs. always-warm cost)
- Failure modes (what if a kernel crashes mid-execution? what if Postgres goes down?)
- Observability (how would you know if kernels are slow?)
- What you'd do differently next time (the pre-mortem, in retrospect)

**This is exactly the conversation you want to have at Datadog.**

### Build order suggestion

Do this in 5 milestones over a few weeks:

1. **Week 1: Backend skeleton.** Go API, Postgres + migrations, JWT auth, CRUD on notebooks/cells. No kernels yet. Local dev with Docker Compose.
2. **Week 2: Kernel POC.** Run a PySpark process inside a Docker container manually. Wire up WebSocket → backend → kernel via Jupyter wire protocol. Echo "execute" to "result" in the simplest case.
3. **Week 3: Frontend.** React app, Monaco editor, WebSocket connection to backend, render outputs.
4. **Week 4: Kubernetes orchestration.** Move kernel from Docker Compose to k3d. Write a tiny Go operator that, given a CR like `Notebook`, creates a kernel Pod + Service.
5. **Week 5: Polish.** OTel tracing, Prometheus metrics, structured logs, CI pipeline, runbook.md, architecture.md.

You don't have to finish to use the project. Even milestone 2 is impressive.

---

## Capstone B — Mini-Datadog

A simplified version of Datadog's product itself: ingest, store, query metrics + traces. Demonstrates you understand observability *infrastructure*, not just observability *culture*.

### Spec

- **Metric ingestion API**: POST `/v1/metrics` with `{ metric, value, tags, timestamp }`.
- **Storage**: time-bucketed columnar files in MinIO; in-memory write buffer flushed every 60s. Inverted index for tag lookups in Redis.
- **Query API**: `GET /v1/query?metric=...&from=...&to=...&group_by=...&agg=avg|p99|sum`.
- **Dashboard UI**: simple React app with line charts (recharts), can pick metric + tag filter + time range.
- **Synthetic load generator**: a Go program that fires fake metrics into the ingestion API at configurable rate.

### Stack

- **Ingest:** Go service, writes to Kafka (single-broker for dev).
- **Compactor:** Go worker, drains Kafka into Parquet files in MinIO, partitioned by hour.
- **Index:** Redis sorted sets per tag value → series IDs.
- **Query:** Go service, Parquet reader (apache/arrow-go), aggregations.
- **UI:** React + Recharts.

### Why this matters

You'll be working *adjacent* to this kind of system at Datadog. Even if the ADP-Notebooks team doesn't write the metrics storage, you'll integrate with it. Building a toy version means you understand the constraints, the access patterns, the failure modes.

### Repo layout

```
mini-datadog/
├── README.md
├── docker-compose.yml
├── ingest/             <- Go ingest service
├── compactor/          <- Go compactor worker
├── query/              <- Go query service
├── ui/                 <- React dashboard
├── loadgen/            <- Go synthetic load
└── manifests/          <- k3d manifests
```

---

## Smaller capstones (pick one if the big ones are too ambitious)

If 4-8 weeks is too much, smaller variants:

- **Mini-pgbouncer** (from Phase 4) — finish the connection-pooling proxy. Great senior interviews come up around this.
- **Distributed rate limiter** — token-bucket, Redis-backed, gRPC + REST APIs, deployed to k3d. Touch every concept.
- **Event-driven order service** — Kafka producer + consumer + Postgres + outbox pattern. Demonstrates understanding of common backend patterns.

---

## What to do with these projects

1. Push to a public GitHub repo.
2. Write a great README — explain decisions, draw the architecture, mention trade-offs.
3. Add a "What I learned / what I'd do differently" section. Senior interviewers read this carefully.
4. Link from your resume and LinkedIn.
5. When asked "tell me about a project," walk through architecture → an interesting trade-off → a failure mode → what you'd do at scale.
