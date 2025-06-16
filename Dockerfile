# 1. 빌드 스테이지
FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# 빌드 타겟을 ARG로 받아서 빌드
ARG TARGET=server-1
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/${TARGET}

# 2. 실행 스테이지
FROM scratch

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]
