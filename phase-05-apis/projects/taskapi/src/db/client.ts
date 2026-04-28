// src/db/client.ts
//
// Postgres connection pool. The `pg` library has its own pool — we configure it.
// In a real Datadog-scale system you'd front this with pgbouncer (which is what
// you built in Phase 4!).

import pg from "pg"

export const pool = new pg.Pool({
    connectionString: process.env.DATABASE_URL ??
        "postgres://taskapi:dev@localhost:5432/taskapi",
    max: 10,                          // max pool size
    idleTimeoutMillis: 30_000,        // close idle conns after 30s
    connectionTimeoutMillis: 5_000,   // fail fast on connect
})

// Helper for transactions. Pass a callback that receives a client; the helper
// handles BEGIN, COMMIT/ROLLBACK, and conn release.
export async function withTx<T>(
    fn: (client: pg.PoolClient) => Promise<T>
): Promise<T> {
    const client = await pool.connect()
    try {
        await client.query("BEGIN")
        const result = await fn(client)
        await client.query("COMMIT")
        return result
    } catch (err) {
        await client.query("ROLLBACK")
        throw err
    } finally {
        client.release()
    }
}
