# syntax=docker/dockerfile:1

FROM node:20-alpine AS web-build
WORKDIR /src/web
RUN corepack enable && corepack prepare pnpm@8.15.9 --activate
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:1.25-bookworm AS server-build
WORKDIR /src
RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*
COPY server/go.mod server/go.sum ./server/
WORKDIR /src/server
RUN go mod download
WORKDIR /src
COPY server/ ./server/
COPY --from=web-build /src/web/dist ./server/internal/frontend/dist
WORKDIR /src/server
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/eibothub ./cmd/api

FROM debian:bookworm-slim AS runtime
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*
COPY --from=server-build /out/eibothub /app/eibothub
COPY deploy/ /app/deploy/
RUN chmod +x /app/eibothub /app/deploy/start-linux.sh
EXPOSE 8080
VOLUME ["/app/data", "/app/storage"]
ENTRYPOINT ["/app/eibothub"]
