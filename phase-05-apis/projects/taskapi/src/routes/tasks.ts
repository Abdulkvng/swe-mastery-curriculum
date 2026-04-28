// src/routes/tasks.ts
import { Router } from "express"
import { z } from "zod"
import { pool } from "../db/client.js"
import { requireAuth } from "../middleware/auth.js"
import { idempotency } from "../middleware/idempotency.js"

const router = Router()
router.use(requireAuth)

// === Schemas ===
const CreateTaskSchema = z.object({
    title: z.string().min(1).max(200),
    body: z.string().max(10_000).optional(),
})

const UpdateTaskSchema = z.object({
    title: z.string().min(1).max(200).optional(),
    body: z.string().max(10_000).optional(),
    completed: z.boolean().optional(),
})

const ListQuerySchema = z.object({
    limit: z.coerce.number().int().min(1).max(100).default(20),
    cursor: z.string().optional(),
})

// === Routes ===

router.post("/", idempotency, async (req, res, next) => {
    try {
        const parsed = CreateTaskSchema.safeParse(req.body)
        if (!parsed.success) {
            return res.status(400).json({
                error: { code: "validation_failed", details: parsed.error.format() },
            })
        }
        const { title, body } = parsed.data
        const { rows } = await pool.query(
            `INSERT INTO tasks(user_id, title, body)
             VALUES ($1, $2, $3)
             RETURNING id, title, body, completed, created_at, updated_at`,
            [req.user!.id, title, body ?? null]
        )
        res.status(201).json(rows[0])
    } catch (e) { next(e) }
})

router.get("/", async (req, res, next) => {
    try {
        const parsed = ListQuerySchema.safeParse(req.query)
        if (!parsed.success) {
            return res.status(400).json({
                error: { code: "validation_failed", details: parsed.error.format() },
            })
        }
        const { limit, cursor } = parsed.data

        // Cursor is the created_at + id of the last seen item, base64-encoded.
        let cursorClause = ""
        const args: any[] = [req.user!.id, limit]
        if (cursor) {
            try {
                const decoded = JSON.parse(Buffer.from(cursor, "base64").toString())
                cursorClause = `AND (created_at, id) < ($3, $4)`
                args.push(decoded.t, decoded.id)
            } catch {/* invalid cursor, ignore */}
        }

        const { rows } = await pool.query(
            `SELECT id, title, body, completed, created_at, updated_at
             FROM tasks
             WHERE user_id = $1 ${cursorClause}
             ORDER BY created_at DESC, id DESC
             LIMIT $2`,
            args
        )

        const nextCursor = rows.length === limit
            ? Buffer.from(JSON.stringify({
                  t: rows[rows.length - 1].created_at,
                  id: rows[rows.length - 1].id,
              })).toString("base64")
            : null

        res.json({ items: rows, next_cursor: nextCursor })
    } catch (e) { next(e) }
})

router.get("/:id", async (req, res, next) => {
    try {
        const id = Number(req.params.id)
        if (!Number.isFinite(id)) return res.status(400).json({ error: { code: "bad_id" } })
        const { rows } = await pool.query(
            `SELECT id, title, body, completed, created_at, updated_at
             FROM tasks WHERE id = $1 AND user_id = $2`,
            [id, req.user!.id]
        )
        // Return 404 not 403 to avoid leaking the task's existence.
        if (rows.length === 0) {
            return res.status(404).json({ error: { code: "not_found" } })
        }
        res.json(rows[0])
    } catch (e) { next(e) }
})

router.patch("/:id", async (req, res, next) => {
    try {
        const id = Number(req.params.id)
        const parsed = UpdateTaskSchema.safeParse(req.body)
        if (!parsed.success) {
            return res.status(400).json({
                error: { code: "validation_failed", details: parsed.error.format() },
            })
        }
        const fields = parsed.data
        const sets: string[] = []
        const values: any[] = []
        let n = 1
        for (const [k, v] of Object.entries(fields)) {
            if (v === undefined) continue
            sets.push(`${k} = $${n++}`)
            values.push(v)
        }
        if (sets.length === 0) {
            return res.status(400).json({ error: { code: "no_fields" } })
        }
        sets.push(`updated_at = now()`)
        values.push(id, req.user!.id)
        const { rows } = await pool.query(
            `UPDATE tasks SET ${sets.join(", ")}
             WHERE id = $${n++} AND user_id = $${n}
             RETURNING id, title, body, completed, created_at, updated_at`,
            values
        )
        if (rows.length === 0) {
            return res.status(404).json({ error: { code: "not_found" } })
        }
        res.json(rows[0])
    } catch (e) { next(e) }
})

router.delete("/:id", async (req, res, next) => {
    try {
        const id = Number(req.params.id)
        const { rowCount } = await pool.query(
            `DELETE FROM tasks WHERE id = $1 AND user_id = $2`,
            [id, req.user!.id]
        )
        if (rowCount === 0) {
            return res.status(404).json({ error: { code: "not_found" } })
        }
        res.status(204).end()
    } catch (e) { next(e) }
})

export default router
