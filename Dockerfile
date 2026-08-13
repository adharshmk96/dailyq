# syntax=docker/dockerfile:1

# 1. Build the static Nuxt SPA.
FROM oven/bun:1 AS ui
WORKDIR /src/dailyq-ui
COPY dailyq-ui/package.json dailyq-ui/bun.lock* ./
RUN bun install --frozen-lockfile
COPY dailyq-ui/ ./
RUN bun run generate

# 2. Compile the Go server with the SPA embedded.
FROM golang:1.26-alpine AS api
WORKDIR /src/dailyq-api
COPY dailyq-api/go.mod dailyq-api/go.sum ./
RUN go mod download
COPY dailyq-api/ ./
# internal/web/dist is gitignored; replace it with the generated SPA.
COPY --from=ui /src/dailyq-ui/.output/public/ ./internal/web/dist/
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/dailyq .

# 3. Runtime.
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata wget \
    && adduser -D -u 10001 dailyq \
    && mkdir -p /data && chown dailyq:dailyq /data
COPY --from=api /out/dailyq /usr/local/bin/dailyq
USER dailyq
WORKDIR /data
VOLUME ["/data"]

ENV DAILYQ_SERVER_HOST=0.0.0.0 \
    DAILYQ_SERVER_PORT=8080 \
    DAILYQ_SERVER_MODE=release \
    DAILYQ_DATABASE_DSN=/data/dailyq.db \
    DAILYQ_LOG_FORMAT=json

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["dailyq"]
CMD ["serve"]
