# Simple dev Dockerfile (Alpine-based, smaller than Debian Bookworm)
FROM golang:1.22-alpine

# Install tools: 7zip (7zz) and p7zip (7z)
RUN apk add --no-cache ca-certificates 7zip p7zip procps \
 && update-ca-certificates \
 && command -v 7zz >/dev/null \
 && command -v 7z >/dev/null

WORKDIR /work
COPY . .

# Build example CLI
RUN go mod download \
 && CGO_ENABLED=0 go build -o /usr/local/bin/casket-example ./cmd/example

# Keep container running for interactive work (BusyBox sleep has no "infinity")
CMD ["tail", "-f", "/dev/null"]
