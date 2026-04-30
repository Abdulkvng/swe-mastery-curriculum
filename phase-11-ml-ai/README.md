# Phase 11 — ML/AI Engineering

> You've done ML at PwC. You've worked with OpenAI Codex Lab. StackSense is an AI gateway. So you know more ML than most engineers your age. What this phase covers is the **engineering side** of AI/ML — the part that's usually missing from a CS curriculum: how models actually get served, how you keep them trained correctly, how AI systems behave in production. The kind of work that bridges your ML background with your SWE trajectory.

**Time:** 5–8 days.

**You'll know you're done when:** you can explain MLOps end-to-end, know what a feature store is and why, can build a prod-ready LLM service with rate limiting and observability, understand how LLMs work under the hood, and understand vector DBs deeply enough to choose between them.

---

## Table of contents

1. [What does this even mean? — "ML engineering" vs "data science"](#what)
2. [Module 11.1 — The full ML lifecycle](#lifecycle)
3. [Module 11.2 — MLOps — what it really is](#mlops)
4. [Module 11.3 — Model serving: batch, real-time, streaming](#serving)
5. [Module 11.4 — Feature stores](#feature-stores)
6. [Module 11.5 — Vector databases & RAG](#vector-dbs)
7. [Module 11.6 — LLMs under the hood](#llm-under-the-hood)
8. [Module 11.7 — LLM serving in production](#llm-serving)
9. [Module 11.8 — Evaluation, drift, monitoring](#eval)
10. [Module 11.9 — Cost & latency optimization](#cost)
11. [🛠️ Project: Production LLM service](#project)
12. [Quiz: LLMs under the hood](#llm-quiz)
13. [Interview question bank](#interview)
14. [What you should now know](#what-you-should-now-know)

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

<a name="llm-under-the-hood"></a>
## Module 11.6 — LLMs under the hood

LLMs look magical, but under the hood they are mostly doing one job:

> Given the tokens I have already seen, predict the most likely next token.

That sounds simple, but it becomes powerful when the model has billions of parameters, sees huge amounts of text during training, and uses the transformer architecture to understand context.

### The simple flow

```text
Text input
  ↓
Tokenizer
  ↓
Token IDs
  ↓
Embeddings
  ↓
Transformer engine
  ↓
Probability distribution over next tokens
  ↓
Next token
  ↓
Repeat until done
```

Example:

```text
Prompt: "The capital of France is"
Model predicts likely next tokens:
- " Paris"  → high probability
- " London" → low probability
- " pizza"  → very low probability
```

The model does not output the whole answer at once. It generates one token, adds it to the context, then predicts the next token again.

### Step 1: Tokenization

The model does not see raw English words. It sees tokens.

A token can be:
- A whole word
- Part of a word
- A space plus a word
- A symbol
- A piece of code

Example:

```text
"LLMs are powerful"
```

Might become something like:

```text
["LL", "Ms", " are", " powerful"]
```

Then each token gets converted into an ID:

```text
[2043, 8172, 389, 6615]
```

The exact IDs depend on the tokenizer.

Why this matters:
- Cost is usually based on tokens.
- Long prompts cost more.
- Long prompts are slower.
- Tokenization can make some words, names, and code more expensive than expected.

### Step 2: Embeddings

Token IDs are not useful by themselves. The model turns each token ID into a vector.

A vector is just a list of numbers.

Example idea:

```text
"king"  → [0.21, -0.45, 0.88, ...]
"queen" → [0.19, -0.39, 0.91, ...]
"banana" → [-0.72, 0.12, 0.03, ...]
```

The point is not that the numbers are meaningful to humans. The point is that the model learns useful relationships between tokens.

Tokens used in similar contexts end up with similar vector patterns.

### Step 3: The transformer engine

The transformer is the engine that makes modern LLMs work.

Before transformers, models struggled with long context. They had trouble connecting words that were far apart. Transformers solved this with attention.

The transformer asks:

> For each token, which other tokens in the prompt should I pay attention to?

Example:

```text
The animal did not cross the road because it was tired.
```

The word `it` probably refers to `animal`, not `road`.

Attention helps the model learn that relationship.

### Attention in simple terms

Each token creates three things:

- **Query:** what am I looking for?
- **Key:** what information do I have?
- **Value:** what should I pass forward if someone pays attention to me?

Simple analogy:

```text
Query = question
Key = label on each piece of information
Value = the actual information
```

The model compares queries and keys to decide what matters.

If a token's query matches another token's key strongly, the model pays more attention to that token's value.

That is how the model connects ideas across the prompt.

### Multi-head attention

The model does not use only one attention pattern. It uses many at the same time.

Each attention head can focus on something different:
- Grammar
- Subject/object relationships
- Code structure
- Names
- Dates
- Earlier instructions
- The tone of the prompt

This is why the model can track multiple things at once.

### Transformer layers

An LLM has many transformer layers stacked on top of each other.

A rough mental model:

```text
Layer 1: basic token patterns
Layer 2: grammar and nearby relationships
Layer 3: sentence meaning
Layer 4+: broader reasoning and task structure
Later layers: answer planning, style, and next-token prediction
```

That is not a perfect scientific breakdown, but it is a good intuition.

Each layer refines the representation of every token.

By the end, the model has a context-aware representation of the prompt.

### Step 4: Predicting the next token

After the transformer processes the prompt, the model produces scores for possible next tokens.

These scores get converted into probabilities.

Example:

```text
Prompt: "The capital of France is"

Possible next tokens:
" Paris"  → 92%
" Lyon"   → 2%
" London" → 1%
other     → 5%
```

The model then chooses a token.

It might choose the most likely token, or it might sample from likely tokens depending on settings like temperature.

### Temperature

Temperature controls randomness.

Low temperature:
- More predictable
- More focused
- Better for factual answers and code

High temperature:
- More creative
- More varied
- Higher chance of weird answers

Simple version:

```text
temperature = 0   → safest / most deterministic
temperature = 0.7 → balanced
temperature = 1+  → more random
```

### Step 5: Training

LLMs learn by reading huge amounts of text and trying to predict missing next tokens.

Training loop:

```text
1. Give model some text
2. Hide the next token
3. Model guesses the next token
4. Compare guess to correct answer
5. Adjust weights slightly
6. Repeat billions or trillions of times
```

At first, the model is terrible. Over time, it learns patterns:
- Grammar
- Facts
- Code syntax
- Reasoning patterns
- Writing styles
- Common explanations
- Cause and effect patterns

It is not storing the internet like a database. It is compressing patterns into weights.

### Step 6: Fine-tuning and alignment

Base models are trained to predict text. That does not automatically make them helpful assistants.

After pretraining, models are usually improved with:

- **Instruction tuning:** train the model to follow instructions.
- **Preference tuning / RLHF:** train the model to prefer answers humans rate as better.
- **Safety tuning:** reduce harmful, low-quality, or policy-breaking behavior.

This is why a chatbot answers directly instead of just continuing your sentence randomly.

### What the model actually "knows"

An LLM does not know things the way a human does.

Better mental model:

> An LLM is a giant pattern engine trained to predict text, but the patterns are rich enough that it can summarize, reason, code, translate, and explain.

It can still be wrong because it is predicting likely text, not checking truth by default.

That is why RAG, tools, evals, and citations matter in production systems.

### Why this matters for StackSense and AI infra

Under the hood details directly affect system design:

| LLM detail | Infra impact |
|---|---|
| Tokens drive computation | Track token usage and cost |
| Longer prompts are slower | Optimize prompts and context |
| Output is generated token by token | Measure TTFT and tokens/sec |
| Context window is limited | Use RAG and memory carefully |
| Models can hallucinate | Add evals, retrieval, and guardrails |
| Providers fail or slow down | Build routing and fallback |
| Different models have different strengths | Route tasks by cost, speed, and quality |

This is why an LLM gateway like StackSense is useful. It gives teams control over cost, latency, quality, fallback, and observability.

---

<a name="llm-serving"></a>
## Module 11.7 — LLM serving in production

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
## Module 11.8 — Evaluation, drift, monitoring

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
## Module 11.9 — Cost & latency optimization

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

### Bonus project: Build a tiny next-token predictor

Build a small toy model that predicts the next character or word from a short text file.

Goal: understand the core LLM idea without needing GPUs.

Requirements:
- Use Python.
- Load a small `.txt` file.
- Split text into characters or simple word tokens.
- Train a tiny model or use a simple Markov chain.
- Given a seed phrase, generate 50-100 tokens.
- Print the top 5 likely next tokens at each step.

What you should learn:
- Tokenization matters.
- Prediction happens one step at a time.
- Better context usually improves output.
- Randomness changes the generated answer.
- Even a simple predictor can produce patterns, but not deep understanding.

Stretch:
- Add a `temperature` setting.
- Compare low temperature vs high temperature outputs.
- Track generated token count and estimated cost like a mini StackSense.

---

<a name="llm-quiz"></a>
## Quiz: LLMs under the hood

### Concept checks

1. What is the main job an LLM is trained to do?
2. Why does a model use tokens instead of raw text?
3. What is an embedding?
4. What problem did transformers help solve compared to older sequence models?
5. In attention, what are query, key, and value in plain English?
6. Why does multi-head attention help?
7. Why does an LLM generate answers one token at a time?
8. What does temperature control?
9. Why can an LLM hallucinate?
10. Why do token counts matter for infra and cost?

### Short answers

1. Explain the difference between pretraining and instruction tuning.
2. Explain why a longer prompt usually costs more and runs slower.
3. Explain why RAG helps with private or fresh data.
4. Explain why LLM latency is often measured with TTFT and tokens/sec.
5. Explain how StackSense could use token tracking to save money.

### Build-thinking questions

1. If users complain your chatbot is slow, what metrics would you inspect first?
2. If your LLM bill suddenly doubles, what would you check?
3. If the model keeps making up facts, what changes would you make to the system?
4. If one provider goes down, how should an LLM gateway respond?
5. If a small model is good enough for easy prompts, how would you route requests?

### Answer key

1. Predict the next token from previous tokens.
2. Models need numerical inputs. Tokenization converts text into manageable pieces with IDs.
3. A learned vector representation of a token, phrase, document, or other input.
4. They made it easier to connect far-apart tokens using attention.
5. Query is what a token is looking for, key is what another token offers, value is the information passed forward.
6. Different heads can focus on different relationships at the same time.
7. Each generated token becomes part of the context for predicting the next token.
8. Randomness in token selection.
9. It predicts likely text and does not automatically verify truth.
10. Tokens drive cost, latency, context limits, and provider billing.

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
8. **Explain how an LLM works under the hood.**
9. **What problem does attention solve in transformers?**
10. **Why do LLMs generate token by token?**
11. **How does continuous batching speed up LLM serving?**
12. **What's quantization? When are the quality trade-offs acceptable?**
13. **You launched a model and the metric started drifting. Walk me through diagnosis.**
14. **Design an LLM gateway like StackSense.**
15. **How do you A/B test a model in production?**
16. **You have one cheap model and one expensive model. How do you decide which to use per request?**
17. **What metrics do you monitor for an LLM in production?**
18. **What's MLOps and how does it differ from DevOps?**

---

## ✅ What you should now know

- [ ] The full ML lifecycle and where it fails
- [ ] What MLOps actually means
- [ ] Batch vs real-time vs streaming serving
- [ ] Feature stores and why they exist
- [ ] Vector DBs, HNSW, the vector DB landscape
- [ ] RAG end to end
- [ ] How LLMs work under the hood
- [ ] Tokens, embeddings, attention, transformer layers, and next-token prediction
- [ ] LLM serving optimizations (vLLM, batching, quantization)
- [ ] LLM gateway concerns (caching, fallback, redaction)
- [ ] Drift monitoring and online eval
- [ ] Cost optimization knobs
- [ ] You've built a production-style LLM service

---

**Next:** [Phase 12 — Frontend Depth](../phase-12-frontend-depth/README.md)
