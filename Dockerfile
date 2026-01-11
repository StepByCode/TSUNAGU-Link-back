FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 本番用ビルド（setupではなくbuildのみ）
RUN mkdir -p bin && go build -o bin/server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/server .

# デフォルト環境変数（機密情報は実行時に渡す）
# DB_HOST, DB_USER, DB_PASSWORD, DB_NAMEは環境変数として設定してください
ENV DB_PORT=5432 \
    DB_SSLMODE=disable \
    SERVER_PORT=8080

EXPOSE 8080

CMD ["./server"]
