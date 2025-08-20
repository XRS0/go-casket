# -- Example Dockerfile for go-casket --
FROM golang:1.22-alpine AS builder

# --- Build stage ---
WORKDIR /src
COPY . .
RUN go build -o /casket-example ./cmd/example

# --- Runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates p7zip unrar && update-ca-certificates

WORKDIR /work
COPY --from=builder /casket-example /usr/local/bin/casket-example

ENTRYPOINT ["casket-example"]
# Default args for quick test (override when running)
CMD ["/work/sample.7z", "/work/out"]
