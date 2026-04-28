// examples/basic/main.go
//
// Example using the connection pool with a real Postgres database.
//
// To run:
//   docker run --rm -d --name pg -e POSTGRES_PASSWORD=pw -p 5432:5432 postgres:16
//   DATABASE_URL='postgres://postgres:pw@localhost:5432/postgres?sslmode=disable' \
//       go run ./examples/basic

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq" // register Postgres driver

	"github.com/kvng/connection-pool/pool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("set DATABASE_URL")
	}

	cfg := pool.DefaultConfig()
	cfg.Connector = func(ctx context.Context) (*sql.DB, error) {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}
		// Each *sql.DB is itself a pool internally; we set it to 1 so OUR
		// pool is in charge.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		return db, db.PingContext(ctx)
	}
	cfg.MinConns = 2
	cfg.MaxConns = 5
	cfg.AcquireTimeout = 3 * time.Second
	cfg.MaxLifetime = 10 * time.Minute
	cfg.MaxIdleTime = 1 * time.Minute
	cfg.HealthInterval = 30 * time.Second

	cfg.Observers = pool.Observers{
		OnAcquire: func(c *pool.Conn) { fmt.Printf("[acquire] conn#%d\n", c.ID()) },
		OnRelease: func(c *pool.Conn) { fmt.Printf("[release] conn#%d\n", c.ID()) },
		OnCreate:  func(c *pool.Conn) { fmt.Printf("[create]  conn#%d\n", c.ID()) },
		OnDestroy: func(c *pool.Conn, reason string) {
			fmt.Printf("[destroy] conn#%d reason=%s\n", c.ID(), reason)
		},
	}

	p, err := pool.NewPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Hammer it with a bunch of concurrent queries.
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			c, err := p.Acquire(ctx)
			if err != nil {
				fmt.Printf("worker %d: acquire failed: %v\n", id, err)
				return
			}
			defer c.Release()

			var version string
			if err := c.DB().QueryRow("SELECT version()").Scan(&version); err != nil {
				fmt.Printf("worker %d: query failed: %v\n", id, err)
				return
			}
			fmt.Printf("worker %d: %s\n", id, version[:30])
		}(i)
	}
	wg.Wait()

	fmt.Printf("\nFinal stats: %+v\n", p.Stats())
}
