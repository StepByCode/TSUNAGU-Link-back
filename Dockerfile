FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 本番用ビルド（setupではなくbuildのみ）
RUN mkdir -p bin && go build -o bin/server ./cmd/server

FROM alpine:latest

# 必要なパッケージとgolang-migrateのインストール
RUN apk --no-cache add ca-certificates postgresql-client curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

WORKDIR /root/

COPY --from=builder /app/bin/server .
COPY --from=builder /app/db/migrations ./migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# 環境変数はランタイム時にCoolifyから渡される
# ビルド時のARGをENVに変換しない（空の値が焼き込まれるのを防ぐ）

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
