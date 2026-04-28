// src/routes/auth.ts
//
// Login + refresh-token endpoints. Real password hashing is bcrypt or argon2;
// here we use a stand-in for brevity. Replace with `bcrypt` before any real use.

import { Router } from "express"
import { z } from "zod"
import crypto from "node:crypto"
import { pool } from "../db/client.js"
import {
    issueAccessToken, generateRefreshToken, hashRefreshToken,
} from "../lib/jwt.js"

const router = Router()

const LoginSchema = z.object({
    email: z.string().email(),
    password: z.string().min(8),
})

// Tiny placeholder — DO NOT use in production. Use bcrypt.
function checkPassword(plain: string, hash: string): boolean {
    return crypto.createHash("sha256").update(plain).digest("hex") === hash
}

router.post("/login", async (req, res, next) => {
    try {
        const parsed = LoginSchema.safeParse(req.body)
        if (!parsed.success) {
            return res.status(400).json({
                error: { code: "validation_failed", details: parsed.error.format() },
            })
        }
        const { email, password } = parsed.data
        const { rows } = await pool.query(
            `SELECT id, email, password_hash FROM users WHERE email = $1`, [email]
        )
        if (rows.length === 0 || !checkPassword(password, rows[0].password_hash)) {
            // Same response for "no user" and "wrong password" — don't leak
            // whether an email is registered.
            return res.status(401).json({ error: { code: "invalid_credentials" } })
        }
        const user = rows[0]
        const access = issueAccessToken(user.id, user.email)
        const { token: refresh, hash } = generateRefreshToken()
        const expiresAt = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)
        await pool.query(
            `INSERT INTO refresh_tokens(token_hash, user_id, expires_at)
             VALUES ($1, $2, $3)`,
            [hash, user.id, expiresAt]
        )
        res.json({ access_token: access, refresh_token: refresh })
    } catch (e) { next(e) }
})

router.post("/refresh", async (req, res, next) => {
    try {
        const { refresh_token } = req.body ?? {}
        if (!refresh_token || typeof refresh_token !== "string") {
            return res.status(400).json({ error: { code: "missing_refresh_token" } })
        }
        const hash = hashRefreshToken(refresh_token)
        const { rows } = await pool.query(
            `SELECT user_id, expires_at, revoked_at
             FROM refresh_tokens WHERE token_hash = $1`,
            [hash]
        )
        if (rows.length === 0
            || rows[0].revoked_at !== null
            || new Date(rows[0].expires_at) < new Date()) {
            return res.status(401).json({ error: { code: "invalid_refresh_token" } })
        }
        const userId = rows[0].user_id
        const u = await pool.query(`SELECT email FROM users WHERE id = $1`, [userId])
        const access = issueAccessToken(userId, u.rows[0].email)
        res.json({ access_token: access })
    } catch (e) { next(e) }
})

router.post("/logout", async (req, res, next) => {
    try {
        const { refresh_token } = req.body ?? {}
        if (!refresh_token) return res.status(204).end()
        const hash = hashRefreshToken(refresh_token)
        await pool.query(
            `UPDATE refresh_tokens SET revoked_at = now()
             WHERE token_hash = $1 AND revoked_at IS NULL`,
            [hash]
        )
        res.status(204).end()
    } catch (e) { next(e) }
})

export default router
