# Project Structure

## Organization Philosophy

**レイヤードアーキテクチャ + internal パッケージパターン**: 責任を明確に分離し、外部への公開を制限することで保守性とカプセル化を実現

```
エントリーポイント (cmd/)
    ↓
Handler層 (internal/handler/) : HTTPリクエスト処理
    ↓
Service層 (internal/service/) : ビジネスロジック
    ↓
Repository層 (internal/repository/) : データアクセス
    ↓
データベース
```

## Directory Patterns

### エントリーポイント
**Location**: `cmd/server/`
**Purpose**: アプリケーションの起動、依存性の構築・注入
**Example**: `main.go` でDB接続、リポジトリ/サービス/ハンドラーのインスタンス化、Echoサーバー起動

### Handler層 (HTTPインターフェース)
**Location**: `internal/handler/`
**Purpose**: HTTPリクエスト/レスポンス処理、入力バリデーション、Serviceレイヤーへの委譲
**Example**: `user_handler.go` で `/api/v1/users` エンドポイントを実装、Echoにルート登録

### Service層 (ビジネスロジック)
**Location**: `internal/service/`
**Purpose**: ビジネスルール、トランザクション制御、複数リポジトリの調整
**Example**: `user_service.go` でユーザー作成時のパスワードハッシュ化、JWT生成

### Repository層 (データアクセス)
**Location**: `internal/repository/`
**Purpose**: データベースCRUD操作、SQLクエリ実行
**Example**: `user_repository.go` で `users` テーブルへのSQL操作をカプセル化

### Model/Entity
**Location**: `internal/model/`
**Purpose**: ドメインモデル、データ構造定義
**Example**: `user.go` で `User` 構造体定義

### Configuration
**Location**: `internal/config/`
**Purpose**: 環境変数読み込み、設定管理
**Example**: `config.go` で `Config` 構造体と `Load()` 関数

### API仕様
**Location**: `api/`
**Purpose**: OpenAPI仕様書、API契約の定義
**Example**: `openapi.yaml` から `oapi-codegen` で型を生成

### Database Migrations
**Location**: `db/migrations/`
**Purpose**: データベーススキーマのバージョン管理
**Pattern**: `{seq}_{name}.up.sql` / `{seq}_{name}.down.sql`

## Naming Conventions

- **Files**: スネークケース (`user_handler.go`, `user_service.go`)
- **Packages**: 小文字単数形 (`handler`, `service`, `repository`)
- **Types**: パスカルケース (`UserHandler`, `UserService`)
- **Functions/Methods**: キャメルケース (`NewUserHandler()`, `CreateUser()`)
- **Constructors**: `New{Type}` パターン (`NewUserHandler`, `NewUserService`)

## Import Organization

```go
import (
    // 標準ライブラリ
    "database/sql"
    "fmt"
    "log"

    // 内部パッケージ (プロジェクト内)
    "github.com/StepByCode/TSUNAGU-Link-back/internal/config"
    "github.com/StepByCode/TSUNAGU-Link-back/internal/handler"

    // サードパーティ
    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"  // ドライバーは blank import
)
```

**Import Order**: 標準ライブラリ → 内部パッケージ → サードパーティ

**Module Path**: `github.com/StepByCode/TSUNAGU-Link-back`

## Code Organization Principles

### Dependency Rule（依存性の方向）
```
Handler → Service → Repository
  ↓         ↓          ↓
外側      中間層      内側
```
依存は外側から内側へのみ。逆方向の依存は禁止（Repository は Service に依存しない）

### Interface Driven Design
各層はインターフェースを通じて連携し、依存性注入で実装を差し替え可能に

### Single Responsibility
1ファイル1責務（例: `user_handler.go` はユーザー関連のHTTPハンドラーのみ）

### Package Cohesion
関連する機能を同一パッケージに配置（例: `handler` パッケージ内に全てのHTTPハンドラー）

### Internal Package Pattern
`internal/` 配下のコードは外部モジュールから import 不可（Goの言語機能）

---
_Document patterns, not file trees. New files following patterns shouldn't require updates_
