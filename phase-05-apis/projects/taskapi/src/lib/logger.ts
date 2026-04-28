// src/lib/logger.ts
import pino from "pino"

export const logger = pino({
    level: process.env.LOG_LEVEL ?? "info",
    // In dev, pretty-print. In prod, raw JSON for log aggregators.
    ...(process.env.NODE_ENV === "production"
        ? {}
        : {
            transport: {
                target: "pino-pretty",
                options: { colorize: true },
            },
        }),
})
