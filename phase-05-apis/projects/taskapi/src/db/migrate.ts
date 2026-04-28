// src/db/migrate.ts
//
// Tiny migration runner. Reads .sql files from migrations/ in order and applies
// any that haven't been applied yet. State is tracked in `_migrations` table.

import { readdir, readFile } from "node:fs/promises"
import { fileURLToPath } from "node:url"
import { dirname, join } from "node:path"
import { pool } from "./client.js"

const __dirname = dirname(fileURLToPath(import.meta.url))

async function run() {
    await pool.query(`
        CREATE TABLE IF NOT EXISTS _migrations (
            name TEXT PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
        )
    `)

    const dir = join(__dirname, "migrations")
    const files = (await readdir(dir)).filter(f => f.endsWith(".sql")).sort()

    for (const file of files) {
        const { rowCount } = await pool.query(
            `SELECT 1 FROM _migrations WHERE name = $1`, [file]
        )
        if (rowCount && rowCount > 0) {
            console.log(`[skip] ${file}`)
            continue
        }
        const sql = await readFile(join(dir, file), "utf8")
        console.log(`[apply] ${file}`)
        await pool.query(sql)
        await pool.query(
            `INSERT INTO _migrations(name) VALUES ($1)`, [file]
        )
    }

    await pool.end()
    console.log("migrations done")
}

run().catch(e => { console.error(e); process.exit(1) })
