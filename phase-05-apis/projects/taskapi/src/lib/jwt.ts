// src/lib/jwt.ts
//
// JWT issue + verify. Two token types: short access (15min) + long refresh (30d).
// In a real system the secret comes from a secret manager, not env vars.

import jwt from "jsonwebtoken"
import crypto from "node:crypto"

const SECRET = process.env.JWT_SECRET ?? "dev-only-do-not-ship-this"

export interface AccessTokenClaims {
    sub: string       // user ID as string
    email: string
}

export function issueAccessToken(userId: number, email: string): string {
    return jwt.sign(
        { sub: String(userId), email } satisfies AccessTokenClaims,
        SECRET,
        { algorithm: "HS256", expiresIn: "15m" }
    )
}

export function verifyAccessToken(token: string): AccessTokenClaims {
    // jwt.verify throws on invalid/expired — let the caller handle.
    return jwt.verify(token, SECRET, { algorithms: ["HS256"] }) as AccessTokenClaims
}

// Refresh tokens are opaque random strings (NOT JWTs).
// We store the SHA-256 hash in DB so a leaked DB doesn't leak tokens.
export function generateRefreshToken(): { token: string; hash: string } {
    const token = crypto.randomBytes(32).toString("hex")
    const hash = crypto.createHash("sha256").update(token).digest("hex")
    return { token, hash }
}

export function hashRefreshToken(token: string): string {
    return crypto.createHash("sha256").update(token).digest("hex")
}
