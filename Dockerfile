FROM golang:1.22-bookworm AS builder

WORKDIR /work
COPY . .

RUN go mod download \
	&& CGO_ENABLED=0 go build -o /usr/local/bin/casket-example ./cmd/example

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
	ca-certificates \
	curl \
	procps \
	xz-utils \
	&& rm -rf /var/lib/apt/lists/*

# Скачиваем и распаковываем свежий 7zz (25.01)
RUN curl -L https://7-zip.org/a/7z2501-linux-arm64.tar.xz -o /tmp/7z.tar.xz \
	&& tar -C /usr/local/bin -xf /tmp/7z.tar.xz \
	&& rm /tmp/7z.tar.xz \
	&& chmod +x /usr/local/bin/7zz

CMD ["tail", "-f", "/dev/null"]
