# 1. 빌드 스테이지
FROM golang:1.22 as builder

WORKDIR /app

# go.mod만 복사
COPY go.mod ./
RUN go mod download

# 전체 프로젝트 복사
COPY . .

# server-1 기준으로 빌드
RUN go build -o app ./cmd/server-1

# 2. 실행 스테이지
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/app .

# config, data 디렉토리 마운트
VOLUME ["/app/config.json", "/app/data"]

# 앱 실행
CMD ["./app"]
