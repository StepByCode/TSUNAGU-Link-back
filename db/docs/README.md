# Database Documentation

このディレクトリには、データベース設計に関するドキュメントや参考資料を保存します。

## ファイル一覧

### `future_schema_draft.sql`

将来的なスキーマ設計案です。現在のOpenAPI仕様（`api/openapi.yaml`）には含まれていない、より複雑なマルチスキーマ設計が記述されています。

**含まれる設計要素**:
- `master_data` スキーマ: grades, organizations, locations
- `users` スキーマ: Zitadel統合を想定したユーザー管理
- `profiles` スキーマ: ユーザープロフィール情報（NFC、組織所属など）
- `attendance_logs` スキーマ: 入退室ログ管理

**注意**: このファイルは参考資料であり、現在のマイグレーションシステム（`db/migrations/`）では自動実行されません。

## 現在の本番スキーマ

現在の本番環境で使用されているスキーマは、`db/migrations/` ディレクトリのマイグレーションファイルで定義されています。

- `000001_create_users_table.up.sql`: シンプルなユーザーテーブル（`public.users`）

このスキーマは `api/openapi.yaml` のAPI仕様と整合性が取れています。

## 将来的な拡張

`future_schema_draft.sql` に記載されている設計を本番環境に適用する場合は、以下の手順を推奨します:

### 方法1: Atlasを使用（推奨）

1. OpenAPI仕様（`api/openapi.yaml`）を更新して新しいエンドポイントを定義
2. スキーマ定義ファイル（`../schema.hcl`）を更新して新しいテーブルを追加
3. Atlasでマイグレーションファイルを自動生成:
   ```bash
   make migrate-generate name=add_master_data_schema
   make migrate-generate name=add_profiles_schema
   make migrate-generate name=add_attendance_logs_schema
   ```
4. 生成されたマイグレーションファイルをレビュー
5. マイグレーションを実行:
   ```bash
   make migrate-up
   ```

### 方法2: 手動作成

1. OpenAPI仕様（`api/openapi.yaml`）を更新
2. 空のマイグレーションファイルを作成:
   ```bash
   make migrate-create name=add_master_data_schema
   make migrate-create name=add_profiles_schema
   make migrate-create name=add_attendance_logs_schema
   ```
3. 各マイグレーションファイルに適切なUP/DOWNスクリプトを手動で記述
4. マイグレーションを実行:
   ```bash
   make migrate-up
   ```

## マイグレーション管理

マイグレーションの作成・実行方法については、プロジェクトルートの `README.md` を参照してください。

このプロジェクトでは、**Atlas**を使用してスキーマ定義（`db/schema.hcl`）から
マイグレーションファイルを自動生成することができます。
