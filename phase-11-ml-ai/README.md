# Phase 11 — ML/AI Engineering

> You've done ML at PwC. You've worked with OpenAI Codex Lab. StackSense is an AI gateway. So you know more ML than most engineers your age. What this phase covers is the **engineering side** of AI/ML — the part that's usually missing from a CS curriculum: how models actually get served, how you keep them trained correctly, how AI systems behave in production. The kind of work that bridges your ML background with your SWE trajectory.

**Time:** 5–8 days.

**You'll know you're done when:** you can explain MLOps end-to-end, know what a feature store is and why, can build a prod-ready LLM service with rate limiting and observability, and understand vector DBs deeply enough to choose between them.

---

## Table of contents

1. [What does this even mean? — "ML engineering" vs "data science"](#what)
2. [Module 11.1 — The full ML lifecycle](#lifecycle)
3. [Module 11.2 — MLOps — what it really is](#mlops)
4. [Module 11.3 — Model serving: batch, real-time, streaming](#serving)
5. [Module 11.4 — Feature stores](#feature-stores)
6. [Module 11.5 — Vector databases & RAG](#vector-dbs)
7. [Module 11.6 — LLM serving in production](#llm-serving)
8. [Module 11.7 — Evaluation, drift, monitoring](#eval)
9. [Module 11.8 — Cost & latency optimization](#cost)
10. [🛠️ Project: Production LLM service](#project)
11. [Interview question bank](#interview)
12. [What you should now know](#what-you-should-now-know)

---

<a name="what"></a>
## 🧠 What does this even mean? — "ML engineering" vs "data science"

A confusing-but-useful axis:

- **Data scientist:** explores data, trains models, builds prototypes. Often Python in notebooks. Output: a trained model + a story.
- **ML engineer:** takes a model and ships it. Builds the pipelines, the serving infra, the monitoring, the retraining loops. Output: a system that keeps working.

Real teams overlap heavily, but the engineering side is what separates "we trained this thing" from "we have a production AI product."

What ML engineers actually do day-to-day:
- Build feature pipelines (data engineering)
- Build training pipelines (orchestration: Airflow, Kubeflow)
- Build serving infrastructure (FastAPI, KServe, Triton)
- Monitor model behavior (latency, accuracy drift)
- Build experimentation platforms (A/B test infra)
- Write Python that doesn't fall over

You already do most of this in StackSense (an AI gateway IS a model-serving / observability layer). This phase makes that knowledge explicit.

---

<a name="lifecycle"></a>
## Module 11.1 — The full ML lifecycle

```
1. PROBLEM FRAMING
   - What's the metric? What's the baseline? What's the cost of error?

2. DATA
   - Source it (logs, events, labeled data)
   - Validate it (schema, ranges, missingness)
   - Label it (Snorkel, MTurk, internal labelers)
   - Store it (data lake / warehouse)

3. FEATURES
   - Engineer them (rolling means, ratios, embeddings)
   - Test them (data leakage? distribution shift?)
   - Serve them (training-serving skew is the #1 ML bug)

4. TRAIN
   - Pick a model class
   - Hyperparameter search
   - Cross-validate
   - Track experiments (MLflow, W&B)

5. EVALUATE
   - Test set, held out from day 1
   - Cohort metrics (segments)
   - Failure analysis: where does it lose?

6. DEPLOY
   - Package (Docker, ONNX, TorchScript)
   - Serve (REST/gRPC, batched, GPU/CPU)
   - Canary / shadow / A/B

7. MONITOR
   - Latency, throughput, error rate
   - Input distribution drift
   - Output distribution drift
   - Business metric (ground truth lags)

8. RETRAIN
   - Schedule or trigger
   - Old version archived; new version validated; staged rollout
```

Most ML "failures" in production aren't model failures — they're step 3, 6, or 7 failures. Engineering, not modeling.

---

<a name="mlops"></a>
## Module 11.2 — MLOps — what it really is

> 📖 **Definition — MLOps:** DevOps for ML. The operational practices that turn one-off models into versioned, monitored, automatable systems.

The four pillars:

### 1. Versioning everything

Code → Git. Data → DVC, LakeFS, or just object-store with content hashes. Models → MLflow, Weights & Biases, or your registry of choice. Hyperparameters → tracked alongside.

If you can't reproduce a model from versioned inputs, you can't trust it. Period.

### 2. Continuous training pipelines

Treat training as a deployable artifact. Code change → trigger pipeline → train → evaluate → register if good → deploy via canary.

Tools:
- **Airflow / Prefect / Dagster** — DAG orchestration
- **Kubeflow / Metaflow** — ML-specific pipelines
- **GitHub Actions / GitLab CI** — for the simpler cases

### 3. Model registry & serving

A registry stores `(name, version, metadata, artifact_path)` per model. Serving infra pulls a specific version. You can pin, roll back, A/B test.

### 4. Monitoring + feedback loops

You can't just "ship and forget." Models degrade — input distributions change (your users change), real world changes (recession, new product). Without monitoring, you don't know.

The full loop: monitoring → label some recent data → retrain → re-deploy.

---

<a name="serving"></a>
## Module 11.3 — Model serving: batch, real-time, streaming

Three modes:

### Batch

Run inference on a backlog. Output to a table.

Use when: predictions can be precomputed (e.g., "top recommended items per user, refreshed daily"). Cheapest by far.

Stack: Spark / Flink for the big-data version; just Python+pandas for small.

### Real-time (online)

User sends a request, model responds with prediction. Latency-critical.

Use when: prediction needs current input (e.g., "is this transaction fraud?", "what should this LLM say?").

Stack:
- **FastAPI / Flask** — for Python models, simple
- **TorchServe / TensorFlow Serving** — purpose-built for PyTorch / TF
- **Triton Inference Server** — NVIDIA's, GPU-optimized, multi-framework
- **KServe** — Kubernetes-native; production
- **vLLM, TGI** — for LLMs specifically

### Streaming

Predictions on a continuous stream of events. Output to a topic / db.

Use when: things flow continuously (clickstream, sensor data) and decisions need to react quickly.

Stack: Flink, Kafka Streams, Spark Structured Streaming.

### Choosing

| Need | Mode |
|---|---|
| "Show top recommendations on home page" | Batch (precompute hourly) |
| "Score this fraud signal NOW" | Real-time |
| "Continuously detect anomalies in metrics" | Streaming |

Most production AI is more batch than people think. Real-time is glamorous but expensive.

---

<a name="feature-stores"></a>
## Module 11.4 — Feature stores

> 📖 **Definition — Feature store:** A specialized DB that serves *features* (the inputs to a model) for both training (batch, historical) and serving (real-time, current).

The core problem they solve: **training-serving skew.**

You train your fraud model with features computed from `events_history` table — "user's avg transaction value over last 7 days." Then in production you compute the same feature on the fly and call the model. Subtle bug: maybe your training feature was inclusive of "today" but your serving feature isn't. Model trained on slightly-different inputs than it sees in prod. Quality drops mysteriously.

Feature stores solve this by **defining features once**, then materializing them to:
- A batch store (data warehouse) for training jobs
- An online store (Redis, Cassandra, DynamoDB) for serving

Both are computed from the same definition, so they match.

### Tools

- **Feast** — open-source, popular
- **Tecton** — managed, Feast originated there
- Cloud-native: Vertex AI Feature Store, Sagemaker Feature Store

---

<a name="vector-dbs"></a>
## Module 11.5 — Vector databases & RAG

> 📖 **Definition — Vector database:** A database optimized for *similarity search* over high-dimensional vectors. Stores embeddings, returns "find me the K closest to this query vector."

Why: embeddings let you represent semantic content (text, images) as vectors. Similar content → nearby vectors. So "find docs similar to this question" becomes "find vectors near this query embedding."

### How they work

Brute-force similarity is O(N × dims). For 10M docs × 1536 dims, way too slow per query.

Solution: **approximate nearest neighbor (ANN)** algorithms. Trade a tiny bit of accuracy for huge speed.

The dominant algorithm:

> 📖 **Definition — HNSW (Hierarchical Navigable Small World):** A multi-layer graph where each node connects to its nearest neighbors. Search starts at the top (sparse, long-range jumps) and zooms in through layers. Sub-millisecond queries on millions of vectors.

Other algorithms: IVF (inverted file), PQ (product quantization), ScaNN.

### The vector DB landscape

| DB | Strengths | Notes |
|---|---|---|
| **Postgres + pgvector** | Already have Postgres? Add the extension | Up to ~10M vectors comfortably |
| **Pinecone** | Managed, easy | Pricey at scale |
| **Weaviate** | Open source, GraphQL, hybrid search | Strong general choice |
| **Qdrant** | Rust, fast, filtering | Self-host friendly |
| **Milvus** | Cloud-scale, GPU support | Heavy operationally |
| **Chroma** | Local-first, very simple API | Great for prototypes |
| **FAISS** | Library, not a DB | Still useful as a building block |

For most apps under 10M vectors: pgvector. Don't overcomplicate.

### RAG (Retrieval-Augmented Generation)

> 📖 **Definition — RAG:** Pattern where you retrieve relevant context from a vector DB and feed it to an LLM to answer a question. Solves the "LLM doesn't know your private data" problem.

```
1. User asks a question
2. Embed the question (text-embedding-3 or similar)
3. Vector DB query: top K documents most similar to the question
4. Build prompt: "Given these documents [...], answer the question."
5. LLM generates answer
6. Optionally: cite sources
```

Failure modes:
- **Bad retrieval.** Top-K docs aren't actually relevant. Often because chunking is wrong (too small/big) or embeddings don't capture the right semantics.
- **LLM ignores context.** Prompts have to be careful — say "cite the specific document" or "say 'I don't know' if not in the context."
- **Stale data.** RAG depends on a fresh index. Build a re-ingestion pipeline.

You did RAG in your AI Research Agent project. Solidify it.

---

<a name="llm-serving"></a>
## Module 11.6 — LLM serving in production

This is StackSense territory. The key concepts:

### Why LLM serving is unusual

- **GPUs are expensive.** A single inference can need a GPU; you're using one of these even at idle.
- **Output is variable-length.** Unlike traditional ML where output size is fixed, LLMs generate token-by-token until done.
- **Memory is tight.** Models + KV cache fill GPU memory; concurrency is limited.
- **Latency: time-to-first-token (TTFT) and tokens-per-second (TPS).** Different from "latency" in normal services.

### Serving frameworks

- **vLLM** — open source, blazing fast (PagedAttention is the secret sauce)
- **TGI (Text Generation Inference)** — Hugging Face's
- **TensorRT-LLM** — NVIDIA's, hardware-optimal
- **Ollama, llama.cpp** — for local / smaller scale

### Optimizations

- **Continuous batching:** receive new requests while existing ones still generating. Massive throughput gain.
- **PagedAttention:** virtual-memory-style management of the KV cache. Removes fragmentation.
- **Quantization:** lower-precision weights (int8, int4) for smaller memory footprint with small quality loss.
- **Speculative decoding:** small fast model proposes tokens; big model verifies in parallel. ~2x speedup.

### LLM gateways (StackSense itself)

A gateway sits between callers and one or more LLM providers (OpenAI, Anthropic, internal models). Concerns:

- **Provider routing** — choose model based on cost, capability, availability
- **Caching** — exact-match (KV cache) and semantic cache (embedding-similar prompts return cached responses)
- **Rate limiting** — both per-user and aggregated upstream
- **Cost tracking** — token-based billing
- **Observability** — token counts, latency, success rate per provider
- **Fallback / failover** — if OpenAI is down, route to Anthropic
- **Prompt redaction** — strip PII before sending upstream
- **Audit log** — for compliance

You're already building this. The list above is what to make sure you have.

---

<a name="eval"></a>
## Module 11.7 — Evaluation, drift, monitoring

### Offline eval

Test set — labeled, held out, never trained on. Compute metrics:
- **Classification:** precision, recall, F1, AUC-ROC, AUC-PR
- **Regression:** MAE, RMSE, MAPE
- **Generative (LLMs):** harder. BLEU/ROUGE rough; LLM-as-judge; human eval gold standard

### Online eval

Once deployed:
- **A/B test** new model vs old. Measure business metric (clicks, conversion, retention).
- **Shadow deployment:** run both, compare, but only use old's output.
- **Canary:** route 1% to new; ramp up.

### Drift monitoring

- **Input drift:** distribution of features changes (PSI, KS test).
- **Output drift:** prediction distribution changes.
- **Concept drift:** the relationship between features and labels changes (recession → fraud patterns shift). Hardest to detect; relies on labeled feedback.

### Tools

- **Evidently** — open source drift detection
- **Arize, Fiddler, Weights & Biases** — managed
- **Datadog APM** — for latency/errors of model service
- **Datadog ML monitoring** — yes, Datadog itself has this

### LLM-specific eval

- **Hallucination rate:** how often does it fabricate?
- **Refusal rate:** how often does it refuse?
- **Toxicity / safety:** evaluated by separate classifier
- **Latency at the prompt-level:** TTFT, total time
- **Cost per request:** tokens × $/token
- **Retrieval recall (for RAG):** did the right doc get retrieved?

---

<a name="cost"></a>
## Module 11.8 — Cost & latency optimization

ML costs add up. Three levers:

### Smaller models

- **Distillation:** train a small model to mimic a big one. Often 90%+ quality at 10x speed.
- **Quantization:** int8/int4 weights. Tiny quality loss, big memory/speed wins.
- **Pruning:** zero out unimportant weights.

### Smarter routing

- Cheap model first; escalate hard queries to big model
- Cache aggressively; semantic cache catches near-duplicates
- Don't call LLM at all when a deterministic rule works

### Better infra

- Continuous batching → ~3-5x throughput per GPU
- Right-size GPUs (A10 vs A100 vs H100)
- Spot instances for batch / training
- Co-locate model + retrieval to avoid round-trips

For LLM costs specifically:
- **Input tokens cost ~3-10x less than output tokens.** Cap output length aggressively.
- **System prompts get repeated** — providers offer prompt caching to dedupe.
- **Streaming reduces perceived latency** — start showing tokens at TTFT, even if total takes longer.

---

<a name="project"></a>
## 🛠️ Project: Production LLM service

A FastAPI service in front of an LLM provider with the production basics: caching, rate limiting, prompt redaction, observability, fallback.

**See `projects/llm-gateway/`.**

### Spec

- POST `/chat` — takes a prompt, returns a completion
- Auth: API key
- Per-key rate limit (token-bucket via Redis)
- PII redaction (regex first, optional spaCy for stronger)
- Exact-match prompt cache (Redis with TTL)
- Provider fallback: OpenAI primary, Anthropic backup
- Token-count metering per request
- Prometheus metrics: requests, latency, cache hit rate, tokens
- OpenTelemetry traces

This is StackSense scaffolding without StackSense's depth. Use it to feel the parts.

---

<a name="interview"></a>
## 🎯 Interview question bank

1. **Walk through the ML lifecycle. Where do most production models fail?**
2. **What's training-serving skew? How do you prevent it?**
3. **Difference between batch, real-time, streaming inference?**
4. **What is a feature store and what problem does it solve?**
5. **Explain RAG. What's the most common failure mode?**
6. **What's vector similarity search? What is HNSW?**
7. **When would you use pgvector vs Pinecone vs Weaviate?**
8. **How does continuous batching speed up LLM serving?**
9. **What's quantization? When are the quality trade-offs acceptable?**
10. **You launched a model and the metric started drifting. Walk me through diagnosis.**
11. **Design an LLM gateway like StackSense.**
12. **How do you A/B test a model in production?**
13. **You have one cheap model and one expensive model. How do you decide which to use per request?**
14. **What metrics do you monitor for an LLM in production?**
15. **What's MLOps and how does it differ from DevOps?**

---

## ✅ What you should now know

- [ ] The full ML lifecycle and where it fails
- [ ] What MLOps actually means
- [ ] Batch vs real-time vs streaming serving
- [ ] Feature stores and why they exist
- [ ] Vector DBs, HNSW, the vector DB landscape
- [ ] RAG end to end
- [ ] LLM serving optimizations (vLLM, batching, quantization)
- [ ] LLM gateway concerns (caching, fallback, redaction)
- [ ] Drift monitoring and online eval
- [ ] Cost optimization knobs
- [ ] You've built a production-style LLM service

---

**Next:** [Phase 12 — Frontend Depth](../phase-12-frontend-depth/README.md)
