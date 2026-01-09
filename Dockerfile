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

# デフォルト環境変数（Coolifyで上書き可能）
ENV DB_HOST=yckco84g0sowg8scggoos44w \
    DB_PORT=5432 \
    DB_USER=tsunagu \
    DB_PASSWORD=N14xBDZdgAgRUYjN5ek63n1OFs8lTWE7I5poLRzF0SCSMhlz32PQx7L3ARZfGcQE \
    DB_NAME=tsunagu_db \
    DB_SSLMODE=disable \
    SERVER_PORT=8080

EXPOSE 8080

CMD ["./server"]
