// src/middleware/idempotency.ts
//
// Honor Idempotency-Key header on POSTs. Stores (key -> response) in Postgres.
// On retry with the same key + same request body hash, return the stored response.

import type { Request, Response, NextFunction } from "express"
import crypto from "node:crypto"
import { pool } from "../db/client.js"

export async function idempotency(req: Request, res: Response, next: NextFunction) {
    const key = req.header("Idempotency-Key")
    if (!key || !req.user) return next()

    const hash = crypto.createHash("sha256")
        .update(JSON.stringify(req.body ?? null))
        .digest("hex")

    const { rows } = await pool.query(
        `SELECT request_hash, response_status, response_body
         FROM idempotency_keys
         WHERE key = $1 AND user_id = $2`,
        [key, req.user.id]
    )

    if (rows.length > 0) {
        const stored = rows[0]
        if (stored.request_hash !== hash) {
            return res.status(409).json({
                error: { code: "idempotency_collision",
                         message: "key already used with different request body" },
            })
        }
        return res.status(stored.response_status).json(stored.response_body)
    }

    // Wrap res.json to capture and persist the response.
    const origJson = res.json.bind(res)
    res.json = (body: any) => {
        // Fire-and-forget the insert; don't block the response.
        pool.query(
            `INSERT INTO idempotency_keys
                (key, user_id, request_hash, response_status, response_body)
             VALUES ($1, $2, $3, $4, $5)
             ON CONFLICT (key) DO NOTHING`,
            [key, req.user!.id, hash, res.statusCode, body]
        ).catch(() => {/* swallow; logged at the error handler level */})
        return origJson(body)
    }
    next()
}
