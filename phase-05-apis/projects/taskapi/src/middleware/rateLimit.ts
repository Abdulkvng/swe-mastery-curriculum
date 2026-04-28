// src/middleware/rateLimit.ts
//
// Per-user rate limiting via Redis. Simple fixed-window. For production,
// upgrade to sliding window or token bucket (covered in Phase 5 module 5.7).

import type { Request, Response, NextFunction } from "express"
import Redis from "ioredis"

const redis = new Redis(process.env.REDIS_URL ?? "redis://localhost:6379")

export function rateLimit(opts: { max: number; windowSec: number }) {
    return async (req: Request, res: Response, next: NextFunction) => {
        // Use authenticated user ID if available, else IP. (Both are imperfect.)
        const id = req.user?.id ?? req.ip ?? "anon"
        const window = Math.floor(Date.now() / 1000 / opts.windowSec)
        const key = `rl:${id}:${window}`

        const count = await redis.incr(key)
        if (count === 1) {
            await redis.expire(key, opts.windowSec)
        }
        res.setHeader("X-RateLimit-Limit", String(opts.max))
        res.setHeader("X-RateLimit-Remaining", String(Math.max(0, opts.max - count)))

        if (count > opts.max) {
            res.setHeader("Retry-After", String(opts.windowSec))
            return res.status(429).json({
                error: { code: "rate_limited", message: "too many requests" },
            })
        }
        next()
    }
}
