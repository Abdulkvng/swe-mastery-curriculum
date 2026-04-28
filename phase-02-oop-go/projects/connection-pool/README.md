# connection-pool

> A production-style database connection pool, written from scratch in Go.

## Why this project

Connection pools are everywhere: `database/sql.DB` already has one; pgbouncer is one; HikariCP (Java) is one. Building one yourself teaches you:

- Concurrency primitives (mutexes, channels, `sync.Once`)
- Resource management (acquire/release, timeouts, context cancellation)
- Lifecycle management (creation, health checks, eviction, shutdown)
- API design (what does the public surface look like?)
- Many design patterns at once (Singleton, Strategy, Observer)

It's also a top-tier interview question. By the end, you'll be able to whiteboard one cold.

## Layout

```
connection-pool/
├── README.md
├── go.mod
├── pool/                    <- the library
│   ├── pool.go              <- main Pool type
│   ├── conn.go              <- Conn wrapper
│   ├── config.go            <- Config + defaults
│   ├── stats.go             <- metrics
│   ├── singleton.go         <- optional Default() pool with sync.Once
│   ├── pool_test.go         <- unit tests
│   └── concurrency_test.go  <- stress tests
└── examples/
    └── basic/
        └── main.go          <- real Postgres usage
```

## Setup

```bash
cd connection-pool
go mod tidy
go test ./...
```

To run the example you need a real Postgres. Easiest:

```bash
docker run --rm -d --name pg \
    -e POSTGRES_PASSWORD=pw \
    -p 5432:5432 \
    postgres:16

# Run example
DATABASE_URL="postgres://postgres:pw@localhost:5432/postgres?sslmode=disable" \
    go run ./examples/basic
```

## Concepts mapped to code

| OOP concept | Where it lives |
|---|---|
| Encapsulation | `Pool` struct's lowercase fields (`idle`, `inUse`, `mu`) |
| Abstraction | `Pool` interface vs concrete `*pool` impl |
| Singleton | `pool/singleton.go` using `sync.Once` |
| Strategy | `Config.HealthCheck` is a pluggable function |
| Observer | `Config.OnAcquire`, `OnRelease`, `OnHealthFail` callbacks |
| SOLID — DIP | We depend on `Connector` interface, not `*sql.DB` directly |
| Constructor | `NewPool(cfg Config) (Pool, error)` |
