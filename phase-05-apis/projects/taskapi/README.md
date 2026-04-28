# TaskAPI

> A real backend, in TypeScript + Express + Postgres + Redis. Auth, rate limiting, idempotency, structured logs, OpenTelemetry traces — the full kit.

## Run locally

```bash
# 1. Spin up Postgres + Redis
docker-compose up -d

# 2. Install + migrate
npm install
npm run migrate

# 3. Start
npm run dev
```

Then:
```bash
# Login (returns a JWT)
curl -X POST localhost:3000/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"k@x.com","password":"hunter2"}'

# Use the token
TOKEN=...
curl localhost:3000/tasks -H "Authorization: Bearer $TOKEN"

# Create a task with idempotency
curl -X POST localhost:3000/tasks \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -H "Idempotency-Key: $(uuidgen)" \
    -d '{"title":"learn distributed systems"}'
```

## What this project teaches

- Express + middleware chaining
- JWT issue/verify with refresh tokens
- Postgres via `pg`, transactions with `BEGIN/COMMIT`
- Redis token-bucket rate limiting
- Idempotency-Key store + retrieval
- Structured logging with `pino`
- OpenTelemetry instrumentation
- Zod-based request validation
- OpenAPI doc

## Project layout

```
taskapi/
├── package.json
├── tsconfig.json
├── docker-compose.yml         <- Postgres + Redis + Jaeger
├── openapi.yaml               <- API spec
├── src/
│   ├── server.ts              <- main entry
│   ├── db/
│   │   ├── client.ts          <- pg Pool
│   │   └── migrations/001_init.sql
│   ├── lib/
│   │   ├── jwt.ts
│   │   ├── logger.ts
│   │   └── tracing.ts
│   ├── middleware/
│   │   ├── auth.ts
│   │   ├── rateLimit.ts
│   │   ├── idempotency.ts
│   │   └── requestId.ts
│   ├── routes/
│   │   ├── auth.ts
│   │   ├── tasks.ts
│   │   └── health.ts
│   └── schemas/
│       └── task.ts
└── tests/
    └── tasks.test.ts
```

> 📝 This README ships with the project skeleton — full files are in this directory. Some advanced parts (full migration runner, complete OpenAPI spec) are sketched; you'll fill them in as you work through the curriculum and adapt them to your StackSense/Datadog work.
