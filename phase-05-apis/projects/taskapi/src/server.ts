// src/server.ts
import express from "express"
import pinoHttp from "pino-http"
import { v4 as uuidv4 } from "uuid"
import { logger } from "./lib/logger.js"
import { rateLimit } from "./middleware/rateLimit.js"
import authRoutes from "./routes/auth.js"
import tasksRoutes from "./routes/tasks.js"

const app = express()
app.use(express.json({ limit: "1mb" }))

// Request ID middleware: generate or propagate one.
app.use((req, res, next) => {
    const id = req.header("X-Request-ID") ?? uuidv4()
    res.setHeader("X-Request-ID", id)
    ;(req as any).requestId = id
    next()
})

// Structured access log.
app.use(pinoHttp({
    logger,
    customProps: req => ({ requestId: (req as any).requestId }),
}))

app.get("/healthz", (_req, res) => res.json({ ok: true }))

app.use("/auth", authRoutes)
app.use("/tasks", rateLimit({ max: 100, windowSec: 60 }), tasksRoutes)

// Error handler — last middleware.
app.use((err: any, req: any, res: any, _next: any) => {
    req.log?.error({ err }, "unhandled")
    res.status(500).json({
        error: { code: "internal", message: "internal server error" },
    })
})

const port = Number(process.env.PORT ?? 3000)
app.listen(port, () => logger.info({ port }, "taskapi listening"))
