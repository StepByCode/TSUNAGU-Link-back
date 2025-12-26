# TSUNAGU-Link-back

TSUNAGU Linkのバックエンドリポジトリです

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Echo v4
- **データベース**: PostgreSQL 16
- **認証**: JWT (将来的にZitadelを統合予定)
- **API仕様**: OpenAPI 3.0

## プロジェクト構成

```
.
├── cmd/
│   └── server/          # エントリーポイント
├── internal/
│   ├── config/          # 設定管理
│   ├── handler/         # HTTPハンドラ
│   ├── model/           # データモデル
│   ├── repository/      # データアクセス層
│   └── service/         # ビジネスロジック層
├── api/
│   └── openapi.yaml     # OpenAPI仕様書
├── db/
│   └── migrations/      # データベースマイグレーション
├── docker-compose.yml   # Docker構成
├── Dockerfile           # アプリケーションコンテナ
└── Makefile             # タスクランナー
```

## セットアップ

### 必要要件

- Go 1.21以上
- Docker & Docker Compose
- Make

### 環境変数

`.env.example`をコピーして`.env`を作成:

```bash
cp .env.example .env
```

### 初期セットアップ

プロジェクトのセットアップ（ツールのインストール、依存関係のダウンロード、Docker起動、マイグレーション実行）:

```bash
make setup
```

## 使い方

### Docker環境での起動

```bash
# PostgreSQLを起動
make docker-up

# アプリケーションをビルドして実行
make run
```

### 開発モード（ホットリロード）

```bash
make dev
```

### その他のコマンド

```bash
# ヘルプを表示
make help

# ビルド
make build

# テスト実行
make test

# コードフォーマット
make fmt

# Dockerログを表示
make docker-logs

# Dockerを停止
make docker-down

# マイグレーション作成
make migrate-create name=create_example_table

# マイグレーション実行
make migrate-up

# マイグレーションロールバック
make migrate-down

# OpenAPIからコード生成
make openapi-gen
```

## API エンドポイント

### ヘルスチェック
- `GET /health` - サービスのヘルスチェック

### 認証
- `POST /api/v1/auth/login` - ログイン

### ユーザー管理
- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users` - ユーザー一覧取得
- `GET /api/v1/users/:id` - ユーザー詳細取得
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除

詳細は `api/openapi.yaml` を参照してください。

## データベースマイグレーション

マイグレーションファイルは `db/migrations/` ディレクトリに配置されています。

新しいマイグレーションを作成:

```bash
make migrate-create name=add_new_table
```

マイグレーションを実行:

```bash
make migrate-up
```

## 開発ガイドライン

### ディレクトリ構成の原則

- `cmd/`: アプリケーションのエントリーポイント
- `internal/`: 外部に公開しないパッケージ
- `api/`: API仕様書
- `db/migrations/`: データベーススキーマの変更履歴

### レイヤーアーキテクチャ

1. **Handler層**: HTTPリクエスト/レスポンスの処理
2. **Service層**: ビジネスロジック
3. **Repository層**: データアクセス

## 今後の予定

- [ ] Zitadel認証基盤の統合
- [ ] より多くのAPIエンドポイント
- [ ] 単体テスト・統合テストの充実
- [ ] CI/CDパイプラインの構築
