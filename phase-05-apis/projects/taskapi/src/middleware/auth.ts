// src/middleware/auth.ts
import type { Request, Response, NextFunction } from "express"
import { verifyAccessToken } from "../lib/jwt.js"

declare global {
    namespace Express {
        interface Request {
            user?: { id: number; email: string }
        }
    }
}

export function requireAuth(req: Request, res: Response, next: NextFunction) {
    const auth = req.header("Authorization")
    if (!auth?.startsWith("Bearer ")) {
        return res.status(401).json({
            error: { code: "unauthenticated", message: "missing bearer token" },
        })
    }
    try {
        const claims = verifyAccessToken(auth.slice(7))
        req.user = { id: Number(claims.sub), email: claims.email }
        next()
    } catch (e) {
        res.status(401).json({
            error: { code: "invalid_token", message: "token invalid or expired" },
        })
    }
}
