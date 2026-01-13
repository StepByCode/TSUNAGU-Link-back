#!/bin/sh
set -e

echo "Starting entrypoint script..."

# 環境変数の確認
echo "Database configuration:"
echo "  DATABASE_HOST: ${DATABASE_HOST}"
echo "  DB_HOST: ${DB_HOST}"
echo "  DB_PORT: ${DB_PORT}"
echo "  DB_NAME: ${DB_NAME}"
echo "  DB_USER: ${DB_USER}"

# 必須環境変数のバリデーション
MISSING_VARS=""
[ -z "${DB_HOST}${DATABASE_HOST}" ] && MISSING_VARS="${MISSING_VARS} DB_HOST/DATABASE_HOST"
[ -z "${DB_PORT}" ] && MISSING_VARS="${MISSING_VARS} DB_PORT"
[ -z "${DB_USER}" ] && MISSING_VARS="${MISSING_VARS} DB_USER"
[ -z "${DB_PASSWORD}" ] && MISSING_VARS="${MISSING_VARS} DB_PASSWORD"
[ -z "${DB_NAME}" ] && MISSING_VARS="${MISSING_VARS} DB_NAME"
[ -z "${DB_SSLMODE}" ] && MISSING_VARS="${MISSING_VARS} DB_SSLMODE"

if [ -n "${MISSING_VARS}" ]; then
    echo "ERROR: Missing required environment variables:${MISSING_VARS}"
    echo "Please check Coolify environment variable configuration."
    exit 1
fi

# DATABASE_HOSTとDB_HOSTのフォールバック処理
ACTUAL_DB_HOST="${DATABASE_HOST:-${DB_HOST}}"
echo "  Using DB host: ${ACTUAL_DB_HOST}"

# データベース接続文字列を構築
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${ACTUAL_DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Waiting for database to be ready..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if pg_isready -h "${ACTUAL_DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    attempt=$((attempt + 1))
    echo "Database not ready yet (attempt $attempt/$max_attempts), waiting..."
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "ERROR: Database did not become ready in time"
    exit 1
fi

echo "Running database migrations..."
if [ -d "/root/migrations" ]; then
    migrate -path /root/migrations -database "${DATABASE_URL}" up
    if [ $? -eq 0 ]; then
        echo "Migrations completed successfully"
    else
        echo "ERROR: Migrations failed"
        exit 1
    fi
else
    echo "WARNING: Migrations directory not found, skipping migrations"
fi

echo "Starting application server..."
exec ./server
