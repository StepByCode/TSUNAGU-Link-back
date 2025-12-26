# Technology Stack

## Architecture

**レイヤードアーキテクチャ**: Handler → Service → Repository の3層構造により、責任を明確に分離

```
Handler層    : HTTPリクエスト/レスポンスの処理、入力検証
Service層    : ビジネスロジック、トランザクション制御
Repository層 : データアクセス、SQL操作
```

## Core Technologies

- **言語**: Go 1.24.7 (1.21以上)
- **Webフレームワーク**: Echo v4 - 高速で拡張性の高いGo製フレームワーク
- **データベース**: PostgreSQL 16 - エンタープライズグレードのRDBMS
- **認証**: JWT (golang-jwt/jwt/v5) - 将来的にZitadel統合予定
- **コンテナ**: Docker & Docker Compose

## Key Libraries

- **Echo v4** (`github.com/labstack/echo/v4`): HTTPルーティング、ミドルウェア
- **lib/pq** (`github.com/lib/pq`): PostgreSQLドライバ
- **golang-jwt/jwt** (`github.com/golang-jwt/jwt/v5`): JWT認証
- **golang-migrate** (ツール): データベースマイグレーション管理
- **oapi-codegen** (ツール): OpenAPIからのコード生成
- **uuid** (`github.com/google/uuid`): UUID生成
- **crypto** (`golang.org/x/crypto`): パスワードハッシュ化（bcrypt）

## Development Standards

### Type Safety
- Go言語の静的型付けを活用
- エラーハンドリングは明示的に行う（panic/recoverは最小限）
- インターフェース駆動設計（依存性注入）

### Code Quality
- `go fmt` でコードフォーマット統一
- エラーは常に処理（`if err != nil` パターン）
- ログは標準ライブラリの `log` パッケージを使用

### Testing
- `go test` による単体テスト
- レースコンディション検出 (`-race` フラグ)
- カバレッジレポート生成対応

## Development Environment

### Required Tools
- Go 1.21以上
- Docker & Docker Compose
- Make
- golang-migrate (マイグレーション)
- oapi-codegen (OpenAPIコード生成)
- Air (ホットリロード、開発時)

### Common Commands
```bash
# Dev (ホットリロード): make dev
# Build: make build
# Run: make run
# Test: make test
# Setup (初回): make setup
# Migrations: make migrate-up / make migrate-down
# OpenAPI生成: make openapi-gen
```

## Key Technical Decisions

### 依存性注入パターン
main.goで依存関係を構築し、各層に注入：
```go
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo, ...)
userHandler := handler.NewUserHandler(userService)
```

### 環境変数による設定管理
`internal/config/config.go` で環境変数を読み込み、デフォルト値を提供

### OpenAPI駆動開発
`api/openapi.yaml` を仕様の真実とし、コード生成でハンドラーの型安全性を確保

### マイグレーションファースト
データベーススキーマ変更は必ずマイグレーションファイルで管理（`db/migrations/`）

### Echoミドルウェア活用
- Logger: リクエストログ
- Recover: パニックリカバリ
- CORS: クロスオリジンリソース共有

---
_Document standards and patterns, not every dependency_
