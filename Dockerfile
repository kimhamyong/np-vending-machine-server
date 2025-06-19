FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY . .

ARG TARGET=server-1
RUN go build -o app ./cmd/${TARGET}

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates tzdata && apt-get clean

COPY --from=builder /app/app .
COPY internal/storage/schema.sql /app/schema.sql

CMD ["./app"]
