FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./src ./src
COPY cmd ./cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go


FROM debian:bookworm

ENV CONTAINER_MODE=1

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/* \
    && curl -sSf https://atlasgo.sh | sh

WORKDIR /app
COPY --from=builder /out/app ./app
COPY migration/versions ./migration/versions

ENTRYPOINT ["./app"]
