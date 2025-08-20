# Multi-stage build: build Go binary, then run it in a minimal image with 7z/unrar

# --- Build stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /src
# Dependencies if needed (git for private/replace modules)
RUN apk add --no-cache git
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/casket-example ./cmd/example

# --- Runtime stage ---
FROM alpine:3.19
RUN apk add --no-cache p7zip unrar ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /out/casket-example /usr/local/bin/casket-example
# Default entrypoint runs the example CLI
ENTRYPOINT ["casket-example"]
# Example:
# docker build -t go-casket .
# docker run --rm -v "$PWD:/work" -w /work go-casket /work/test.rar /work/out
