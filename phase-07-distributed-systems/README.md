# Phase 7 — Distributed Systems & System Design

> Single-machine engineering ends here. Distributed systems is what happens when "just add another server" becomes an architectural decision instead of a vibe. CAP, consensus, sharding, replication, message queues, eventual consistency — all the words that get thrown around in system design interviews. This phase makes them real.
>
> System design interviews are a *language* — there's a vocabulary, a framework, and patterns. We'll cover the framework, then walk through 6 canonical designs (URL shortener, Twitter timeline, rate limiter, distributed cache, time-series database, Jupyter-notebook backend) so you have analogues for any prompt.

**Time:** 7–10 days, plus continuous practice.

**You'll know you're done when:** you have a 7-step framework you can run on autopilot in any system-design interview, you can pick the right database/queue/cache for a workload and defend it, you understand the CAP/PACELC trade-offs concretely, and you've designed (on paper or whiteboard) at least 6 systems end-to-end.

---

## Table of contents

1. [What does this even mean? — "Distributed system"](#what-distrib-means)
2. [Module 7.1 — Why distributed systems exist (and the fallacies)](#module-71--why)
3. [Module 7.2 — Latency numbers every engineer should know](#module-72--latency)
4. [Module 7.3 — CAP, PACELC, and what they actually mean](#module-73--cap)
5. [Module 7.4 — Replication](#module-74--replication)
6. [Module 7.5 — Sharding & partitioning](#module-75--sharding)
7. [Module 7.6 — Consistent hashing](#module-76--consistent-hash)
8. [Module 7.7 — Consensus: Paxos, Raft (intro)](#module-77--consensus)
9. [Module 7.8 — Message queues: Kafka deep-ish dive](#module-78--kafka)
10. [Module 7.9 — Caching: layers, invalidation, hot keys](#module-79--caching)
11. [Module 7.10 — Load balancers and service discovery](#module-710--lb)
12. [Module 7.11 — Failure modes you'll actually see](#module-711--failures)
13. [Module 7.12 — The 7-step system design framework](#module-712--framework)
14. [Module 7.13 — Six canonical designs](#module-713--designs)
15. [Papers worth reading](#papers)
16. [Interview question bank](#interview-questions)
17. [What you should now know](#what-you-should-now-know)

---

<a name="what-distrib-means"></a>
## 🧠 What does this even mean? — "Distributed system"

> 📖 **Definition — Distributed system:** A set of independent computers that, to the user, appear as one coherent system. They communicate over a network and coordinate to provide a service.

That's the textbook definition. The honest one: **a system where another computer you don't control can fail and break yours.**

Single-process: when something goes wrong, it goes wrong here. Distributed: a thing 5000 miles away you've never met can disappear and your code now needs to know what to do.

Leslie Lamport: *"A distributed system is one in which the failure of a computer you didn't even know existed can render your own computer unusable."*

That's the daily emotional reality. Designing distributed systems = designing for partial failure.

---

<a name="module-71--why"></a>
## Module 7.1 — Why distributed systems exist (and the fallacies)

### Why we go distributed

Three reasons. In rough order of historical importance:

1. **Capacity.** A single machine can't hold/serve all the data. Google, Facebook, Datadog metrics — the data volume forces multi-machine.
2. **Availability.** Hardware fails. If your service can run on N machines and survive N-1 failing, you have meaningful uptime.
3. **Latency.** Users in Tokyo shouldn't wait 200 ms for a packet to round-trip from Virginia. CDNs, edge compute.

Sometimes also: regulatory (data must stay in the EU), org structure (different teams own different services).

### The 8 fallacies of distributed computing (Sun, 1994 — still painfully relevant)

Things engineers wrongly assume about networks. Memorize these:

1. The network is reliable. ❌
2. Latency is zero. ❌
3. Bandwidth is infinite. ❌
4. The network is secure. ❌
5. Topology doesn't change. ❌
6. There is one administrator. ❌
7. Transport cost is zero. ❌
8. The network is homogeneous. ❌

Every distributed system bug, traced deep enough, comes from believing one of these.

---

<a name="module-72--latency"></a>
## Module 7.2 — Latency numbers every engineer should know

Jeff Dean's famous numbers, modernized for 2026:

| Operation | Time | Visualization |
|---|---|---|
| L1 cache reference | 1 ns | a heartbeat |
| Branch mispredict | 3 ns | |
| L2 cache reference | 4 ns | |
| Mutex lock/unlock | 17 ns | |
| Main memory reference | 100 ns | |
| Compress 1 KB with Snappy | 2 µs | |
| Send 2 KB over 10 Gbps network | 4 µs | |
| Read 1 MB sequentially from memory | 10 µs | |
| Round trip in same datacenter | 500 µs | |
| Read 1 MB sequentially from SSD | 1 ms | |
| Disk seek (HDD) | 5 ms | |
| Round trip US west↔east coast | 40 ms | |
| Round trip US↔Europe | 100 ms | |
| Round trip US↔Asia | 200 ms | |

Mental shortcut to memorize:

```
L1: 1ns
L2: 4ns
RAM: 100ns  (100x L1)
Datacenter RTT: 500us  (500,000x L1, 5000x RAM)
SSD seek: 1ms
Cross-country: 40ms
Cross-ocean: 200ms
```

Use these in interviews when sizing things. *"Disk read is ~1 ms; in-memory cache is ~100 ns; that's a 10,000x speedup, which is why we cache."*

---

<a name="module-73--cap"></a>
## Module 7.3 — CAP, PACELC, and what they actually mean

### CAP

> 📖 **Definition — CAP theorem:** In a system that experiences a network partition (some nodes can't reach others), you must choose between **Consistency** (every read returns the most recent write) and **Availability** (every request gets a response). You can't have both. You always have Partition tolerance because, as Coda Hale wrote, "*partition tolerance is not optional*."

Most explanations of CAP are wrong or oversimplified. The honest version:

- **CP system:** during a partition, refuses to serve requests that can't be guaranteed consistent. Examples: Postgres + sync replication, MongoDB (in default mode), Zookeeper, etcd.
- **AP system:** during a partition, keeps serving but might return stale data; resolves conflicts later. Examples: DynamoDB, Cassandra (default), CouchDB.

CAP applies *only during a partition*. Most of the time there's no partition. Which leads to PACELC.

### PACELC — the more useful framing

> 📖 **Definition — PACELC:** **If P**artition, choose **A** or **C**. **E**lse (no partition), choose **L**atency or **C**onsistency.

The "ELC" part is what actually shapes most system design:
- **EL** (latency-prioritizing) systems serve reads from replicas, possibly stale. Fast but eventually consistent.
- **EC** (consistency-prioritizing) systems serve reads from the primary or with strict ordering. Slower but always fresh.

Cassandra: PA / EL. Always available, always low-latency, eventually consistent.
Postgres: PC / EC. Strong consistency normally; refuses to serve writes during partition.

Datadog Husky: PA / EL. Metrics are append-mostly; a 30-second lag in queryability is fine.

### Eventual consistency — what it actually means

> 📖 **Definition — Eventual consistency:** If no new updates are made to an item, eventually all reads will return the last-updated value.

Notice the weak guarantee. In between, reads can return any version. Anomalies you see:
- Read your own writes: You write, immediately read, get OLD value. Workaround: route reads to primary for the user who just wrote.
- Monotonic reads: Read says X. Read again says older-than-X. Workaround: pin reads to the same replica for a session.
- Consistent prefix: Replicas process operations out of order; observers see "reply before question." Workaround: causal consistency / vector clocks.

**The biggest mistake** is assuming eventual = "100ms latency." It often is. But there's no upper bound. During a partition, eventual could mean hours.

---

<a name="module-74--replication"></a>
## Module 7.4 — Replication

We touched this in Phase 4 from the DB perspective. Distributed-systems view:

### Single-leader (primary/replica)

One leader accepts writes; replicates to followers. Followers serve reads (or stay hot for failover).

- ✅ Simple model. ACID-ish.
- ✅ Read scaling.
- ❌ Single write point. Throughput capped by leader.
- ❌ Failover is dangerous. Two leaders if not careful.

Examples: Postgres streaming, MySQL primary/replica, Redis primary/replica.

### Multi-leader

Multiple nodes accept writes; replicate to each other. Conflicts resolved via "last-write wins" (clock-based, prone to clock-skew bugs) or CRDTs (conflict-free replicated data types).

- ✅ Geo-distributed writes.
- ❌ Conflict resolution is hard.
- ❌ Most use cases don't actually need this.

Examples: CouchDB, BDR for Postgres, multi-region active-active in DynamoDB Global Tables.

### Leaderless

Every replica accepts reads and writes. Quorums for consistency: read R replicas, write W replicas; if R+W>N, you can't read stale.

- ✅ Highly available, no failover.
- ❌ Tunable consistency = decisions to make.

Example: Cassandra. DynamoDB.

### Synchronous vs asynchronous replication

- **Sync:** writes wait for replica ack. Stronger durability, higher latency.
- **Async:** writes commit immediately, replicas catch up. Faster, may lose data on leader crash.
- **Semi-sync:** wait for ≥1 of N replicas. The pragmatic middle.

---

<a name="module-75--sharding"></a>
## Module 7.5 — Sharding & partitioning

> 📖 **Definition — Sharding:** Splitting a dataset across multiple servers, where each server owns a disjoint subset.

You shard when one server can't fit all the data, or all the QPS, or both.

### Strategies

**Range-based:** users with id 1-1M → shard A; 1M-2M → shard B. Simple but hot spots if data is unevenly distributed.

**Hash-based:** `shard_id = hash(key) % N`. Even distribution. But adding a shard requires moving most of the data (all `% N` results change when N changes).

**Consistent hashing:** see next module. Solves the "adding shards moves everything" problem.

**Directory-based:** a metadata service tracks which shard owns which key. Most flexible; adds a hop.

### What sharding hurts

- **Cross-shard queries:** "Show me all users named Alice" — must hit every shard.
- **Cross-shard transactions:** ACID across shards = distributed transactions. Expensive.
- **Joins:** the relational `JOIN` collapses without partition co-location.

The system design strategy: shard by the dimension that most queries filter on. If 90% of queries are "user X's stuff," shard by user_id. The 10% of cross-user queries get a denormalized view or a search index.

---

<a name="module-76--consistent-hash"></a>
## Module 7.6 — Consistent hashing

The scheme that DynamoDB, Cassandra, memcached clusters, CDN edge routing all use.

### The naive problem

You have N cache servers. `server = hash(key) % N`. You add one server. Now `% N+1` produces different results for almost every key. You just invalidated the entire cache.

### The fix

Imagine a ring (mod 2^32). Hash each server to a point on the ring. Hash each key to a point on the ring. The key belongs to the first server clockwise from the key's point.

```
              0
              │
       Server A
              │
       k1 ───►│  (k1's server)
              │
       Server C
              │
       k2 ───►│  (k2's server)
              │
       Server B
              │
              2^32
```

When you add Server D, only the keys between D and the previous server (going counterclockwise) need to move. Other keys stay put. **That's why it's "consistent" — adding/removing nodes minimally disrupts the mapping.**

### Virtual nodes

Naive consistent hashing has uneven distribution (random points on a ring → random gaps). Solution: each physical server gets many virtual node positions on the ring. Better balance.

### Used in

- Cassandra/DynamoDB partitioning
- memcached clusters (client-side)
- Akka's cluster sharding
- Many CDNs for edge selection

We'll skip the full implementation but it's a great optional project.

---

<a name="module-77--consensus"></a>
## Module 7.7 — Consensus: Paxos, Raft (intro)

> 📖 **Definition — Consensus:** Multiple nodes agreeing on a single value, even in the face of network failures and node crashes.

Used for:
- Leader election ("which node is in charge right now?")
- Distributed locks (etcd, Zookeeper)
- Configuration management (Kubernetes uses etcd, which uses Raft)
- Total-ordered logs (Kafka uses ZK or KRaft for cluster coordination)

### Paxos — the famous, infamous one

Leslie Lamport, 1998. Solves consensus in a partially synchronous network with crashes. Algorithmically correct, **operationally a nightmare** to understand and implement. The original paper is famously unreadable.

### Raft — Paxos for humans

Stanford, 2014. Designed to be **understandable**. Splits consensus into three sub-problems: leader election, log replication, safety. Used by etcd, Consul, TiKV, CockroachDB.

The 30-second sketch of Raft:

1. **Leader election.** Each node has a random election timeout (150-300ms). If no heartbeat from a leader in that time, the node becomes a *candidate*, increments its term, votes for itself, asks others for votes. Whoever gets a majority becomes leader.
2. **Log replication.** Clients send commands to the leader. Leader appends to its log, replicates to followers. Once a majority have it, leader marks it *committed* and applies it.
3. **Safety.** A new leader must have all committed entries from previous terms. Enforced by election rules (only a node with up-to-date log can win).

Visualize: a class election. The teacher (leader) writes assignments on the board (log). Students (followers) copy them. If teacher is sick, kids hold an election. New teacher must be one who had the most-recent notes.

You don't need to implement Raft for an interview. Knowing it exists, that it's the modern consensus algorithm, and that you'd use a library (etcd, Consul) rather than write it — that's enough.

---

<a name="module-78--kafka"></a>
## Module 7.8 — Message queues: Kafka deep-ish dive

> 📖 **Definition — Kafka:** A distributed, durable, high-throughput log. Producers append; consumers read at their own pace. Originally built at LinkedIn (2011); now Apache.

Datadog uses Kafka extensively for metric ingestion.

### Concepts

- **Topic:** a named log. Like a table.
- **Partition:** a topic is divided into partitions (the unit of parallelism). Each partition is an ordered, immutable sequence of records. Each record gets an offset.
- **Producer:** writes records to a topic. Picks which partition (round-robin, hash by key, custom).
- **Consumer:** reads from one or more partitions. Tracks its own offset.
- **Consumer group:** consumers in the same group split partitions among themselves. Two consumers in different groups both see all records.
- **Broker:** a Kafka server. A cluster has many brokers. Each partition is replicated across brokers (config: replication factor, often 3).
- **ISR (in-sync replicas):** the set of replicas currently caught up with the leader. Leader can ack writes after waiting for `min.insync.replicas` of them.

### Why this design wins

- **Persistent log, not transient queue.** You can re-read history. Replay events to rebuild state.
- **Partition = parallelism unit.** Add partitions → more consumer parallelism.
- **Sequential disk writes.** Way faster than random; modern disks do ~500 MB/sec sequential.
- **Zero-copy reads.** Kafka uses `sendfile()` to ship pages from disk to the socket without going through user space.

### Trade-offs to know

- **Ordering only within a partition.** Cross-partition ordering is your problem. Most apps shard by user/account/entity, so order within that entity is preserved.
- **Exactly-once is hard.** Kafka provides "at-least-once" by default. "Exactly-once" via idempotent producers + transactions is real but complex.
- **Operationally heavy.** Running Kafka well requires expertise. Managed services (Confluent, AWS MSK) are increasingly the norm.

### Kafka vs alternatives

- **RabbitMQ:** classic queue, push-based, transient. Good for traditional task queues. No long-term replay.
- **NATS:** modern lightweight pub/sub. Optional persistence (JetStream).
- **AWS SQS:** managed queue. Simple. Limited compared to Kafka.
- **Redis Streams:** Kafka-lite for small scale.
- **Pulsar:** Apache project, separates compute from storage.

---

<a name="module-79--caching"></a>
## Module 7.9 — Caching: layers, invalidation, hot keys

### The cache hierarchy

```
Browser cache (static assets)
   ↓
CDN (Cloudflare, CloudFront)
   ↓
Reverse proxy cache (Varnish, nginx)
   ↓
Application cache (in-process LRU)
   ↓
Distributed cache (Redis, Memcached)
   ↓
Database
```

Each layer absorbs traffic from the next. A hit at any layer = a saved trip.

### Invalidation — the hard part

Phil Karlton: *"There are only two hard things in computer science: cache invalidation and naming things."*

The strategies:

- **TTL (time to live).** Cache for N seconds; expire. Simple, sometimes-stale data.
- **Write-through:** every write goes to cache AND DB. Cache always fresh.
- **Write-behind:** write to cache; flush to DB asynchronously. Fast writes, risk of loss.
- **Cache-aside:** read tries cache → fall through to DB → populate cache. Most common.
- **Explicit invalidation:** on writes, the writer also invalidates affected cache keys.

Real systems combine these. Datadog's metric tag service uses TTL + explicit invalidation on tag changes.

### Hot keys

Most caches assume even key distribution. A hot key (one key getting 1M req/sec) breaks that — the shard owning it gets crushed.

Mitigations:
- **Add randomness:** key + small random suffix; smear across shards.
- **Replicate hot keys:** consciously cache the key on multiple nodes.
- **Local L1:** per-server in-process LRU in front of distributed cache.

Pinterest published a paper on this. Twitter does too (timeline cache for celebrity accounts).

### Thundering herd

Cache entry expires → 1000 concurrent requests all miss → 1000 hit the DB → DB collapses.

Mitigations:
- **Stale-while-revalidate:** serve stale entry while one request refreshes.
- **Per-key locks:** only one request fetches; others wait.
- **Probabilistic refresh:** refresh slightly before TTL with low probability (X-Fetch algorithm).

---

<a name="module-710--lb"></a>
## Module 7.10 — Load balancers and service discovery

> 📖 **Definition — Load balancer:** Sits between clients and a pool of servers, distributing requests. Provides scaling, failure isolation, sometimes additional features (rate limiting, TLS termination, request routing).

### L4 vs L7 (recap from Phase 1)

- **L4** load balancers operate on TCP/UDP. Fast, dumb. AWS NLB, ipvs.
- **L7** load balancers operate on HTTP. Can route by path/header/cookie. AWS ALB, nginx, Envoy, HAProxy.

### Algorithms

- **Round-robin** — equal share. Simple.
- **Least-connections** — send to whoever has fewest open connections.
- **Weighted** — assign weights for unequal capacity.
- **IP hash** — same client always goes to same server (sticky sessions).
- **Power of two choices** — pick 2 random servers, send to whichever has fewer load. Surprisingly close to optimal in practice.

### Service discovery

In a static world, the LB knows about servers via config files. In dynamic worlds (Kubernetes, autoscaling), servers come and go.

Patterns:
- **Client-side:** clients query a service registry (Consul, etcd) for healthy instances. Load balance themselves.
- **Server-side:** clients hit a stable address; an LB or sidecar dispatches.
- **DNS-based:** service has a DNS name; resolves to current IPs. Cheap; high TTL caching causes lag.

Kubernetes does service discovery via DNS + iptables/IPVS rules per node. Service meshes (Istio, Linkerd) do client-side via sidecars.

---

<a name="module-711--failures"></a>
## Module 7.11 — Failure modes you'll actually see

The catalog of "things that will kill your service":

1. **Network partition.** Two halves of your system can't reach each other. Both think the other is dead. Both might keep operating ("split brain"). Quorum-based systems (need majority) prevent this.

2. **Slow responses (worse than fast failures).** Server is up but answering at 10 sec instead of 50 ms. Without timeouts, callers stack up, threads exhausted, system collapses. **The first defense for any RPC: a timeout.**

3. **Cascading failure.** Service A is slow. Service B (which calls A) waits. B's queue fills. Callers of B wait. The whole system bricks. Mitigations: timeouts, circuit breakers, bulkheads.

4. **Thundering herd.** Already covered. Cache misses, retry storms, cron jobs all firing at midnight.

5. **GC pauses / stop-the-world.** A 5-second JVM GC stops your service. Health checks fail; LB removes node. New requests hit other nodes. They GC. Now they fail too. Domino.

6. **Clock skew.** "Last-write-wins" using clocks with NTP drift = silently lost writes. Logical clocks (vector clocks, Lamport timestamps) solve this. Spanner uses TrueTime (atomic clocks!) to skip the problem.

7. **Hot shard.** One shard's load spikes, the others coast. Throughput limited by the hot one.

8. **Dependency death.** You depend on a third party. They go down. You go down. Mitigations: cached fallbacks, graceful degradation, feature flags to disable broken paths.

9. **Configuration push gone wrong.** Pushed bad config → all nodes accept it → all nodes break simultaneously. Canary deploys; blast radius limits.

10. **Database migration.** Schema change locks a table for 2 minutes. App stalls. Mitigation: pt-online-schema-change (MySQL), `CREATE INDEX CONCURRENTLY` (Postgres), gradual rollouts.

Datadog's value prop is detecting many of these. Knowing them = knowing why customers buy.

---

<a name="module-712--framework"></a>
## Module 7.12 — The 7-step system design framework

For any system design interview prompt, run these in order. Out loud.

### 1. Clarify requirements (5 minutes)

- Functional: what features? What's the API?
- Non-functional: how many users? Read/write ratio? Latency target? Consistency requirements? Geo distribution?

If they say "design Twitter," ask: "Just timelines? Search too? Direct messages? Globally distributed?"

### 2. Capacity estimation (5 minutes)

Numbers you'll need: QPS, storage, bandwidth.

```
500M users
Each posts 1 tweet/day average → 500M tweets/day
500M / 86400 ≈ 6K writes/sec average; peak ~5x = 30K/sec

Each user reads timeline 10x/day → 5B reads/day → 60K reads/sec average
Read:write ≈ 10:1.

Each tweet: 280 chars + metadata ≈ 1 KB
500M * 1 KB = 500 GB/day, 180 TB/year for tweets alone
```

You'll often round generously. They want to see you reason, not nail the number.

### 3. API design (3 minutes)

```
POST /tweets        body: {text}
GET  /timeline      query: ?limit=20&cursor=...
POST /follow/:user
```

Stay at the level of "what do clients call?" Don't go into protocols yet.

### 4. High-level architecture (5 minutes)

Boxes and arrows. Web → LB → API servers → cache + DB + queue + workers.

```
[Client] → [CDN] → [API Gateway] → [Service A]
                                  → [Service B]
                                       ↓
                               [Cache] [DB] [Kafka]
                                              ↓
                                          [Workers]
```

### 5. Data model & storage (10 minutes)

Decide:
- Which DB(s)? Postgres, Cassandra, Redis, S3, ElasticSearch?
- Schema (just the key tables)
- Indexes
- Partitioning / sharding key

This is where most senior interviewers spend most time. They want to see you understand the access patterns and pick storage to match.

### 6. Identify bottlenecks & scale (10 minutes)

Pick 2-3 hot points and zoom in.

- Hot timeline reads → cache layer (Redis cluster, replicate hot keys)
- Write fanout for celebrities → hybrid push/pull (push for normal users, pull for celebs)
- Search → ElasticSearch cluster
- Analytics queries → separate columnar store (BigQuery, Snowflake)

### 7. Address failures + observability (5 minutes)

- What if the cache is cold?
- What if the DB primary fails?
- What if a region goes down?
- What metrics do you watch? What do you page on?
- How do you deploy without downtime?

You won't always have time for all 7 steps. But hit at least 1, 2, 4, 5, 6.

---

<a name="module-713--designs"></a>
## Module 7.13 — Six canonical designs

Each below is a sketch. Practice expanding each into a 30-minute interview answer.

### A. URL shortener (the "warmup")

- Functional: POST URL → short code; GET short code → 301 to long URL.
- Estimation: 100M URLs/day. Read:write = 100:1. ~70 GB storage/year.
- API: `POST /shorten {url}` → `{code}`. `GET /:code`.
- Architecture: Stateless API → Redis (hot codes) → Postgres (durable).
- Code generation: base62 encode of an auto-increment ID? Or random 7-char string + collision retry. Both have trade-offs.
- Scale: easy — reads from Redis, write through to Postgres. Postgres can handle the writes.
- Failures: Redis loss = warm cache from DB. Postgres replication for durability.

### B. Twitter-style timeline

- The hardest part: how do users see "their feed" of people they follow?
- Two approaches:
  - **Pull (read fanout):** on read, fetch recent tweets from each followed user, merge. Cheap writes, expensive reads.
  - **Push (write fanout):** on write, fan out to all followers' inboxes. Cheap reads, expensive writes (Lady Gaga writes once → 100M inboxes).
- Hybrid: push for normal users, pull for celebrities. Merge both at read time.
- Storage: tweets in a key-value store (Cassandra/DynamoDB), inbox in Redis with TTL.

### C. Distributed rate limiter (Redis-backed)

- Per-user, per-IP, per-API-key keys.
- Algorithm: token bucket via Redis Lua. Atomic. ~50 µs latency.
- Sharding: by client key (`hash(user_id) % N`).
- Geo: each region has its own Redis cluster; counters are local. Don't try global synchronous counts.

### D. Distributed cache (memcached/Redis cluster)

- Consistent hashing for key distribution.
- Replication factor 2 for hot keys.
- Eviction: LRU (which we built in Phase 3!).
- Failure: virtual nodes + replication; on node loss, traffic distributes evenly.

### E. Time-series database (preview of Datadog Husky)

- Workload: 10M points/sec ingest. Queries: aggregations over 5min/1hr/1day windows by tags.
- Architecture:
  - **Ingest:** load-balanced API → Kafka (durable buffer) → ingestion workers → write to columnar files in object store
  - **Storage:** Parquet/ORC files in S3, partitioned by `metric x hour x shard`
  - **Index:** inverted index on tag values (which series IDs match `service:web`?)
  - **Query:** map-reduce style. Coordinator splits query by partition, fans out to workers, merges results.
  - **Rollups:** background jobs aggregate raw points into 1min, 1hr, 1day rollups for older data.
- Trade-offs: AP / EL. Eventual consistency on the last 30 seconds. Strong rollups for older queries.

### F. Jupyter-style notebook backend (your team!)

- Frontend (React) ↔ backend ↔ kernel pool.
- A kernel = a Python/Spark process running user code.
- Backend orchestrates kernel lifecycle in Kubernetes (one pod per kernel).
- Comm: WebSocket front to back; ZeroMQ back to kernel (Jupyter wire protocol).
- Storage: notebooks in Postgres (metadata) + object store (cells, outputs).
- Sharing: per-user RBAC. Real-time collab via OT/CRDT (advanced).
- Why Datadog cares: data engineers explore data in notebooks, build queries, share insights.

We'll build a mini version of this in Phase 9.

---

<a name="papers"></a>
## Papers worth reading

Reading classic distributed systems papers is what separates pretenders from people who can actually design these. Read 1-2 a month, forever.

- **MapReduce** (Dean, Ghemawat 2004) — the paper that started Hadoop and modern data engineering.
- **The Google File System** (Ghemawat 2003) — distributed storage; precursor to HDFS.
- **Bigtable** (Chang 2006) — Google's wide-column store. Inspired HBase, Cassandra.
- **Dynamo: Amazon's Highly Available Key-value Store** (DeCandia 2007) — defined "AP" systems for a generation. Inspired Cassandra, DynamoDB.
- **Raft** (Ongaro, Ousterhout 2014) — consensus made understandable.
- **The Tail at Scale** (Dean, Barroso 2013) — why p99 matters more than mean. Required reading at any scale-out company.
- **Designing Data-Intensive Applications** (Kleppmann book) — the textbook every distributed systems engineer recommends. Read it cover to cover. Twice.

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

These are meta-questions that come up *outside* the system design itself.

1. **Walk me through a recent outage.** *(If you don't have one, talk through a hypothetical one and how you'd debug.)*

2. **CAP theorem in one minute.**

3. **PACELC — why is it more useful than CAP?**

4. **Strong vs eventual consistency. When is each acceptable?**

5. **Explain consistent hashing.**

6. **What does Raft actually do?**

7. **Why is `now()` dangerous in distributed systems?**

8. **What's a Bloom filter? When would you use one?**

9. **What's a vector clock?**

10. **Walk me through a Kafka commit.**

11. **Explain the difference between at-most-once, at-least-once, exactly-once delivery.**

12. **You see p50 latency at 50 ms but p99 at 5 seconds. What's happening?** *(GC pauses, queue buildup behind a slow node, dependent service latency, hot shard, retries, head-of-line blocking. Standard answers.)*

13. **You need to deploy a code change to 1000 servers. How?** *(Canary → small percentage → monitor → ramp. Feature flags for rollback without redeploy.)*

14. **Idempotency — give 3 examples of why you'd need it.**

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] Why distributed; the 8 fallacies
- [ ] Latency numbers in your bones
- [ ] CAP, PACELC, eventual consistency anomalies
- [ ] Replication strategies
- [ ] Sharding and partitioning + consistent hashing
- [ ] Consensus at a sketch level (Raft)
- [ ] Kafka concepts and trade-offs
- [ ] Caching strategies and pitfalls
- [ ] Load balancers, service discovery
- [ ] The 10 failure modes
- [ ] The 7-step system design framework
- [ ] Six canonical designs you can fluently expand
- [ ] You've started the paper-reading habit

---

**Next:** [Phase 8 — Datadog Stack](../phase-08-datadog-stack/README.md)
