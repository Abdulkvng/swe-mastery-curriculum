# Phase 8 — Datadog Stack

> This phase is targeted prep for your specific team: **Advanced Data Platform — Notebooks**. We're going to dive into the technologies you'll encounter on day one: Apache Spark + just-enough Scala, Jupyter's actual internals (which most engineers using Jupyter have never thought about), Kubernetes + Docker (the substrate everything runs on), and the observability culture/vocabulary Datadog is famous for. By the end of this phase you'll be able to walk into your first standup and follow the conversation.

**Time:** 7–10 days.

**You'll know you're done when:** you can run a Spark job in Scala, explain Jupyter's kernel-protocol architecture in detail, deploy a service to local Kubernetes, articulate Datadog's observability product mental model, and use SLI/SLO/SLA terminology correctly.

---

## Table of contents

1. [What does this even mean? — "ADP" and "Notebooks"](#what-adp-means)
2. [Module 8.1 — Apache Spark architecture](#module-81--spark)
3. [Module 8.2 — Just-enough Scala for Spark](#module-82--scala)
4. [Module 8.3 — PySpark vs Scala Spark](#module-83--pyspark)
5. [Module 8.4 — Jupyter internals: kernels, ZeroMQ, the wire protocol](#module-84--jupyter)
6. [Module 8.5 — Docker, properly](#module-85--docker)
7. [Module 8.6 — Kubernetes essentials](#module-86--k8s)
8. [Module 8.7 — k8s operators and CRDs](#module-87--operators)
9. [Module 8.8 — Observability culture: RED, USE, SLI/SLO/SLA](#module-88--observability)
10. [Module 8.9 — OpenTelemetry deeper](#module-89--otel)
11. [Module 8.10 — Datadog the product, mental model](#module-810--datadog-product)
12. [🛠️ Project: Spark word-count + JSON log analysis](#project-spark)
13. [🛠️ Project: Mini-Jupyter kernel](#project-jupyter)
14. [🛠️ Project: Deploy a service to local k3d](#project-k3d)
15. [Day-1 vocabulary cheat sheet](#vocab)
16. [What you should now know](#what-you-should-now-know)

---

<a name="what-adp-means"></a>
## 🧠 What does this even mean? — "ADP" and "Notebooks"

**ADP** stands for **Advanced Data Platform**. It's the team at Datadog responsible for letting customers (and Datadog's own engineers) explore, query, and analyze the firehose of data Datadog ingests.

**Notebooks** specifically are Jupyter-style interactive documents where engineers can write code (Python, often using PySpark), run queries against Datadog's data, visualize results, and share findings. They're the lingua franca of data work.

So "ADP — Notebooks" means: the team that builds and runs the system letting users write live, executable, shareable analyses against Datadog's data, backed by Spark and Datadog's internal storage (Husky etc.).

What technologies you'll therefore touch:
- **Spark** for distributed query / processing
- **Jupyter** as the user interface protocol
- **Kubernetes** for orchestrating user kernel pods
- **Postgres** for metadata (notebook contents, permissions, etc.)
- **Redis** for caching, session state
- **Go and TypeScript** for backend services and frontend
- **Docker** to package everything
- **GitLab CI** for the build/test/deploy loop

This phase covers all of those.

---

<a name="module-81--spark"></a>
## Module 8.1 — Apache Spark architecture

> 📖 **Definition — Apache Spark:** A distributed compute engine for large-scale data processing. Originally MapReduce successor (2014). Now the standard for ETL, analytics, ML pipelines on big data. Datadog uses Spark heavily; you'll write Spark jobs.

### The component model

```
┌─────────────────────────────────────────────────────┐
│                  Driver Program                      │
│  - Builds the DAG (Directed Acyclic Graph) of work  │
│  - SparkContext / SparkSession                       │
└──────────────┬──────────────────────────────────────┘
               │
       ┌───────▼────────┐
       │ Cluster Manager│   (YARN, Mesos, Kubernetes, Standalone)
       └───────┬────────┘
               │
   ┌───────────┼───────────┐
   ▼           ▼           ▼
┌──────┐    ┌──────┐    ┌──────┐
│Exec 1│    │Exec 2│    │Exec 3│   (JVM processes on worker nodes)
│ Tasks│    │ Tasks│    │ Tasks│
└──────┘    └──────┘    └──────┘
```

- **Driver** — your program. Builds the plan. Owns the SparkSession.
- **Cluster manager** — spawns executors on worker nodes.
- **Executors** — run the actual computation. Hold partitions of data in memory.
- **Tasks** — units of work. One task per partition per stage.

### RDDs, DataFrames, Datasets — three abstractions

- **RDD (Resilient Distributed Dataset):** the original Spark API. Like a distributed list. Low-level; you write `map`, `filter`, `reduce`. **Used rarely now.**
- **DataFrame:** a table with named columns. Like a distributed Pandas DataFrame. Spark optimizes the query plan.
- **Dataset:** like DataFrame but with compile-time types (Scala/Java only). Best of both.

In practice: **use DataFrames** unless you have a specific reason to drop down.

### Lazy evaluation + the DAG

Spark transformations are lazy. `.filter()` doesn't compute anything. It records "filter" in the DAG. Only when you call an *action* (`.count()`, `.collect()`, `.write.save()`) does Spark actually run the DAG.

This is key. Spark optimizes the whole DAG before executing — pushing filters down, combining adjacent maps, deciding shuffle vs not, etc.

### Narrow vs wide transformations

- **Narrow:** each output partition depends on one input partition. `map`, `filter`. Cheap, no network.
- **Wide:** each output partition depends on multiple input partitions. `groupBy`, `join`, `reduceByKey`. Triggers a **shuffle** — data moves across the network.

Shuffles are expensive. The single biggest Spark performance lever is: minimize shuffles, minimize shuffle size.

### Stages and tasks

When Spark sees a wide transformation, it cuts the DAG into **stages**. Each stage is a sequence of narrow transformations between two shuffles. Within a stage, Spark fuses operations into a single pass over each partition.

Number of tasks per stage = number of partitions. Default 200 for shuffles (configurable: `spark.sql.shuffle.partitions`).

### Caching

```scala
val df = spark.read.parquet("...")
df.cache()                     // marks it cached
df.count()                     // first action populates the cache
df.filter(...).show()          // uses cache
df.groupBy(...).agg(...).show() // uses cache again
df.unpersist()                 // free the cache
```

Cache when you reuse a DataFrame multiple times. Costs memory; only worth it if recomputation cost > memory cost.

---

<a name="module-82--scala"></a>
## Module 8.2 — Just-enough Scala for Spark

You don't need to be a Scala expert. You need to read other people's Spark Scala and write modifications. ~80% of Datadog Spark code is in Scala.

### Hello world

```scala
object Hello {
    def main(args: Array[String]): Unit = {
        println("Hello, Kvng!")
    }
}
```

### Basics

```scala
// Variables
val x: Int = 42        // immutable (use this 99% of the time)
var y: String = "hi"   // mutable (rare)

val name = "Kvng"      // type inferred

// Functions
def add(a: Int, b: Int): Int = a + b
def add(a: Int, b: Int): Int = {
    val sum = a + b
    sum   // last expression is the return value
}

// Anonymous functions (lambdas)
val double = (n: Int) => n * 2
val triple: Int => Int = n => n * 3

// Collections (immutable by default)
val nums = List(1, 2, 3, 4, 5)
val doubled = nums.map(_ * 2)
val evens = nums.filter(_ % 2 == 0)
val sum = nums.reduce(_ + _)
val sumByFold = nums.foldLeft(0)(_ + _)

// Tuples
val pair: (String, Int) = ("Kvng", 21)
println(pair._1)  // Kvng

// Maps
val ages = Map("Alice" -> 30, "Bob" -> 25)
ages.get("Alice")     // Some(30)
ages.getOrElse("X", -1)
```

### Case classes (your DataFrame friends)

```scala
case class User(id: Long, name: String, age: Int)

val u = User(1L, "Kvng", 21)
println(u.name)               // "Kvng"

// Pattern matching
u match {
    case User(_, "Kvng", _)  => println("It's Kvng!")
    case User(_, n, age) if age >= 18 => println(s"$n is an adult")
    case _ => println("someone")
}
```

Case classes are the bread and butter of Spark Datasets — they map to DataFrame schema automatically.

### Option (no nulls!)

Scala discourages nulls. Use `Option[T]`:

```scala
val maybeName: Option[String] = Some("Kvng")
val nope: Option[String] = None

maybeName match {
    case Some(n) => println(n)
    case None    => println("no name")
}

// Or use map / getOrElse
val len = maybeName.map(_.length).getOrElse(0)
```

### Spark in Scala

```scala
import org.apache.spark.sql.SparkSession

object WordCount {
    def main(args: Array[String]): Unit = {
        val spark = SparkSession.builder()
            .appName("wordcount")
            .master("local[*]")     // use all local cores
            .getOrCreate()

        import spark.implicits._   // unlocks $"col" syntax and toDF/toDS

        val df = spark.read.text("README.md")
            .as[String]
            .flatMap(_.split("\\W+"))
            .filter(_.nonEmpty)
            .groupByKey(identity)
            .count()
            .orderBy($"count(1)".desc)

        df.show(20, truncate = false)

        spark.stop()
    }
}
```

That's enough Scala to read Datadog Spark code on day one. You'll learn more by reading.

---

<a name="module-83--pyspark"></a>
## Module 8.3 — PySpark vs Scala Spark

Customers write notebooks in Python. So the *user-facing* path is PySpark. The internal implementation is often Scala. You'll need both.

### The same job in PySpark

```python
from pyspark.sql import SparkSession
from pyspark.sql.functions import col, count, desc, explode, split

spark = SparkSession.builder.appName("wordcount").getOrCreate()

df = (
    spark.read.text("README.md")
    .select(explode(split(col("value"), r"\W+")).alias("word"))
    .filter(col("word") != "")
    .groupBy("word")
    .agg(count("*").alias("c"))
    .orderBy(desc("c"))
)

df.show(20, truncate=False)
spark.stop()
```

### Why PySpark is sometimes slower than Scala

PySpark JVM ↔ Python communication uses sockets + serialization (Arrow speeds it up). Scala compiles to JVM bytecode, runs directly. For pure DataFrame ops, the gap is small (the optimizer plans Scala-equivalent code). For UDFs (user-defined functions in Python), the gap is large — the JVM ships data to a Python worker per row.

Mitigations: vectorized UDFs (using pandas/Arrow), or rewrite in Scala if hot.

---

<a name="module-84--jupyter"></a>
## Module 8.4 — Jupyter internals: kernels, ZeroMQ, the wire protocol

Most data scientists treat Jupyter as magic. You're going to *be* the people who maintain that magic.

### The component model

```
┌──────────────┐         ┌─────────────────┐
│  Browser     │◄───WS──►│  Jupyter Server │
│  (notebook   │         │                  │
│   UI)        │         │                  │
└──────────────┘         └────────┬─────────┘
                                  │
                            ┌─────▼─────┐
                            │  Kernel   │  (separate process)
                            │  process  │  Python, R, Julia, ...
                            └───────────┘
```

The Jupyter Server is a Python process. The browser talks to it over WebSocket. The Server talks to one or more **kernels** (separate processes) over **ZeroMQ**.

### Why separate processes?

- Crash isolation. Bad code can't take down the server.
- Multiple languages. Same server, different kernels.
- Resource control. Each kernel can be its own k8s pod with cgroup limits.

### The Jupyter Wire Protocol

Each kernel exposes 5 ZeroMQ sockets:

- **Shell** (REQ/REP) — `execute_request`, `execute_reply`. The main code-execution channel.
- **IOPub** (PUB) — broadcasts of stdout, stderr, display data. Multiple subscribers.
- **Stdin** (REQ/REP) — for `input()` calls.
- **Control** (REQ/REP) — interrupt, shutdown messages.
- **Heartbeat** (REQ/REP) — Server pings; kernel pongs.

Each message is a multi-part ZMQ message: identities, delimiter, signature (HMAC), header (JSON), parent header, metadata, content.

```json
// header
{
    "msg_id": "uuid",
    "session": "uuid",
    "username": "kvng",
    "msg_type": "execute_request",
    "version": "5.3"
}
// content (for execute_request)
{
    "code": "1 + 1",
    "silent": false,
    "store_history": true,
    "user_expressions": {},
    "allow_stdin": true,
    "stop_on_error": true
}
```

The kernel responds with `execute_reply` on Shell + many `stream`, `display_data`, `execute_result` messages on IOPub.

### Why this matters for your team

When you maintain a notebook backend at scale, you care about:
- How quickly can we provision a kernel? (Pod creation latency, image pull time.)
- How do we route a user's WebSocket to their kernel pod? (Routing layer, sticky sessions.)
- How do we limit a runaway kernel? (cgroup memory limits, idle timeout, hard timeout.)
- How do we share notebooks while keeping kernels isolated? (Shared notebook content; per-user kernel.)
- How do we connect a kernel to PySpark to query Datadog data?

We'll build a tiny version in `projects/mini-jupyter-kernel/`.

---

<a name="module-85--docker"></a>
## Module 8.5 — Docker, properly

> 📖 **Definition — Docker:** A platform for packaging an application and its dependencies into a *container* — a lightweight, isolated environment that runs identically anywhere Docker runs.

### What's in a container

Linux containers use kernel features:
- **Namespaces:** isolate views of the system (PID namespace = process can only see its own children; network namespace = its own network stack; etc.)
- **cgroups:** resource limits (CPU, memory, IO).
- **Layered filesystem (overlay2):** the image is built up of layers, stacked.

A container is NOT a VM. There's no separate kernel. All containers on a host share the host's kernel. They're isolated processes with extra walls.

### Dockerfile basics

```dockerfile
FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

- `FROM` — base image
- `WORKDIR` — set cwd
- `COPY` — copy from build context
- `RUN` — execute at build time, creates a new layer
- `EXPOSE` — documents which port (informational)
- `CMD` — what runs when the container starts

### Layer caching is the pattern

Each `RUN`/`COPY`/`ADD` creates a layer. Docker caches layers; if the layer's inputs haven't changed, it reuses. Order layers from least-changing to most:

```dockerfile
# Bad — every code change reinstalls deps
COPY . .
RUN pip install -r requirements.txt

# Good — code changes don't bust deps cache
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
```

### Multi-stage builds

```dockerfile
# Stage 1: build
FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/server

# Stage 2: runtime (tiny)
FROM gcr.io/distroless/static
COPY --from=builder /out/app /app
ENTRYPOINT ["/app"]
```

Final image has just the binary. ~10 MB instead of 1 GB. Smaller = faster pull = faster scaling.

### Useful commands

```bash
docker build -t myapp:0.1 .
docker run -d -p 8080:8080 --name myapp myapp:0.1
docker ps
docker logs -f myapp
docker exec -it myapp /bin/sh
docker stop myapp && docker rm myapp
docker images
docker system prune -a       # nuke everything (careful)
```

---

<a name="module-86--k8s"></a>
## Module 8.6 — Kubernetes essentials

> 📖 **Definition — Kubernetes (k8s):** A container orchestrator. You declare "I want N copies of this image running, with these resources, behind this load balancer," and k8s makes it happen, keeps it running, scales it, restarts on failure.

### The core objects

- **Pod** — the smallest unit. One or more containers running together (sharing network, sometimes volumes). 99% of pods have one container.
- **Deployment** — a controller that manages a set of pods, ensures N replicas exist, handles rolling updates.
- **Service** — a stable virtual IP + DNS name in front of a set of pods. Balances traffic.
- **Ingress** — exposes services to outside the cluster (typically HTTP routing).
- **ConfigMap** — non-secret config (env vars, files).
- **Secret** — secret config (passwords, API keys, certs).
- **PersistentVolume / PersistentVolumeClaim** — durable storage.
- **Namespace** — virtual cluster within a cluster. Auth boundaries.
- **StatefulSet** — like Deployment but with stable identities/storage. For databases, kafka, etc.
- **DaemonSet** — one pod per node. Used for node-level agents (Datadog Agent runs as a DaemonSet).
- **Job / CronJob** — run-to-completion workloads.

### A simple deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: webapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: webapp
  template:
    metadata:
      labels:
        app: webapp
    spec:
      containers:
        - name: web
          image: kvng/webapp:0.1
          ports:
            - containerPort: 8080
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 256Mi
          readinessProbe:
            httpGet: { path: /healthz, port: 8080 }
            initialDelaySeconds: 5
          livenessProbe:
            httpGet: { path: /healthz, port: 8080 }
            initialDelaySeconds: 30
            periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: webapp
spec:
  selector:
    app: webapp
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
```

```bash
kubectl apply -f deployment.yaml
kubectl get pods
kubectl logs -f deployment/webapp
kubectl describe pod <pod-name>
kubectl exec -it <pod> -- /bin/sh
kubectl rollout status deployment/webapp
kubectl rollout undo deployment/webapp
```

### Resources, requests, limits

- **Request** = "give me at least this much, please." Scheduler uses it to pick a node.
- **Limit** = "don't go above this." Hard cap. Pods over their memory limit get OOMKilled.

CPU is `1` = 1 core; `500m` = 0.5 cores. Memory in `Mi` (mebibyte) or `Gi`.

### Probes

- **Readiness probe** — am I ready to serve traffic? Failed → removed from Service, but pod stays alive.
- **Liveness probe** — am I alive? Failed → pod is restarted.
- **Startup probe** — for slow-starting apps; suspends liveness checks until passing.

Misconfigured probes are a top-3 cause of pod issues. Liveness probe too aggressive → restart loops. Readiness too lax → traffic to broken pods.

### Local k8s with k3d

```bash
brew install k3d
k3d cluster create dev --servers 1 --agents 2
kubectl get nodes
# you have a local 3-node cluster.
```

---

<a name="module-87--operators"></a>
## Module 8.7 — k8s operators and CRDs

> 📖 **Definition — Custom Resource Definition (CRD):** A way to extend k8s with your own object types. Define `kind: Notebook`, write a controller that watches for them, and now `kubectl get notebooks` works.

> 📖 **Definition — Operator:** A controller that manages an application's lifecycle in k8s — installation, upgrades, backup, scaling — using CRDs to expose it declaratively.

Used everywhere:
- Postgres operator → manages Postgres clusters
- Cert-manager → automates TLS cert issuance via Let's Encrypt
- Prometheus operator → manages Prometheus instances
- Strimzi → Kafka operator
- Datadog's own operator → installs/configures the Datadog Agent

For ADP-Notebooks: imagine a `Notebook` CRD. When a user creates a notebook, the operator:
- Creates a kernel Pod
- Sets up a Service routing
- Mounts the notebook content from S3
- Watches for idle timeout, scales to zero
- Watches for delete, cleans up

That's a real shape for a notebook backend. We'll sketch one in `projects/mini-adp-notebooks/` later in Phase 9.

---

<a name="module-88--observability"></a>
## Module 8.8 — Observability culture: RED, USE, SLI/SLO/SLA

This is Datadog's home turf. You should be fluent in this vocabulary before day one.

### The metric methods

#### RED method (for services)

For every service, monitor:
- **Rate** — requests per second
- **Errors** — error rate
- **Duration** — latency distribution (p50, p90, p99)

That's the basic SRE dashboard for any service. Three queries; covers 90% of "is it healthy?"

#### USE method (for resources)

For every resource (CPU, memory, disk, network):
- **Utilization** — % time it's busy
- **Saturation** — queue depth / waiting work
- **Errors** — error count

A CPU at 100% with low saturation = busy but coping. A CPU at 70% with deep saturation = behind on work.

### SLI, SLO, SLA — the trio

> 📖 **Definition — SLI (Service Level Indicator):** A measurement. e.g., "p99 latency for /api/tasks GET." Just a number.

> 📖 **Definition — SLO (Service Level Objective):** A target for an SLI. e.g., "p99 latency < 200 ms, 99.9% of the time per month."

> 📖 **Definition — SLA (Service Level Agreement):** A contractual commitment to customers. Includes consequences if violated. e.g., "99.95% uptime monthly or 10% credit."

SLI < SLO < SLA. The SLA is what you tell customers. The SLO is what you target internally (with a buffer). The SLI is the measurement.

### Error budget

If your SLO is 99.9% (meaning 0.1% of requests can fail), that 0.1% is your **error budget** for the period.

- Burning slowly? Keep shipping.
- Burning fast? Slow down. More testing. Pause feature work.
- Burned through it? Freeze; reliability work only until next period.

This makes "should we deploy on Friday?" objective rather than a feeling.

### The four golden signals (Google SRE)

(Reprised from Phase 5.)
1. Latency
2. Traffic
3. Errors
4. Saturation

Same idea as RED + a saturation knob.

---

<a name="module-89--otel"></a>
## Module 8.9 — OpenTelemetry deeper

We hit OTel in Phase 5. Going deeper here.

### The data model

OTel has three signal types:
- **Traces** — distributed traces (spans, parent-child relationships).
- **Metrics** — time-series numerics.
- **Logs** — structured events.

The genius: **all three share the same context** (trace ID, span ID, attributes). One request's trace, metrics, and logs are all queryable together by trace ID.

### Spans deeper

```
Request comes in. Root span: "POST /tasks"
├── child: "auth.verify_token"  (12ms)
├── child: "db.tx"
│   ├── child: "db.query insert tasks"  (8ms)
│   └── child: "db.commit"  (3ms)
└── child: "kafka.publish task_created"  (5ms)
Total: 35ms
```

Each span has: name, parent_id, start time, end time, attributes, status, events. Events are point-in-time markers within a span ("cache miss", "retrying").

### Trace propagation

Across service boundaries, the trace context goes in headers:
```
traceparent: 00-<trace-id>-<span-id>-01
tracestate: ...
```

When service A calls service B, A injects these headers; B extracts them; B's spans become children of A's span. End-to-end visibility.

### Sampling

At 1M req/sec, you can't store every span. Sampling:
- **Head-based:** decide at the start. Keep 1%. Fast but you might miss the rare important ones.
- **Tail-based:** buffer all spans for a request, decide after (e.g., "always keep errors, slow ones, and 1% of normal").

Datadog's APM defaults to head-based but supports tail-based for Pro tier customers.

### Datadog and OTel

Datadog ingests OTel-formatted spans/metrics/logs natively. You can use the OTel SDK and ship to Datadog without their proprietary library. Or use `dd-trace` (their library) which uses OTel under the hood now.

---

<a name="module-810--datadog-product"></a>
## Module 8.10 — Datadog the product, mental model

You should know what Datadog *does* before joining. Here's the map:

### The product surfaces

- **Infrastructure Monitoring** — host metrics, container metrics, "is my server healthy?"
- **APM (Application Performance Monitoring)** — distributed tracing across your services.
- **Log Management** — centralized log ingestion, search, indexing.
- **Real User Monitoring (RUM)** — frontend perf and error tracking from real browsers.
- **Synthetics** — automated probe tests (browser + API) from Datadog's POPs.
- **Network Performance Monitoring (NPM)** — tcp-level visibility into service-to-service traffic.
- **Database Monitoring (DBM)** — slow query analysis, blocking, etc.
- **Security (Cloud Security, App Security)** — CSPM, threat detection, IAST.
- **CI Visibility** — flaky tests, build performance.
- **Notebooks (your team!)** — interactive analyses against all of the above.
- **Dashboards** — custom visualizations.
- **Monitors / Alerts** — paged when SLOs slip.

The product strategy: ingest *everything* an engineer might need to debug a production system; give them tools to slice it; correlate signals across types.

### The data backend, conceptually

```
Customers' agents/SDKs ────┐
                           ▼
                       Edge / Gateway
                           │
                  ┌────────┼────────┐
                  ▼        ▼        ▼
              Metrics    Traces    Logs
              (Husky)    (?)      (?)
                  │        │        │
                  └────────┼────────┘
                           ▼
                  Query / Aggregation
                           ▼
                   UI / Notebooks
```

You don't need to know exact internals on day 1, but you should know it's many specialized stores not one giant DB.

---

<a name="project-spark"></a>
## 🛠️ Project: Spark word-count + JSON log analysis

Two small Spark jobs. One in Scala, one in PySpark.

**See `projects/spark-jobs/` for code.**

1. **wordcount-scala** — word count on a text file. Demonstrates RDD-style + DataFrame style.
2. **log-analysis-pyspark** — given a JSON-lines log file, compute: requests per minute, top 10 endpoints by p99 latency, error rate per service.

---

<a name="project-jupyter"></a>
## 🛠️ Project: Mini-Jupyter kernel

A tiny "kernel" that speaks the Jupyter protocol over ZeroMQ. Receives Python code, executes, sends back results.

**See `projects/mini-jupyter-kernel/` for code.** Sketch:
- Listen on Shell + IOPub + Heartbeat sockets
- Implement `execute_request`/`execute_reply`
- Use Python's `compile`/`exec` to run code; capture stdout
- Demonstrate by connecting from a Jupyter Notebook frontend

---

<a name="project-k3d"></a>
## 🛠️ Project: Deploy a service to local k3d

Take the TaskAPI from Phase 5, package it as a Docker image, deploy to local k3d with Postgres + Redis as separate services. Add Datadog Agent (or just OTel collector → console) for observability.

**See `projects/k3d-deploy/` for manifests.**

---

<a name="vocab"></a>
## Day-1 vocabulary cheat sheet

Internalize these. They'll come up in your first week.

- **Husky** — Datadog's internal columnar metric storage.
- **Bits AI** — Datadog's LLM-powered assistant for debugging.
- **Watchdog** — Datadog's anomaly-detection product.
- **DBM** — Database Monitoring.
- **APM** — Application Performance Monitoring.
- **CWS** — Cloud Workload Security.
- **CSM** — Cloud Security Management.
- **RUM** — Real User Monitoring.
- **SLO** — Service Level Objective.
- **Toil** — repetitive ops work (vs sustainable engineering work).
- **Page** — get woken up by an alert (verb).
- **Runbook** — the doc you follow when paged.
- **Postmortem** — what we write after an incident.
- **Blameless postmortem** — Datadog culture: focus on systems, not people.
- **Bits** — Datadog's mascot. (Yes, also the AI's namesake.)

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] Spark architecture: driver, executors, DAG, stages, shuffles
- [ ] Read & write Scala Spark and PySpark
- [ ] Jupyter's component model and wire protocol
- [ ] Docker images, layers, multi-stage builds
- [ ] Kubernetes core objects (Deployment, Service, ConfigMap, etc.)
- [ ] CRDs, operators — what they're for
- [ ] RED, USE, SLI/SLO/SLA — fluent
- [ ] OpenTelemetry signals, trace propagation, sampling
- [ ] Datadog product surface map
- [ ] Datadog vocabulary

---

**Next:** [Phase 9 — Capstones](../phase-09-capstones/README.md)
