FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 本番用ビルド（setupではなくbuildのみ）
RUN mkdir -p bin && go build -o bin/server ./cmd/server

FROM alpine:latest

# ビルド時引数を受け取る（Coolifyから渡される）
ARG DATABASE_HOST
ARG DB_HOST
ARG DB_PORT=5432
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG DB_SSLMODE=disable
ARG SERVER_PORT=8080
ARG JWT_SECRET
ARG JWT_EXPIRY_HOURS=24

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

# ARGを実行時環境変数（ENV）に変換
ENV DATABASE_HOST=${DATABASE_HOST} \
    DB_HOST=${DB_HOST} \
    DB_PORT=${DB_PORT} \
    DB_USER=${DB_USER} \
    DB_PASSWORD=${DB_PASSWORD} \
    DB_NAME=${DB_NAME} \
    DB_SSLMODE=${DB_SSLMODE} \
    SERVER_PORT=${SERVER_PORT} \
    JWT_SECRET=${JWT_SECRET} \
    JWT_EXPIRY_HOURS=${JWT_EXPIRY_HOURS}

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
