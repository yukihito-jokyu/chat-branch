FROM golang:1.25.4-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/server/main.go

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

RUN groupadd -r appuser && useradd -r -g appuser appuser
USER appuser

COPY --from=builder /app/main /usr/local/bin/server
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/config ./config

EXPOSE 1323

CMD ["server"]
