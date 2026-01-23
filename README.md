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

### 手動でマイグレーションを作成

空のマイグレーションファイルを作成して、手動でSQLを記述:

```bash
make migrate-create name=add_new_table
```

### Atlasで自動生成（推奨）

スキーマ定義ファイル（`db/schema.hcl`）を編集してから、マイグレーションを自動生成:

```bash
# 1. db/schema.hcl を編集して新しいテーブルやカラムを追加

# 2a. マイグレーション名を指定して生成
make migrate-generate name=add_new_table

# 2b. 対話的に名前を入力（nameを省略）
make migrate-generate

# 2c. 完全自動（タイムスタンプ名で生成、後でリネーム推奨）
make migrate-generate-auto

# 3. 生成されたマイグレーションファイルを確認
# 4. 問題なければ実行
make migrate-up
```

### マイグレーションを実行

```bash
make migrate-up
```

### マイグレーションをロールバック

```bash
make migrate-down
```

### その他のマイグレーションコマンド

```bash
# マイグレーション状態を確認
make migrate-status

# マイグレーションファイルを検証
make migrate-lint

# 現在のデータベーススキーマを表示
make schema-inspect
```

### マイグレーションツール

このプロジェクトは以下のツールを使用しています:

- **golang-migrate**: マイグレーションの実行
- **Atlas**: マイグレーションファイルの自動生成と検証

詳細は以下を参照してください:
- [Atlasマイグレーション自動生成ガイド](docs/ATLAS_GUIDE.md)
- 設定ファイル: `atlas.hcl`
- スキーマ定義: `db/schema.hcl`

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

## CI/CD

このプロジェクトではGitHub Actionsを使用したCI/CDパイプラインを構築しています。

### 自動実行されるワークフロー

#### CI（継続的インテグレーション）
- **トリガー**: PRの作成・更新、main/developブランチへのpush
- **実行内容**:
  - テスト実行（PostgreSQL付き）
  - Lintチェック（golangci-lint）
  - ビルド検証
  - カバレッジレポート（Codecov）

#### Claude PR Review（自動コードレビュー）
- **トリガー**: PRに`claude-review`ラベルを追加、または`@claude`メンション
- **実行内容**:
  - コード品質チェック
  - セキュリティ脆弱性の検出
  - Goベストプラクティスの確認
  - データベースマイグレーションの検証
  - 仕様駆動開発（SDD）との整合性確認

### Claude Code Actionのセットアップ

Claude AIによる自動PRレビューを有効にするには：

1. **リポジトリ管理者がセットアップを実行**:
   ```bash
   claude
   /install-github-app
   ```

2. **認証方法を選択してシークレットを設定**:

   **方法A: Claude Pro/Maxサブスクリプションを使用（推奨）**
   ```bash
   # OAuthトークンを生成
   claude setup-token
   ```
   - GitHub Settings → Secrets and variables → Actions
   - `CLAUDE_CODE_OAUTH_TOKEN`を追加（生成されたトークン）
   - サブスクリプション内で動作するため、追加料金不要

   **方法B: API Key認証（従量課金）**
   - GitHub Settings → Secrets and variables → Actions
   - `ANTHROPIC_API_KEY`を追加（Anthropic APIキー）

3. **レビューの依頼方法**:

   **方法A: ラベルを追加（推奨）**
   - PRに`claude-review`ラベルを追加するとレビューが自動実行されます
   - ラベルがない場合は[リポジトリ設定](https://github.com/settings)でラベルを作成してください

   **方法B: @claudeメンション**
   ```
   @claude このPRをレビューしてください
   ```

詳細は[docs/CLAUDE_CODE_SETUP.md](docs/CLAUDE_CODE_SETUP.md)または[Claude Code GitHub Actions](https://code.claude.com/docs/ja/github-actions)を参照してください。

## 今後の予定

- [ ] Zitadel認証基盤の統合
- [ ] より多くのAPIエンドポイント
- [ ] 単体テスト・統合テストの充実
- [x] CI/CDパイプラインの構築
