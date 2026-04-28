# Phase 5 — APIs & Backend

> An API is the contract between two pieces of software. Most of your career will be spent designing them, calling them, debugging them, securing them, and figuring out why they got slow. This phase turns "I can write a Flask endpoint" into "I can architect an authenticated, rate-limited, instrumented backend with proper RPC contracts and graceful failure handling."
>
> We'll build the same backend twice: once as REST/JSON in TypeScript, once as gRPC over protocol buffers. You'll see why teams pick one over the other. Then we layer on JWT auth, Redis-backed rate limiting, structured logging, and a Postgres-backed data layer (using the connection-pool ideas from Phase 2 + Postgres knowledge from Phase 4).

**Time:** 7–10 days.

**You'll know you're done when:** you can design a clean REST API given a use case, defend "REST vs gRPC" with concrete trade-offs, build a production-style backend with auth and rate limiting, and explain idempotency, retries, and timeouts well enough to pass a Datadog systems-design round.

---

## Table of contents

1. [What does this even mean? — "API"](#what-api-means)
2. [Module 5.1 — TypeScript and Node, just enough](#module-51--ts-node)
3. [Module 5.2 — REST: principles, not dogma](#module-52--rest)
4. [Module 5.3 — Designing a REST API: the rules of thumb](#module-53--rest-design)
5. [Module 5.4 — JSON Schema, validation, and OpenAPI](#module-54--openapi)
6. [Module 5.5 — Authentication: sessions, JWT, OAuth, API keys](#module-55--auth)
7. [Module 5.6 — Authorization: RBAC, ABAC, the difference](#module-56--authz)
8. [Module 5.7 — Rate limiting, idempotency, retries, timeouts](#module-57--reliability)
9. [Module 5.8 — gRPC and protobuf](#module-58--grpc)
10. [Module 5.9 — When REST, when gRPC, when GraphQL](#module-59--when-what)
11. [Module 5.10 — Observability for APIs](#module-510--observability)
12. [🛠️ Project: TaskAPI — a real backend](#project-taskapi)
13. [Exercises](#exercises)
14. [Interview question bank](#interview-questions)
15. [What you should now know](#what-you-should-now-know)

---

<a name="what-api-means"></a>
## 🧠 What does this even mean? — "API"

> 📖 **Definition — API (Application Programming Interface):** A contract that says "if you call me this way, I'll respond that way."

The word "API" is overloaded. It can mean:

- **Library API** — the methods exposed by a library you import (`json.Marshal()`).
- **Network API** — endpoints exposed over HTTP/gRPC/etc that other services or clients call (`POST /users`).
- **Operating system API** — system calls (`open()`, `read()`, `fork()`).

When SWEs say "API" without qualifier, they usually mean network API. That's what this phase is about.

The point of an API isn't the protocol. It's the **contract**. A good API:
- Says exactly what it accepts and returns
- Says exactly how it can fail
- Doesn't change in ways that break existing callers
- Is hard to misuse

Most production bugs come from APIs that did one of those poorly.

---

<a name="module-51--ts-node"></a>
## Module 5.1 — TypeScript and Node, just enough

We're using TypeScript + Node + Express because it's what you'll see at most non-Big-Tech backends, what's easy to build & demo, and what your StackSense work probably uses for its dashboard. Skip if you're already fluent.

### Setup

```bash
mkdir taskapi && cd taskapi
npm init -y
npm install express pg ioredis jsonwebtoken zod pino pino-http
npm install -D typescript @types/node @types/express @types/jsonwebtoken @types/pg ts-node tsx
npx tsc --init --target es2022 --module nodenext --moduleResolution nodenext \
    --strict --esModuleInterop --resolveJsonModule --outDir dist
```

### Type system that actually helps

```ts
// Inferred types — TS figures out the type from usage.
const name = "Kvng"        // string
const age = 21             // number

// Explicit types
let users: string[] = []
let m: Map<string, number> = new Map()

// Object shapes (interfaces)
interface User {
    id: number
    email: string
    name: string
    deletedAt: Date | null    // union with null
}

// Optional fields
interface Config {
    host: string
    port?: number             // optional (port?: number means port: number | undefined)
}

// Type aliases
type UserId = number
type Result<T> = { ok: true; value: T } | { ok: false; error: string }

// Discriminated unions — your best friend for state machines
type Connection =
    | { state: "idle" }
    | { state: "connecting"; startedAt: Date }
    | { state: "open"; socket: Socket }
    | { state: "closed"; reason: string }

function describe(c: Connection): string {
    switch (c.state) {
        case "idle": return "no connection"
        case "connecting": return `connecting since ${c.startedAt}`
        case "open": return `connected`
        case "closed": return `closed: ${c.reason}`
        // TS warns if you forget a case ✅
    }
}
```

### Async/await

```ts
// async functions return Promises automatically
async function fetchUser(id: number): Promise<User> {
    const res = await fetch(`/users/${id}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    return res.json()
}

// Run in parallel
const [user, posts] = await Promise.all([
    fetchUser(42),
    fetchPosts(42),
])

// Errors propagate via try/catch
try {
    const u = await fetchUser(42)
} catch (e) {
    console.error("oops:", e)
}
```

### Express in 30 seconds

```ts
import express from "express"

const app = express()
app.use(express.json())

app.get("/healthz", (req, res) => {
    res.json({ ok: true })
})

app.post("/users", async (req, res) => {
    const { name, email } = req.body
    // ... validate, insert, etc
    res.status(201).json({ id: 42, name, email })
})

app.listen(3000, () => console.log("listening"))
```

That's enough TS+Express. We'll layer on validation, auth, etc.

---

<a name="module-52--rest"></a>
## Module 5.2 — REST: principles, not dogma

> 📖 **Definition — REST (Representational State Transfer):** An architectural style for distributed systems. Roy Fielding's PhD thesis, 2000. Boils down to: stateless requests, resource-oriented URLs, uniform interface (HTTP methods).

Most APIs called "REST" today are actually "JSON over HTTP that vaguely follows REST conventions." Pure REST (HATEOAS, hypermedia) is rare in practice. That's fine — **most of "REST" that matters is good HTTP + good URLs.**

### The actual constraints

1. **Client-server.** Client and server are separate.
2. **Stateless.** Each request stands alone — no server-side session state needed to interpret it. Scale by adding servers behind a load balancer.
3. **Cacheable.** Responses say whether they can be cached.
4. **Layered system.** Client doesn't know if it's talking to the origin or a proxy/CDN.
5. **Uniform interface.** Resources, methods, representations.

The stateless one is the big one. It's why REST scales horizontally and why session storage moved to Redis: the server processes are interchangeable.

### Resource-oriented URLs

Bad:
```
GET /getUserById?id=42
POST /createUser
POST /updateUser?id=42
```

Good:
```
GET    /users/42        (get one)
GET    /users           (list)
POST   /users           (create)
PUT    /users/42        (replace)
PATCH  /users/42        (partial update)
DELETE /users/42        (delete)
```

URLs are nouns. HTTP methods are verbs.

### Nesting

```
GET  /users/42/posts          (posts by user 42)
POST /users/42/posts          (create a post for user 42)
GET  /posts/99                (get one post directly)
GET  /posts/99/comments       (comments on post 99)
```

Two rules:
- Nest at most 1–2 levels deep. Beyond that, URLs get unreadable. Prefer `/posts/99/comments` over `/users/42/posts/99/comments`.
- Resources should also be reachable at the top level when possible.

---

<a name="module-53--rest-design"></a>
## Module 5.3 — Designing a REST API: the rules of thumb

### Status codes — pick the right one

The full chart from Phase 1 still applies. For APIs specifically:

| Scenario | Status |
|---|---|
| OK, returning data | 200 |
| Created (POST) | 201 |
| Accepted (async, will process) | 202 |
| Deleted, no body | 204 |
| Invalid input (validation failed) | 400 |
| No / bad auth | 401 |
| Auth OK but not allowed | 403 |
| Resource doesn't exist | 404 |
| Conflict (duplicate, version mismatch) | 409 |
| Too much data, can't process semantically | 422 |
| Rate limited | 429 |
| Server bug | 500 |
| Upstream service down | 502 / 503 / 504 |

### Response envelope

There's a 30-year holy war about whether to wrap responses (`{"data": ...}`) or not. Pick one and be consistent. Modern style I'd recommend:

```json
// Success
{
    "id": 42,
    "name": "Kvng",
    "email": "k@x.com"
}

// Error
{
    "error": {
        "code": "validation_failed",
        "message": "email must be a valid address",
        "details": [{ "field": "email", "issue": "format" }]
    }
}
```

The error shape is the important part. Always have:
- A short `code` for clients to switch on
- A human-readable `message` for debugging
- Optional `details` for field-level validation errors

### Pagination

```
GET /posts?limit=20&cursor=eyJpZCI6MTIzfQ
```

Cursor-based ≫ offset-based. Why: offset breaks when items are inserted/deleted concurrently (you skip or duplicate items). A cursor pointing to "the post with id 123" stays correct even if rows shuffle.

```json
{
    "items": [...],
    "next_cursor": "eyJpZCI6MTQzfQ"
}
```

### Filtering, sorting

```
GET /posts?author=42&status=published&sort=-created_at&limit=20
```

The `-` prefix means descending. Be careful with full-text-style query DSLs — they explode in scope. Keep it dumb until you have a clear reason to add complexity.

### Versioning

Three flavors:
- **URL:** `GET /v1/users/42`. Simple, visible, cache-friendly.
- **Header:** `Accept: application/vnd.taskapi.v1+json`. Cleaner URLs but harder to debug.
- **Query param:** `GET /users/42?api_version=1`. Don't.

URL versioning is the boring, correct default. Use semantic versioning sparingly: only major-version bumps when you break compatibility. Adding fields ≠ a new version.

### Idempotency

> 📖 **Definition — Idempotent:** Calling it twice has the same effect as calling it once.

GET, PUT, DELETE are inherently idempotent. POST and PATCH usually aren't.

For payment-like POSTs, **idempotency keys** are the pattern:

```
POST /payments
Idempotency-Key: 7c6fe5ad-9f23-4e8a-9f3d-5c4d3b2a1e0f
```

The server stores `(key → response)`. If a retry comes in with the same key, return the same response without re-charging. Stripe popularized this pattern; everyone uses it now.

---

<a name="module-54--openapi"></a>
## Module 5.4 — JSON Schema, validation, and OpenAPI

You **never** trust input. Every request body, query param, header gets validated.

### Validation in TypeScript with Zod

```ts
import { z } from "zod"

const CreatePostSchema = z.object({
    title: z.string().min(1).max(200),
    body: z.string().max(10_000).optional(),
    tags: z.array(z.string()).max(10).optional(),
})

type CreatePost = z.infer<typeof CreatePostSchema>

app.post("/posts", async (req, res) => {
    const result = CreatePostSchema.safeParse(req.body)
    if (!result.success) {
        return res.status(400).json({
            error: {
                code: "validation_failed",
                message: "invalid request body",
                details: result.error.format(),
            },
        })
    }
    const post: CreatePost = result.data
    // ... insert
})
```

`z.infer<>` gives you a TypeScript type derived from the schema. One source of truth.

### OpenAPI — your API's machine-readable contract

> 📖 **Definition — OpenAPI:** A YAML/JSON spec describing every endpoint, method, request shape, response shape, error case. Tools generate docs, client SDKs, server stubs from it.

```yaml
openapi: 3.1.0
info:
  title: TaskAPI
  version: 1.0.0
paths:
  /tasks:
    get:
      summary: List tasks
      parameters:
        - name: limit
          in: query
          schema: { type: integer, maximum: 100 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TaskList'
    post:
      summary: Create a task
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateTask'
      responses:
        '201':
          description: Created
components:
  schemas:
    Task:
      type: object
      required: [id, title, completed]
      properties:
        id: { type: integer }
        title: { type: string }
        completed: { type: boolean }
    TaskList:
      type: object
      properties:
        items: { type: array, items: { $ref: '#/components/schemas/Task' } }
        next_cursor: { type: string, nullable: true }
    CreateTask:
      type: object
      required: [title]
      properties:
        title: { type: string, minLength: 1, maxLength: 200 }
```

You can:
- Render to a docs site (Swagger UI, Redoc)
- Generate TypeScript types (`openapi-typescript`)
- Generate Python client (`openapi-generator`)
- Mock the server for frontend devs

OpenAPI is the standard. Every serious API ships one.

---

<a name="module-55--auth"></a>
## Module 5.5 — Authentication: sessions, JWT, OAuth, API keys

> 📖 **Definition — Authentication (authn):** Proving who you are. **Authorization (authz):** Proving what you're allowed to do. Two different things; both end with "auth" in casual talk.

### Sessions — the boring traditional way

1. User submits username + password.
2. Server checks → creates a random session ID, stores `{session_id → user_id}` in Redis or DB.
3. Server sets `Set-Cookie: session=<id>; HttpOnly; Secure; SameSite=Lax`.
4. Each subsequent request → cookie is sent → server looks up session.

Pros: easy to invalidate (delete from store), small cookie.
Cons: stateful (need shared session store), every request needs a Redis/DB hit.

### JWT (JSON Web Tokens) — the modern standard for stateless APIs

> 📖 **Definition — JWT:** A self-contained token that contains claims (user ID, roles, expiry, etc.) and is cryptographically signed. The server can verify it without storing anything.

A JWT has three parts, base64-url separated by dots:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI0MiIsImV4cCI6MTcyMTk5OTk5OX0.AbCdE...
[ header                            ].[ payload                              ].[ signature ]
```

Decoded:
```json
// Header
{ "alg": "HS256", "typ": "JWT" }
// Payload (claims)
{ "sub": "42", "exp": 1721999999, "role": "admin" }
// Signature
HMAC_SHA256(base64url(header) + "." + base64url(payload), SECRET)
```

The signature ensures: **anyone with the secret can verify, but no one can forge a new token.**

Verifying: re-compute the signature with your secret. If it matches, the token is genuine. Then check `exp` (expiry).

Pros: stateless. No DB lookup per request. Scales horizontally without shared store.
Cons: **you can't easily revoke a token before expiry** (it's valid until `exp`). Mitigations: short expiry + refresh tokens; a small revocation list in Redis for emergencies.

### Two algorithm classes for JWT

- **Symmetric (HS256):** one secret, used both to sign and verify. Server-internal. Don't share with clients.
- **Asymmetric (RS256, ES256):** sign with a private key, verify with a public key. Useful when one service issues tokens (auth service) and many services verify them — only the auth service needs the private key.

### Implementing JWT in Node

```ts
import jwt from "jsonwebtoken"

const SECRET = process.env.JWT_SECRET!

// Sign
function issue(userId: number, role: string): string {
    return jwt.sign(
        { sub: String(userId), role },
        SECRET,
        { algorithm: "HS256", expiresIn: "15m" }
    )
}

// Verify (in middleware)
function authMiddleware(req: Request, res: Response, next: NextFunction) {
    const auth = req.header("Authorization")
    if (!auth?.startsWith("Bearer ")) {
        return res.status(401).json({ error: { code: "unauthenticated" } })
    }
    try {
        const payload = jwt.verify(auth.slice(7), SECRET) as JwtPayload
        ;(req as any).user = { id: Number(payload.sub), role: payload.role }
        next()
    } catch (e) {
        res.status(401).json({ error: { code: "invalid_token" } })
    }
}
```

### Refresh tokens

Pattern:
- Issue a short-lived **access token** (15 min, JWT)
- Issue a long-lived **refresh token** (30 days, opaque random string, stored in DB)
- Client uses access token for normal requests
- When access token expires, client uses refresh token to get a new access token
- Refresh tokens can be revoked (DB row)

This gives you JWT's scalability with revocation as an escape hatch.

### OAuth 2.0 — when third parties are involved

> 📖 **Definition — OAuth 2.0:** A delegation protocol. "User, please authorize my app to read your Google calendar." Three-legged flow: user, your app, third-party service. End result: your app gets an access token to act on the user's behalf.

You won't implement OAuth from scratch — use a library (Passport.js, Auth0, Clerk, Supabase Auth, etc.). But know:
- **Authorization code flow** with PKCE — modern default for web apps
- **Client credentials flow** — service-to-service, no human user
- The redirect URI is critical (it's how the auth server returns the code to your app)

### API keys

The crudest auth: a long random string. Client sends `Authorization: ApiKey k_abcd1234...`. Server looks up `(key → user/permissions)`.

Use cases: server-to-server calls, programmatic access to your API, internal tools.

Best practices:
- Generate with a CSPRNG (`crypto.randomBytes(32).toString("hex")`)
- Store **only the hash** in your DB (like passwords) — if the DB leaks, attackers can't use the keys
- Show the key to the user **once** at creation
- Allow rotation and revocation
- Optionally scope keys to specific permissions

---

<a name="module-56--authz"></a>
## Module 5.6 — Authorization: RBAC, ABAC, the difference

After you know **who** the user is, decide **what they can do**.

### RBAC (Role-Based Access Control)

Users have **roles**. Roles have **permissions**. Code checks roles.

```ts
const ROLE_PERMISSIONS = {
    admin:  ["users:read", "users:write", "posts:read", "posts:write"],
    editor: ["posts:read", "posts:write"],
    viewer: ["posts:read"],
}

function can(user: User, permission: string): boolean {
    return ROLE_PERMISSIONS[user.role]?.includes(permission) ?? false
}

app.delete("/posts/:id", authMiddleware, (req, res) => {
    if (!can(req.user, "posts:write")) {
        return res.status(403).json({ error: { code: "forbidden" } })
    }
    // ...
})
```

Simple, ubiquitous. Works for most apps.

### ABAC (Attribute-Based Access Control)

Decisions consider **attributes** of the user, the resource, and the context. More flexible, more complex.

Example: "user can edit a post if (a) they're an admin, OR (b) they're the post's author, OR (c) they're a member of the same team and the post is < 1 hour old."

```ts
function canEditPost(user: User, post: Post): boolean {
    if (user.role === "admin") return true
    if (post.authorId === user.id) return true
    if (
        user.teamId === post.teamId &&
        Date.now() - post.createdAt.getTime() < 60 * 60 * 1000
    ) {
        return true
    }
    return false
}
```

ABAC explodes in complexity fast. Tools like OPA (Open Policy Agent) externalize the rules.

### The "ownership" check is everywhere

In practice, most authz is "is this object owned by the user?" Examples:

```ts
const post = await db.post.findById(id)
if (post.authorId !== req.user.id) {
    return res.status(404).json({ error: { code: "not_found" } })
    // ^ return 404, not 403, to avoid leaking that the resource exists
}
```

That last line is subtle but correct: returning 403 leaks that the resource exists. 404 hides it.

---

<a name="module-57--reliability"></a>
## Module 5.7 — Rate limiting, idempotency, retries, timeouts

These are the four reliability primitives. Get them right and your API survives bursts, retries, and bad clients.

### Rate limiting

> 📖 **Definition — Rate limiting:** Capping how many requests a given client can make in a window.

Why: prevent abuse, prevent overload, enforce fair use, save money.

Algorithms:

#### Fixed window

"100 requests per minute" — counter resets at the top of each minute. Simple, but allows bursts at the boundary (200 requests in 2 seconds: 100 in second 59, 100 in second 0).

#### Sliding window log

Keep a list of timestamps; count those in the last 60 seconds. Accurate but expensive (memory).

#### Sliding window counter

Hybrid. Keep current and previous window; weight the previous by how much of it falls in the rolling window. Cheap and reasonably accurate.

#### Token bucket

A bucket holds up to N tokens; refills at R per second. Each request consumes 1 token. Out of tokens → reject. Allows bursts (up to N) but limits sustained rate (R/sec). Used by AWS, Stripe.

#### Leaky bucket

Like token bucket but smooths out bursts. Requests queued; drained at constant rate. Strict rate, may add latency.

### Implementing rate limit in Redis

```ts
// Fixed window — simplest
async function fixedWindowAllow(
    redis: Redis,
    key: string,
    max: number,
    windowSec: number
): Promise<boolean> {
    const current = await redis.incr(key)
    if (current === 1) {
        await redis.expire(key, windowSec)
    }
    return current <= max
}

// Use it
const ok = await fixedWindowAllow(redis, `rl:${userId}:${minute}`, 100, 60)
if (!ok) return res.status(429).set("Retry-After", "60").end()
```

### Idempotency — covered already, key headers

```
Idempotency-Key: <client-generated UUID>
```

Server stores `(key → request hash + response)` for some retention (e.g., 24h). On retry:
- Same key, same payload → return stored response
- Same key, different payload → reject with 409 (collision)

### Retries

When a downstream call fails, should you retry?

Rules:
- **Only retry on idempotent operations** (GET, PUT, DELETE; POSTs only with idempotency keys)
- **Only retry on transient errors** (5xx, network errors, 429). Never on 4xx.
- **Use exponential backoff with jitter:** `delay = base * 2^n + random(0, base)`. Without jitter, all clients sync up and DDoS your service when it recovers ("thundering herd").
- **Cap the retries** (e.g., 3-5). Beyond that, surface the error.
- **Set a deadline** for the entire operation, not just per-attempt.

```ts
async function withRetry<T>(fn: () => Promise<T>, maxAttempts = 3): Promise<T> {
    let lastErr: any
    for (let i = 0; i < maxAttempts; i++) {
        try {
            return await fn()
        } catch (e: any) {
            if (!isRetryable(e)) throw e
            lastErr = e
            const delay = 100 * 2 ** i + Math.random() * 100   // ms
            await new Promise(r => setTimeout(r, delay))
        }
    }
    throw lastErr
}
```

### Timeouts

Every external call needs a timeout. Always.

A request without a timeout can hang forever. Hung requests pile up. Threads/event-loop exhaust. Cascade failure. This is how outages start.

```ts
const controller = new AbortController()
const timeoutId = setTimeout(() => controller.abort(), 3000)
try {
    const res = await fetch(url, { signal: controller.signal })
    // ...
} finally {
    clearTimeout(timeoutId)
}
```

In Go: `context.WithTimeout`. In Java: `CompletableFuture.orTimeout`. Every language has one. Use it.

### Circuit breakers

When a downstream is failing, stop calling it for a while. Recover when it's healthy.

```
states: closed → open → half-open → closed
- closed: normal; if N% fail in window → open
- open: fail fast without calling; after T seconds → half-open
- half-open: allow one trial; success → closed; fail → open
```

Library: `opossum` (Node), `gobreaker` (Go), `Resilience4j` (Java).

---

<a name="module-58--grpc"></a>
## Module 5.8 — gRPC and protobuf

> 📖 **Definition — gRPC:** A binary RPC framework over HTTP/2, using **Protocol Buffers** (protobuf) for message format. Created at Google, open-sourced 2015. Datadog uses it heavily for internal services.

> 📖 **Definition — Protocol Buffers (protobuf):** A schema-first, binary serialization format. You define messages in a `.proto` file, generate code for any language. Smaller and faster than JSON.

### A protobuf example

```proto
// task.proto
syntax = "proto3";

package taskapi;

import "google/protobuf/timestamp.proto";

service TaskService {
    rpc CreateTask(CreateTaskRequest) returns (Task);
    rpc GetTask(GetTaskRequest) returns (Task);
    rpc ListTasks(ListTasksRequest) returns (ListTasksResponse);
    rpc StreamUpdates(StreamUpdatesRequest) returns (stream TaskUpdate);
}

message Task {
    int64 id = 1;
    string title = 2;
    bool completed = 3;
    google.protobuf.Timestamp created_at = 4;
}

message CreateTaskRequest {
    string title = 1;
}

message GetTaskRequest {
    int64 id = 1;
}

message ListTasksRequest {
    int32 limit = 1;
    string cursor = 2;
}

message ListTasksResponse {
    repeated Task items = 1;
    string next_cursor = 2;
}

message StreamUpdatesRequest {}

message TaskUpdate {
    enum Kind {
        KIND_UNKNOWN = 0;
        CREATED = 1;
        UPDATED = 2;
        DELETED = 3;
    }
    Kind kind = 1;
    Task task = 2;
}
```

The numbers (`= 1`, `= 2`) are **field tags** — they go on the wire, not the names. NEVER reuse a tag for a different field; that breaks backward compatibility.

### gRPC's four call types

1. **Unary** — one request, one response. Like HTTP/REST.
2. **Server streaming** — one request, stream of responses. Useful for "subscribe to updates."
3. **Client streaming** — stream of requests, one response. Useful for uploads.
4. **Bidirectional streaming** — both sides stream simultaneously. Chat-like.

### Why gRPC over REST

- **Smaller payloads** — binary, ~30-50% of equivalent JSON
- **Faster** — protobuf is faster to encode/decode than JSON
- **Schema-enforced** — both client and server are generated from the same .proto. No "is the field nullable?" debates.
- **Streaming** — first-class, where REST has to bolt on WebSockets/SSE
- **Cross-language** — generate Go server, TS client, Python client, all from one .proto

### Why NOT gRPC

- **Browser support is awful** — needs gRPC-Web proxy
- **Debug tooling weaker than REST** — can't `curl` it (you need `grpcurl`)
- **Operational complexity** — more moving parts (codegen, protoc, etc.)
- **Public APIs are still REST** — gRPC is mostly internal-only

### Generating code

```bash
# Install protoc + Go plugin
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go server stubs
protoc --go_out=. --go-grpc_out=. proto/task.proto
```

You now have `Task`, `CreateTaskRequest`, etc. as Go structs, plus a `TaskServiceServer` interface to implement, plus `TaskServiceClient` for callers.

### Implementing a gRPC server in Go

```go
type taskServer struct {
    pb.UnimplementedTaskServiceServer
    db *sql.DB
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
    if req.Title == "" {
        return nil, status.Error(codes.InvalidArgument, "title is required")
    }
    var id int64
    err := s.db.QueryRowContext(ctx,
        `INSERT INTO tasks(title) VALUES ($1) RETURNING id`,
        req.Title).Scan(&id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "db: %v", err)
    }
    return &pb.Task{Id: id, Title: req.Title, Completed: false}, nil
}

func main() {
    lis, _ := net.Listen("tcp", ":50051")
    grpcServer := grpc.NewServer()
    pb.RegisterTaskServiceServer(grpcServer, &taskServer{db: ...})
    grpcServer.Serve(lis)
}
```

gRPC error codes (`codes.NotFound`, `codes.InvalidArgument`, `codes.PermissionDenied`, ...) map cleanly to HTTP status codes.

---

<a name="module-59--when-what"></a>
## Module 5.9 — When REST, when gRPC, when GraphQL

| Use case | Pick |
|---|---|
| Public API for web/mobile apps | **REST** (or GraphQL) |
| Internal microservices, performance matters | **gRPC** |
| Many distinct clients (mobile, web, third-party) | **REST** |
| Single client team, full control of stack | **gRPC** is fine |
| Browser-first realtime | **WebSocket** or **Server-Sent Events** |
| Mobile app needs few queries with many fields | **GraphQL** can shine |
| API Gateway exposing microservices | **REST or gRPC-Gateway** (translates REST → gRPC) |
| Event-driven, fire-and-forget | **Kafka / NATS / SQS**, not RPC at all |

### GraphQL — brief mention

> 📖 **Definition — GraphQL:** A query language for APIs. Clients specify exactly which fields they need. Server returns exactly that. One endpoint, one POST per query. Schema-first.

Pros: fewer round trips for nested data; clients fetch exactly what they need; great DX.

Cons: complex caching (every query is unique); easy to write expensive queries (N+1, deep nesting); auth is per-resolver (more error-prone).

GraphQL has a place — content-heavy apps with diverse clients (Shopify, GitHub, Apollo customers). For most internal APIs, REST or gRPC wins on simplicity.

---

<a name="module-510--observability"></a>
## Module 5.10 — Observability for APIs

You're going to Datadog. Observability is in the air. Three pillars:

> 📖 **Definition — Observability — three pillars:**
> - **Metrics:** numeric time-series. "request rate," "p99 latency."
> - **Logs:** discrete events. "user 42 created a task."
> - **Traces:** the path of one request through many services. "POST /tasks → auth → db.insert → publish to kafka."

### Structured logging

Plain `console.log("user", userId, "did", action)` doesn't scale. Use structured logs (JSON) so log aggregators can parse them.

```ts
import pino from "pino"
const log = pino()

log.info({ userId: 42, action: "create_task", taskId: 99 }, "task created")
// {"level":30,"time":1234,"userId":42,"action":"create_task","taskId":99,"msg":"task created"}
```

Now you can query `level >= warn AND userId = 42`.

### Request IDs — propagation

Every request gets a unique ID. Pass it through ALL downstream calls (HTTP `X-Request-ID` header, gRPC metadata). Include it in every log line. When a customer reports a bug at 3:42pm, you grep the request ID and see the entire journey.

### The four golden signals

Google SRE book — what to alert on:
1. **Latency** — how slow are requests
2. **Traffic** — how many requests
3. **Errors** — error rate
4. **Saturation** — how full the system is (CPU, memory, queue depth)

For HTTP: track p50, p90, p99 latency separately. Mean is misleading — a few slow outliers can stay invisible.

### OpenTelemetry — the standard

> 📖 **Definition — OpenTelemetry (OTel):** Open standard for emitting metrics, logs, and traces. Vendor-neutral — works with Datadog, Honeycomb, Grafana Tempo, Jaeger, etc.

```ts
import { trace } from "@opentelemetry/api"

const tracer = trace.getTracer("taskapi")

app.post("/tasks", async (req, res) => {
    const span = tracer.startSpan("create_task")
    try {
        const task = await db.tasks.insert(req.body)
        span.setAttributes({ "task.id": task.id })
        res.status(201).json(task)
    } catch (e) {
        span.recordException(e)
        span.setStatus({ code: 2 /* ERROR */ })
        throw e
    } finally {
        span.end()
    }
})
```

A span is a unit of work. It has a start time, end time, attributes, and a parent (forming a tree per request). Datadog APM ingests OTel spans natively. **Master OTel before your internship — it's the lingua franca of Datadog product.**

---

<a name="project-taskapi"></a>
## 🛠️ Project: TaskAPI — a real backend

A backend you could ship. TypeScript + Express + Postgres + Redis + JWT + rate limiting + structured logging + OpenTelemetry traces.

**See `projects/taskapi/` — full code.**

### Spec

Endpoints:
```
POST   /auth/login              → { access_token, refresh_token }
POST   /auth/refresh            → new access_token
POST   /tasks                   (auth) → 201 Created
GET    /tasks                   (auth) cursor pagination
GET    /tasks/:id               (auth, ownership)
PATCH  /tasks/:id               (auth, ownership)
DELETE /tasks/:id               (auth, ownership)
GET    /healthz                 (public)
```

Cross-cutting:
- Zod validation on every body/query
- JWT auth middleware
- Per-user rate limit: 100 req/min via Redis
- Idempotency-Key support on POST /tasks
- Structured logs with request IDs
- OpenTelemetry spans on each handler
- Postgres for tasks, users, idempotency records
- Redis for rate limiting + refresh token revocation list
- Docker-compose for local Postgres+Redis
- OpenAPI spec checked into the repo

### Layout

```
projects/taskapi/
├── README.md
├── package.json
├── tsconfig.json
├── docker-compose.yml
├── openapi.yaml
├── src/
│   ├── server.ts              <- Express app
│   ├── db/
│   │   ├── client.ts
│   │   └── migrations/
│   │       └── 001_init.sql
│   ├── routes/
│   │   ├── auth.ts
│   │   ├── tasks.ts
│   │   └── health.ts
│   ├── middleware/
│   │   ├── auth.ts            <- JWT verify
│   │   ├── rateLimit.ts       <- Redis token bucket
│   │   ├── idempotency.ts
│   │   ├── requestId.ts
│   │   └── errorHandler.ts
│   ├── lib/
│   │   ├── jwt.ts
│   │   ├── logger.ts          <- pino
│   │   └── tracing.ts         <- OpenTelemetry init
│   └── schemas/
│       ├── task.ts            <- Zod schemas
│       └── auth.ts
└── tests/
    └── tasks.test.ts
```

### Bonus: gRPC twin

```
projects/taskapi-grpc/
├── proto/
│   └── task.proto
├── server/
│   └── main.go                <- Go gRPC server
└── client-ts/
    └── client.ts              <- TS gRPC client
```

Same business logic, two RPC styles. You'll feel the difference.

---

<a name="exercises"></a>
## Exercises

1. **Idempotency design.** Sketch the DB schema for an idempotency-key store. What's the TTL? What index do you need? What happens on a hash collision?

2. **JWT vs sessions essay (one page).** Pick a scenario (mobile app, server-to-server, browser SPA). Argue for one. Mention revocation, scaling, refresh tokens.

3. **Rate-limit Lua.** Implement sliding-window-counter rate limiting as a Redis Lua script. Test that two concurrent clients can't both squeeze in over the limit.

4. **OpenAPI from existing API.** Take any of your past projects (StackSense, etc.), write an OpenAPI spec for at least 5 endpoints, render it with Swagger UI.

5. **gRPC streaming demo.** Build a tiny gRPC service with one server-streaming RPC: `SubscribeMetrics(filter) returns (stream MetricPoint)`. Stream synthetic points every 100ms. Connect a TS client.

6. **Resilience drill.** Take a service that calls `https://httpbin.org/status/500`. Add: timeout, exponential backoff with jitter, max 3 retries, circuit breaker. Demonstrate it doesn't melt down when the upstream is broken.

7. **Tracing.** Add OTel spans to TaskAPI and ship them to Jaeger (run with Docker). Open the trace UI; find a slow request; identify which span is slow.

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

1. **Design an API for a TODO app.** *(Standard warm-up. Walk through resources, methods, status codes, pagination, auth.)*

2. **REST vs gRPC. When each? Trade-offs?**

3. **What does idempotency mean? Which HTTP methods are idempotent?**

4. **How does JWT work? What's a JWT vulnerable to?** *(`alg: none` attack if the lib doesn't pin algo. Stolen-token replay if not short-lived. Inability to revoke unless you maintain a revocation list.)*

5. **Implement rate limiting using Redis. What algorithm, why?**

6. **What status code do you return for: a successful create, a missing resource, a permission error, an invalid request body, a server bug, a downstream timeout?**

7. **Why HTTPS for an API even if the API doesn't handle PII?** *(Header tampering, request-smuggling, ISP snooping, public Wi-Fi MITM, search engines indexing weird HTTP-only URLs.)*

8. **How do you do API versioning? Trade-offs?**

9. **Walk through what happens server-side when a request comes in and the auth middleware accepts it.** *(Token verify → load user from DB or cache → attach to request → call next → handler runs → response → log → emit trace.)*

10. **Your API just got 10x more traffic. What breaks first?** *(Database connections; rate limiters; downstream services without timeouts; logs IO; memory.)*

11. **Describe an N+1 query problem in an API context.**

12. **What's a "thundering herd" and how do you avoid it?** *(All clients retry simultaneously after an outage. Solution: jittered backoff, circuit breakers, server-side request coalescing.)*

13. **CORS — what is it, who enforces it, why does it exist?**

14. **Compare cursor pagination to offset pagination.**

15. **Describe how you'd add observability to an existing legacy API.** *(Structured logs, request IDs, p50/p90/p99 latency, error rate, traces — start with logs+ID+latency, layer on traces last.)*

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] REST principles and pragmatic API design
- [ ] HTTP status codes for every common scenario
- [ ] JSON Schema, Zod, OpenAPI
- [ ] Sessions vs JWT — when each, refresh tokens
- [ ] OAuth flow at a conceptual level
- [ ] RBAC and ABAC
- [ ] Rate limiting algorithms (fixed/sliding window, token bucket)
- [ ] Idempotency-Key pattern
- [ ] Retries: when, how, jittered backoff
- [ ] Timeouts and circuit breakers
- [ ] gRPC + protobuf basics, four call types
- [ ] When to pick REST vs gRPC vs GraphQL
- [ ] The three pillars of observability + four golden signals
- [ ] OpenTelemetry essentials
- [ ] You've built TaskAPI end-to-end

---

**Next:** [Phase 6 — Concurrency & OS](../phase-06-concurrency-os/README.md)
