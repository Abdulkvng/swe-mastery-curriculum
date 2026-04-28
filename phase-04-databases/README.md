# Phase 4 — Databases Deep Dive

> Every backend stores state somewhere. The database is where bugs hide, where outages start, and where senior engineers earn their pay. At Datadog you'll work with Postgres (metadata), Redis (caching), and time-series databases (Husky, their custom one). At Apple, you'll see Cassandra, Foundation DB, and many proprietary systems. Across all of them: the *fundamentals* are the same — ACID, isolation, indexes, replication. Master those once, and the surface vocabulary changes but the mental model holds.
>
> This phase goes deep on Postgres (relational), Redis (key-value, in-memory), and the index data structures (B-trees, hash indexes, LSM trees) that make them fast. You'll extend the Phase 2 connection pool into a mini-pgbouncer.

**Time:** 6–9 days.

**You'll know you're done when:** you can explain ACID with examples; you can predict whether a query will use an index by reading EXPLAIN; you can defend why you'd pick Postgres vs Redis vs DynamoDB vs Cassandra for a given workload; and you've built a simplified Postgres connection proxy.

---

## Table of contents

1. [What does this even mean? — "Database"](#what-database-means)
2. [Module 4.1 — Relational vs NoSQL: the actual differences](#module-41--relational-vs-nosql)
3. [Module 4.2 — ACID, properly](#module-42--acid)
4. [Module 4.3 — Transaction isolation levels](#module-43--isolation)
5. [Module 4.4 — Indexes: B-tree, hash, GIN, BRIN](#module-44--indexes)
6. [Module 4.5 — Postgres deep dive](#module-45--postgres)
7. [Module 4.6 — Query planning and EXPLAIN](#module-46--explain)
8. [Module 4.7 — Replication, partitioning, sharding](#module-47--replication)
9. [Module 4.8 — Redis deep dive](#module-48--redis)
10. [Module 4.9 — Time-series databases (preview of Datadog Husky)](#module-49--tsdb)
11. [Module 4.10 — When to pick what](#module-410--pick-what)
12. [🛠️ Project: Mini-pgbouncer](#project-pgbouncer)
13. [🛠️ Project: B-tree from scratch](#project-btree)
14. [Exercises](#exercises)
15. [Interview question bank](#interview-questions)
16. [What you should now know](#what-you-should-now-know)

---

<a name="what-database-means"></a>
## 🧠 What does this even mean? — "Database"

A **database** is software that stores data and lets you query it. That's the boring definition. The interesting one: a database is software that solves the four problems your application doesn't want to solve itself.

The four problems:

1. **Durability.** If the power goes out mid-write, did the data survive? Files alone don't promise this without `fsync` and careful writing.
2. **Concurrent access.** Two users editing the same row at the same time. Without coordination, one wins, one loses, sometimes both lose.
3. **Searchable structure.** "Find all users named Sarah who signed up last month" — without indexes, you'd scan every row.
4. **Schema and integrity.** Foreign keys, types, constraints — guarantees that bad data can't sneak in.

Databases are the place where your application stops trusting the OS to handle data, and starts trusting code that's been beaten on by 30+ years of weird bugs and edge cases. That's why we use them instead of just writing files.

---

<a name="module-41--relational-vs-nosql"></a>
## Module 4.1 — Relational vs NoSQL: the actual differences

You've heard "relational vs NoSQL" 200 times. Let's be precise.

### Relational (SQL)

> 📖 **Definition — Relational database:** Data is organized into *tables* (rows + columns) with strict schemas. Tables relate to each other via *foreign keys*. Queries are written in **SQL**. Examples: Postgres, MySQL, SQLite, Oracle, SQL Server.

Strengths:
- **ACID transactions** — multiple writes succeed or fail atomically
- **Joins** — combine data across tables in one query
- **Strong consistency** — what you read is what was last committed
- **Mature ecosystem** — every tool, library, hire knows SQL

Weaknesses:
- **Vertical scaling first** — scaling beyond one big server is operationally painful
- **Schema migrations** — changing tables on a live system is non-trivial
- **Some workloads are awkward** — heavily nested documents, graph queries, time-series

### NoSQL (a fuzzy umbrella for "not strictly relational")

NoSQL is four different things people lump together:

| Type | Examples | Use when |
|---|---|---|
| **Key-value** | Redis, DynamoDB, etcd | Simple lookups, caches, sessions |
| **Document** | MongoDB, CouchDB, Firestore | Nested objects, flexible schemas |
| **Wide-column** | Cassandra, ScyllaDB, HBase, Bigtable | Massive write volumes, time-series-ish |
| **Graph** | Neo4j, Dgraph | Heavily connected data (social networks, fraud detection) |

NoSQL databases generally trade *some* consistency for scale and/or flexibility. They are NOT necessarily faster than Postgres — Postgres on modest hardware regularly outperforms small NoSQL clusters. The real question is operational: at what scale does each break, and what do you give up?

### The "right answer" for most companies

Postgres until you can't. Then Postgres + Redis. Then Postgres + Redis + a specialized DB for the workload that's hurting (search → Elastic, time-series → Influx/Datadog Husky, queue → Kafka). That progression covers ~90% of real systems.

---

<a name="module-42--acid"></a>
## Module 4.2 — ACID, properly

ACID is four guarantees a transactional database promises. Memorize the letters AND have a 30-second example for each.

### A — Atomicity

> A transaction is *all-or-nothing*. Either every statement commits, or none of them do.

Example: transferring $100 from account A to account B requires two writes:
```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 'A';
UPDATE accounts SET balance = balance + 100 WHERE id = 'B';
COMMIT;
```

If the server crashes between these statements, atomicity guarantees the first update is also rolled back. Money doesn't vanish.

How it's implemented: **write-ahead log (WAL)**. Before changing the actual data files, the database appends the change to a log file (`fsync`'d to disk). On crash recovery, the DB replays the log and rolls back any incomplete transactions.

### C — Consistency

> A transaction takes the database from one valid state to another. Constraints (foreign keys, NOT NULL, CHECK) are enforced.

Note: "consistency" in ACID is **different** from "consistency" in CAP theorem. ACID-C is about constraint enforcement. CAP-C is about replica agreement. Don't confuse them — interviewers do this on purpose to see if you notice.

### I — Isolation

> Concurrent transactions don't see each other's intermediate state.

This is the deepest letter. There are *levels* of isolation, with trade-offs against performance. Module 4.3 covers them.

### D — Durability

> Once a transaction is committed, it survives crashes, power loss, and the server falling out a window.

Implementation: WAL again, plus periodic checkpoints to flush the dirty pages to the actual data files.

The fast/slow knob: how many `fsync()` calls per commit. Postgres' default is "fsync the WAL on every commit" — durable but limits commit throughput. Some apps (Datadog metrics intake, for example) accept "write to memory, fsync on a schedule" — much faster, with a small window of data loss on crash.

### When ACID isn't enough — distributed transactions

ACID is well-defined on a *single* database. Spread the transaction across two databases (or two services), and you need extra protocols: two-phase commit (slow, blocking), three-phase commit, or pragmatic "saga" patterns (compensating transactions). We'll touch this in Phase 7.

---

<a name="module-43--isolation"></a>
## Module 4.3 — Transaction isolation levels

The SQL standard defines four levels. Each prevents specific concurrency bugs. Stronger = safer but slower.

### The bugs each level prevents

| Bug | What it means | Example |
|---|---|---|
| **Dirty read** | T1 reads T2's *uncommitted* changes. T2 rolls back. T1 saw garbage. | T1 reads balance=200 (T2 set it but rolled back). T1 acts on bad data. |
| **Non-repeatable read** | T1 reads a row twice. T2 updates and commits in between. T1 sees two different values. | T1 reads balance=100, T2 sets to 50, T1 reads again, sees 50. |
| **Phantom read** | T1 runs a range query twice. T2 inserts a new matching row. T1's second query has more rows. | T1: `SELECT * WHERE age > 30` → 5 rows. T2 inserts a 31yo. T1 re-runs → 6 rows. |
| **Lost update / write skew** | Two transactions read the same data, both write — one's write effectively erases the other's logic. | Two booking systems read "1 seat left," both reserve it, both succeed. |

### The four levels

| Level | Dirty | Non-repeatable | Phantom |
|---|---|---|---|
| READ UNCOMMITTED | ⚠️ | ⚠️ | ⚠️ |
| READ COMMITTED | ✅ | ⚠️ | ⚠️ |
| REPEATABLE READ | ✅ | ✅ | ⚠️* |
| SERIALIZABLE | ✅ | ✅ | ✅ |

*Postgres' "REPEATABLE READ" actually prevents phantoms via MVCC snapshots — stronger than the SQL standard requires.

### What Postgres actually does

Postgres uses **MVCC (Multi-Version Concurrency Control)**:

> 📖 **Definition — MVCC:** Each row update creates a *new version* of the row. Readers see a consistent snapshot of the database; they don't block writers, writers don't block readers.

Each row has hidden columns `xmin` (transaction that created it) and `xmax` (transaction that deleted/superseded it). When you query, Postgres filters rows visible to YOUR transaction's snapshot.

Cost of MVCC: **bloat**. Old row versions stick around until VACUUM (Postgres' garbage collector) removes them. A heavily-updated table can balloon in size if VACUUM falls behind. This is one of the most common Postgres operational issues at scale.

### Defaults

- Postgres: READ COMMITTED
- MySQL InnoDB: REPEATABLE READ
- Oracle: READ COMMITTED

In application code:
```sql
BEGIN;
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- statements
COMMIT;
```

> 🎯 **Interview trap:** "Two transactions both read the same balance, both subtract $50, both commit. What happened?" — In READ COMMITTED, both succeed and you've lost $50. Fix: SERIALIZABLE, or use `SELECT ... FOR UPDATE` to take a row lock, or use an atomic update (`UPDATE accounts SET balance = balance - 50 WHERE id = ?`).

---

<a name="module-44--indexes"></a>
## Module 4.4 — Indexes: B-tree, hash, GIN, BRIN

> 📖 **Definition — Index:** A separate data structure on disk that lets the database find rows matching a predicate without scanning every row.

Without an index, `SELECT * FROM users WHERE email = 'k@x.com'` scans every row in `users` — O(n). With a B-tree index on `email`, it's O(log n).

The trade-off: indexes speed up reads but slow down writes (every insert/update has to update the indexes too) and use disk space.

### B-tree (the default everywhere)

> 📖 **Definition — B-tree:** A self-balancing tree where each node has many keys (not just two). Optimized for disk I/O — minimize the number of nodes you have to read.

Why "many keys per node"? Because each disk read pulls a whole *page* (typically 8 KB in Postgres, 16 KB in InnoDB). You want each page to contain as much useful info as possible. A binary tree with one key per node would mean ~30 disk reads for a billion-row table; a B-tree with 100 keys per node needs ~5.

What B-trees are good for:
- Equality lookups: `WHERE email = ?`
- Range queries: `WHERE created_at > '2025-01-01'`
- Prefix matching: `WHERE name LIKE 'Kv%'` (but NOT `LIKE '%kv%'`)
- ORDER BY (the index already maintains sort order)

We'll implement a B-tree from scratch in `projects/btree-from-scratch/`.

### Hash indexes

Constant-time equality lookup. Worse than B-tree for ranges.

In Postgres, hash indexes exist but aren't commonly used because B-tree is good enough at equality and also handles ranges. Used more in in-memory KV stores.

### GIN (Generalized Inverted Index)

For "many values per row." Postgres uses GIN for:
- Full-text search (`tsvector`)
- Array containment (`tags @> ARRAY['datadog']`)
- JSON key/value lookup (`jsonb`)

Conceptually: an inverted index. `tag → list of row IDs containing that tag`. Same idea as Lucene/Elasticsearch.

### BRIN (Block Range Index)

For very large tables where data has natural ordering — typically time-series. Stores min/max values for each block range. Tiny on disk, but only useful when correlation between physical layout and logical order is high.

Datadog uses BRIN-like indexes heavily for time-series ingestion, where new data always lands at "now."

### LSM trees (covered conceptually)

> 📖 **Definition — LSM (Log-Structured Merge) tree:** Writes go to an in-memory sorted table; periodically flushed to immutable on-disk files; periodically compacted (merged + deduplicated). Used by Cassandra, RocksDB, LevelDB, ScyllaDB, HBase, ClickHouse storage layer.

Why: writes are *append-only* (fast). Reads check memtable first, then on-disk files in newest-first order.

Trade-off: read amplification (a read may have to check several files). Compactions consume CPU/IO. Tombstones for deletes (a delete is just a marker that gets compacted away later).

LSM is the dominant data structure for write-heavy storage in the last decade. Datadog's Husky uses an LSM-flavored storage engine.

### Composite (multi-column) indexes

```sql
CREATE INDEX users_state_age ON users(state, age);
```

Useful for queries filtering on both columns. **Order matters**:
- `WHERE state = 'CA' AND age = 21` → uses index ✅
- `WHERE state = 'CA'` → uses index (prefix) ✅
- `WHERE age = 21` → does NOT use index ❌ (no `state` predicate)

Rule of thumb: order columns from "most filtering" to "least," and equality columns before range columns.

### Partial indexes

```sql
CREATE INDEX users_active ON users(id) WHERE deleted_at IS NULL;
```

Smaller, faster. Only indexes the rows you'll actually query.

### Covering indexes

```sql
CREATE INDEX users_email_inc ON users(email) INCLUDE (name, signup_date);
```

The included columns sit in the index but aren't part of the search key. Lets queries return without touching the table itself ("index-only scan").

---

<a name="module-45--postgres"></a>
## Module 4.5 — Postgres deep dive

Postgres is the relational database to know in 2026. Datadog uses it for metadata. Apple uses it widely. It's the "default" choice you should be able to defend.

### Setup locally

```bash
brew install postgresql@16
brew services start postgresql@16

# Create a DB and user
createdb mydb
psql mydb
```

### Schema design — the basics

```sql
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE posts (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    body        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX posts_user_id ON posts(user_id);
CREATE INDEX posts_created_at ON posts(created_at);
```

Notes:
- `BIGSERIAL` = auto-incrementing 8-byte int. Use `BIGINT GENERATED ALWAYS AS IDENTITY` if you want SQL-standard syntax.
- `TIMESTAMPTZ` = timestamp WITH timezone. **Always use this**, never plain `TIMESTAMP`. Saves you from a lifetime of timezone bugs.
- `ON DELETE CASCADE` = when a user is deleted, their posts auto-delete.

### Queries

```sql
-- All posts by a user, newest first
SELECT id, title, created_at
FROM posts
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 20;

-- Inner join
SELECT u.name, p.title
FROM users u
JOIN posts p ON p.user_id = u.id
WHERE u.deleted_at IS NULL;

-- Aggregation
SELECT u.id, u.name, COUNT(p.id) AS post_count
FROM users u
LEFT JOIN posts p ON p.user_id = u.id
GROUP BY u.id, u.name
ORDER BY post_count DESC
LIMIT 10;

-- Window function: rank posts within each user's history
SELECT
    user_id,
    title,
    created_at,
    ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) AS rn
FROM posts;
```

### Postgres-specific superpowers

#### JSONB

```sql
CREATE TABLE events (
    id    BIGSERIAL PRIMARY KEY,
    data  JSONB NOT NULL
);

INSERT INTO events(data) VALUES ('{"type":"login","ip":"1.2.3.4"}');

SELECT data->>'type', data->'ip'
FROM events
WHERE data @> '{"type":"login"}';

CREATE INDEX events_data ON events USING GIN (data);
```

JSONB is a binary, indexed JSON. Real Postgres tables can absorb document-style data without going NoSQL.

#### Arrays

```sql
CREATE TABLE products (
    id    BIGSERIAL PRIMARY KEY,
    tags  TEXT[]
);

INSERT INTO products(tags) VALUES (ARRAY['datadog', 'observability']);
SELECT * FROM products WHERE 'datadog' = ANY(tags);
```

#### CTEs (Common Table Expressions)

```sql
WITH active_users AS (
    SELECT id, name FROM users WHERE deleted_at IS NULL
),
recent_posts AS (
    SELECT user_id, COUNT(*) AS c FROM posts
    WHERE created_at > now() - interval '7 days'
    GROUP BY user_id
)
SELECT au.name, COALESCE(rp.c, 0) AS posts_this_week
FROM active_users au
LEFT JOIN recent_posts rp ON rp.user_id = au.id
ORDER BY posts_this_week DESC;
```

#### LISTEN / NOTIFY

Pub/sub built into Postgres. Useful for cache invalidation, lightweight queueing.

```sql
LISTEN cache_invalidate;
-- Another connection:
NOTIFY cache_invalidate, 'user:42';
```

#### Range types

```sql
SELECT * FROM bookings
WHERE during && tstzrange('2025-04-27 10:00', '2025-04-27 11:00');
```

`&&` is "overlaps." Useful for booking systems, calendar collisions.

---

<a name="module-46--explain"></a>
## Module 4.6 — Query planning and EXPLAIN

A query goes through:
1. **Parse** — text to AST
2. **Plan** — pick an execution strategy (which indexes, which join algorithm, which order)
3. **Execute** — run the plan

`EXPLAIN ANALYZE` shows the chosen plan AND actual runtime.

```sql
EXPLAIN ANALYZE
SELECT * FROM users WHERE email = 'k@x.com';
```

Output (annotated):
```
Index Scan using users_email_key on users
    (cost=0.29..8.31 rows=1 width=140)
    (actual time=0.024..0.025 rows=1 loops=1)
  Index Cond: (email = 'k@x.com'::text)
Planning Time: 0.142 ms
Execution Time: 0.045 ms
```

What to look for:
- **Index Scan** ✅ vs **Seq Scan** ⚠️ — sequential = full table read
- **rows=** estimate vs actual — big mismatch → stale statistics, run `ANALYZE`
- **Sort** with `external merge` → query spilling to disk; consider raising `work_mem`
- **Nested Loop** vs **Hash Join** vs **Merge Join** — three join algorithms; planner picks based on row estimates and indexes

### Common Postgres performance issues

1. **Sequential scan on a large table** → missing or unusable index. Check `pg_stat_user_tables.seq_scan`.

2. **Bloat** → vacuum falling behind. Check `pg_stat_user_tables.n_dead_tup`. Solution: tune autovacuum, occasionally `VACUUM FULL` (locks the table — schedule during off-hours).

3. **Lock contention** → `SELECT * FROM pg_locks JOIN pg_stat_activity USING (pid)` to find blocking queries. Long-running transactions are usually the culprit.

4. **N+1 queries** → an app loop that runs one query per item. Fix in app code: batch into a single `WHERE id IN (...)` or join.

5. **Connection exhaustion** → too many open connections; each consumes ~10MB. Solution: connection pooler (pgbouncer, or your project from Phase 2!).

---

<a name="module-47--replication"></a>
## Module 4.7 — Replication, partitioning, sharding

### Replication

> 📖 **Definition — Replication:** Keeping copies of the data on multiple machines. The "primary" (or "leader") accepts writes; "replicas" (or "followers") receive a stream of changes.

Why: **read scaling**, **failover**, **geo-distribution**, **backups**.

Two flavors:
- **Synchronous:** primary waits for replica acknowledgment before commit. Stronger durability, higher commit latency.
- **Asynchronous:** primary commits, then ships. Lower latency, small window of potential data loss on primary crash.

Most production systems use *semi-synchronous*: primary waits for at least one of N replicas.

Replication can be **logical** (replays SQL statements / row changes) or **physical** (ships WAL bytes verbatim). Physical is faster but requires identical Postgres versions; logical is more flexible (subscribe to specific tables, replicate across versions).

### Partitioning (within one DB)

> 📖 **Definition — Partitioning:** Splitting a single logical table into multiple physical tables under the hood, by some key (usually time or hash).

```sql
CREATE TABLE events (
    id          BIGINT,
    created_at  TIMESTAMPTZ NOT NULL,
    data        JSONB
) PARTITION BY RANGE (created_at);

CREATE TABLE events_2026_q1 PARTITION OF events
    FOR VALUES FROM ('2026-01-01') TO ('2026-04-01');
CREATE TABLE events_2026_q2 PARTITION OF events
    FOR VALUES FROM ('2026-04-01') TO ('2026-07-01');
```

Queries on `events` automatically dispatch to the right partition. Old partitions can be dropped in O(1) (just `DROP TABLE`) — much faster than `DELETE WHERE created_at < ...`.

### Sharding (across multiple DBs)

> 📖 **Definition — Sharding:** Distributing data across multiple database servers, where each server holds a *subset* of the data. Different from replication — replicas have full copies; shards have disjoint pieces.

Example sharding strategies:
- **Range:** users with `id < 1M` → shard 1, `1M ≤ id < 2M` → shard 2. Simple but uneven.
- **Hash:** `shard_id = hash(user_id) % num_shards`. Even distribution but hard to add new shards.
- **Consistent hashing:** maps keys onto a ring; adding a shard only moves a fraction of keys. Used by Cassandra, DynamoDB.
- **Directory:** a separate service tracks which shard owns which key. Flexible but adds a hop.

Sharding hurts: cross-shard queries, cross-shard transactions, joins. Postgres native sharding (Citus) and "logical sharding in the app layer" are the two common patterns.

### CAP and PACELC (cliffs notes — full treatment in Phase 7)

CAP: in a network partition, you must choose **C**onsistency (refuse to serve potentially-stale reads) or **A**vailability (serve them). PACELC: even when there's no partition, choose between **L**atency and **C**onsistency.

Postgres: CP. Cassandra: AP. Dynamo: AP (tunable). Spanner: tries to be CP with very low latency via atomic clocks. Datadog Husky: AP (eventual consistency on metric ingestion is fine — last 30 seconds of data may not be queryable yet).

---

<a name="module-48--redis"></a>
## Module 4.8 — Redis deep dive

> 📖 **Definition — Redis:** An in-memory key-value store. Single-threaded. Wickedly fast. Supports rich data structures (lists, sets, sorted sets, hashes, streams, HyperLogLog, geospatial).

Used for: caching, session storage, rate limiting, job queues, leaderboards, real-time counters, pub/sub, distributed locks.

### Setup

```bash
brew install redis
brew services start redis
redis-cli
```

### Core commands

```bash
# Strings (KV)
SET user:42 '{"name":"Kvng"}'
GET user:42
SETEX session:abc 3600 '...'    # set with expiration in seconds
INCR counter:visits
INCRBY counter:bytes 1024

# Lists (linked list, push/pop both ends)
LPUSH queue:tasks "task1"
RPUSH queue:tasks "task2"
LPOP queue:tasks
BRPOP queue:tasks 0              # blocking pop (worker pattern)

# Sets
SADD tags:post:42 "datadog" "observability"
SISMEMBER tags:post:42 "datadog"
SINTER tags:post:42 tags:post:43

# Sorted sets (ranked)
ZADD leaderboard 1500 "kvng"
ZADD leaderboard 1200 "alice"
ZRANGE leaderboard 0 9 REV WITHSCORES   # top 10
ZINCRBY leaderboard 50 "kvng"

# Hashes (object-like)
HSET user:42 name "Kvng" age 21
HGET user:42 name
HGETALL user:42

# Pub/sub
SUBSCRIBE channel:events
PUBLISH channel:events "hello"

# TTL
EXPIRE user:42 600
TTL user:42
```

### Persistence options

Redis has two persistence modes (you can use both):

- **RDB:** point-in-time snapshots. Compact files, fast restart. Loses up to N seconds of writes on crash.
- **AOF:** append-only log of every write. Rewrite/compact periodically. Configurable fsync (every write, every second, never). Stronger durability, larger files.

Production Redis at any serious scale runs replicas + Redis Sentinel (failover) or Redis Cluster (sharding).

### Common Redis patterns

#### Cache-aside

```go
func GetUser(id int64) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    if cached, err := redis.Get(key); err == nil {
        return decode(cached), nil
    }
    u, err := db.QueryUser(id)
    if err != nil { return nil, err }
    redis.SetEX(key, 300, encode(u))   // 5min TTL
    return u, nil
}
```

#### Rate limiting (token bucket)

```lua
-- KEYS[1] = key, ARGV[1] = max, ARGV[2] = window seconds
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
end
if current > tonumber(ARGV[1]) then
    return 0   -- rate-limited
end
return 1
```

Atomic via Redis' single-threaded execution. We'll use this in Phase 5.

#### Distributed lock

`SET key value NX EX 30` — set if not exists, 30s expiry. Acquire = success. Release = delete (with check). The famous `Redlock` algorithm for multi-node Redis is contested in correctness — for most apps, single-instance Redis lock is fine.

#### Streams (Kafka-lite, in Redis)

```
XADD events * type login user_id 42
XREAD STREAMS events 0
```

Used by Datadog and many companies as a lightweight job queue / event log when full Kafka is overkill.

### Redis is single-threaded — implications

One core processes commands. This is *deliberate*: no locks needed inside Redis, simple consistency model, predictable latency.

Implication: a single slow command (`KEYS *` on a million-key DB, big `LRANGE`, expensive Lua script) blocks *every other client*. NEVER run `KEYS *` in production. Use `SCAN` instead — cursored, non-blocking.

---

<a name="module-49--tsdb"></a>
## Module 4.9 — Time-series databases (preview of Datadog Husky)

You'll work directly with Datadog's time-series infrastructure on the ADP team. Knowing the shape of these systems is critical.

### What's special about time-series data

- **Append-mostly:** new data lands at "now." Old data rarely changes.
- **Bursty:** ingestion volume can be 10M points/sec.
- **Queries are aggregations:** "average CPU per host over the last 24h, grouped by service."
- **Retention tiers:** keep raw 7 days, 1-min rollups for 30 days, 1-hour rollups for 1 year.
- **Cardinality matters:** "unique series count" (`metric x tags`) drives cost. Low cardinality → cheap. High cardinality (per-user metrics) → very expensive.

### Common architectural patterns

- **Columnar storage** — store all values for one metric contiguously. Compresses brilliantly (similar values).
- **Time-bucketed files** — one file per metric per hour (or per shard). Old hours are read-only and can be aggressively compressed/rolled up.
- **Inverted index on tags** — `service:web AND env:prod` → bitmap of matching series IDs.
- **Pre-aggregation** — store rollups (avg/min/max/sum/count per minute) so queries don't scan raw points.

### TSDBs you should know

| TSDB | Notes |
|---|---|
| **Datadog Husky** | Internal, custom. Object-store backed. Cost-efficient at scale. |
| **Prometheus** | Open source, pull-based ("Prometheus scrapes targets"). Used massively. |
| **InfluxDB** | Open source, push-based. Newer engine "IOx" is Apache Arrow + Parquet. |
| **TimescaleDB** | Postgres extension. SQL on time-series. |
| **VictoriaMetrics** | Prometheus-compatible, more efficient. Open source. |

We'll dive deeper into Husky's architecture in Phase 8.

---

<a name="module-410--pick-what"></a>
## Module 4.10 — When to pick what

The "what database should I use?" question, demystified.

| Need | Pick |
|---|---|
| Transactional app (orders, users, billing) | **Postgres** |
| Cache | **Redis** |
| Session store | **Redis** |
| Queue, light | **Redis Streams**, **Postgres SKIP LOCKED** |
| Queue, heavy | **Kafka**, **SQS** |
| Full-text search | **Elasticsearch**, **Postgres tsvector** for small scale |
| Time-series metrics | **Datadog**, **Prometheus**, **InfluxDB** |
| Massive write throughput, eventual consistency OK | **Cassandra**, **DynamoDB** |
| Document store with flexible schema | **MongoDB**, or **Postgres JSONB** |
| Graph queries (friend-of-friend) | **Neo4j** |
| Analytics over huge datasets | **Snowflake**, **BigQuery**, **ClickHouse**, **Redshift** |
| Distributed transactions across globe | **Spanner**, **CockroachDB** |
| Embedded / mobile | **SQLite** |

> 🎯 **Interview tip:** When asked "which database for X?", restate the workload first ("read-heavy, eventually consistent, ~10K writes/sec, 1KB rows, query patterns are mostly key-lookup"), then justify your pick with 2-3 concrete reasons.

---

<a name="project-pgbouncer"></a>
## 🛠️ Project: Mini-pgbouncer

Extend the Phase 2 connection pool into a real **TCP proxy** that speaks the Postgres wire protocol just enough to share a small pool of upstream Postgres connections among many client connections.

### Why this is the project

- Forces you to read a real wire protocol (Postgres') and parse it
- Combines networking (Phase 1), OOP (Phase 2), DB knowledge (this phase)
- Real production tool — pgbouncer literally exists and Datadog uses it

### Spec

```
clients (many)  ───► [mini-pgbouncer:5433] ───► [postgres:5432]
                                         (small pool of N upstream conns)
```

- **Listen** on a TCP port for client connections
- For each client, negotiate the Postgres startup message
- **Acquire** an upstream conn from the pool (your pool from Phase 2)
- For "transaction pooling" mode: hold the upstream conn only for the duration of an explicit transaction, then return it
- For each forwarded message: read from one side, write to the other
- On client disconnect: release the upstream conn

This is a STRICTLY simplified design — real pgbouncer has many more features. The point is to feel the proxy.

### Layout

```
projects/mini-pgbouncer/
├── README.md
├── go.mod
├── main.go              <- proxy main
├── proxy/
│   ├── server.go        <- TCP listen + connection handling
│   ├── protocol.go      <- enough Postgres wire protocol to do startup
│   └── session.go       <- one client session, with pool acquire/release
└── ../../../phase-02-oop-go/projects/connection-pool/   <- imported
```

> Heads up: this project is sketched here; full implementation will be a stretch goal because it requires a fair bit of Postgres-protocol parsing (messages have a 1-byte type tag + 4-byte length prefix). The skeleton is in `projects/mini-pgbouncer/`.

---

<a name="project-btree"></a>
## 🛠️ Project: B-tree from scratch (Go)

Build a working in-memory B-tree with `Insert`, `Search`, `Delete`, and `Range`.

This is what powers your Postgres index, just without the disk part. Once you understand it, "BRIN," "GIN," "LSM" become specializations of "tree-shaped index."

### Spec

- Generic `BTree[K, V]` over comparable keys
- Configurable order `t` (each node has between `t-1` and `2t-1` keys)
- `Insert(k, v)` — splits full nodes on the way down
- `Search(k)` — O(log n)
- `Delete(k)` — handles all three deletion cases (leaf with extra keys, leaf with min keys, internal node)
- `Range(lo, hi)` — return matching key/value pairs in order

Skeleton in `projects/btree-from-scratch/`. Tests included.

---

<a name="exercises"></a>
## Exercises

1. **Schema design.** Design a Postgres schema for a Twitter-clone: users, tweets, follows, likes. Include indexes. Justify each one.

2. **EXPLAIN reading.** Spin up Postgres with 1M rows in a `users` table. Run a query without an index, then add one, watch EXPLAIN go from `Seq Scan` to `Index Scan`. Compare execution time.

3. **Isolation hands-on.** In two `psql` sessions, reproduce: dirty read (won't happen in Postgres — explain why); non-repeatable read at READ COMMITTED; phantom prevention at REPEATABLE READ.

4. **Lost update fix three ways.** Take the bank-transfer scenario. Solve it with: (a) `SELECT ... FOR UPDATE` row lock, (b) atomic `UPDATE balance = balance - 50`, (c) SERIALIZABLE retry-on-conflict.

5. **Cache patterns in Redis.** Implement: cache-aside, write-through, write-behind. Discuss when each is appropriate.

6. **Rate limiter.** Build a token-bucket rate limiter as a Lua script in Redis. Test from Go with `go-redis`.

7. **Time-series storage.** Sketch (no code) an architecture for storing 100K metrics/sec, 30-day retention, with sub-second query latency for "last 1h, group by host." Identify the top 3 design choices.

8. **Replication trade-off essay (one page).** Pick a real system (Postgres streaming replication, MySQL semi-sync, MongoDB replica sets) and explain its consistency model, failover behavior, and one specific failure mode (e.g., "what happens to in-flight writes if the primary crashes mid-commit?").

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

1. **Explain ACID with a concrete example for each letter.**

2. **Walk me through the four isolation levels and what each prevents.**

3. **What is MVCC? Why does Postgres need VACUUM?**

4. **Difference between READ COMMITTED and REPEATABLE READ in Postgres specifically.**

5. **What's the difference between an index and a table?**

6. **When would you NOT want an index?** *(Heavy write tables, low-cardinality columns, indexes that won't be used.)*

7. **Why is a B-tree better than a binary search tree for disk?** *(Many keys per node = fewer disk seeks. CPUs read pages, not bits.)*

8. **Explain a hash join vs a nested loop join vs a merge join.**

9. **You have a query taking 5 seconds. How do you diagnose?** *(EXPLAIN ANALYZE, check for seq scans, check pg_stat_statements, check locks, check stale ANALYZE stats.)*

10. **What's the difference between sharding and replication?**

11. **What's MVCC bloat and how do you mitigate it?**

12. **When should you use Redis instead of Postgres?**

13. **You're building a leaderboard. Postgres or Redis? Defend.** *(Redis ZSET — `ZADD` and `ZRANGE` are O(log N). Postgres can do it but Redis is purpose-built for this.)*

14. **What's the trade-off between RDB and AOF in Redis?**

15. **What's CAP theorem? What's PACELC?**

16. **You're designing a metrics ingestion pipeline at 1M points/sec. What database?** *(Time-series: Husky/Influx/Prometheus. Defend on cardinality, query patterns, retention, compression.)*

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] Relational vs NoSQL, four flavors of NoSQL
- [ ] ACID — concrete example per letter
- [ ] Four isolation levels, three concurrency anomalies
- [ ] How MVCC works in Postgres + bloat
- [ ] B-tree, hash, GIN, BRIN indexes — when each
- [ ] LSM trees conceptually
- [ ] Postgres-specific features: JSONB, arrays, CTEs, window functions
- [ ] Reading EXPLAIN ANALYZE
- [ ] Replication, partitioning, sharding distinctions
- [ ] Redis data structures and use patterns
- [ ] Time-series database basics (preview of Datadog Husky)
- [ ] Decision tree for "which DB for what"

---

**Next:** [Phase 5 — APIs & Backend](../phase-05-apis/README.md)
