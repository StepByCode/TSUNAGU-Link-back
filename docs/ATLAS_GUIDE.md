# Atlas マイグレーション自動生成ガイド

このドキュメントでは、Atlasを使用してデータベースマイグレーションファイルを自動生成する方法を説明します。

## 概要

Atlasは、スキーマ定義ファイル（`db/schema.hcl`）からデータベースマイグレーションファイルを自動生成するツールです。このプロジェクトでは、golang-migrate互換の形式でマイグレーションファイルを生成します。

## セットアップ

### 1. Atlasのインストール

```bash
make install-tools
```

または、直接インストール:

```bash
curl -sSf https://atlasgo.sh | sh
```

### 2. Dockerの起動

Atlasは開発用データベースを使用してスキーマの差分を計算します:

```bash
make docker-up
```

## 基本的な使い方

### マイグレーションの自動生成

1. **スキーマ定義ファイルを編集**

   `db/schema.hcl` を開いて、新しいテーブルやカラムを追加します。

   例: `posts` テーブルを追加
   ```hcl
   table "posts" {
     schema = schema.public

     column "id" {
       null    = false
       type    = uuid
       default = sql("uuid_generate_v4()")
     }

     column "title" {
       null = false
       type = varchar(255)
     }

     column "content" {
       null = false
       type = text
     }

     column "user_id" {
       null = false
       type = uuid
     }

     column "created_at" {
       null    = false
       type    = timestamp
       default = sql("CURRENT_TIMESTAMP")
     }

     primary_key {
       columns = [column.id]
     }

     foreign_key "posts_user_id_fkey" {
       columns     = [column.user_id]
       ref_columns = [table.users.column.id]
       on_delete   = CASCADE
     }

     index "idx_posts_user_id" {
       columns = [column.user_id]
     }
   }
   ```

2. **マイグレーションファイルを生成**

   ```bash
   make migrate-generate name=create_posts_table
   ```

   これにより、以下のファイルが自動生成されます:
   - `db/migrations/000002_create_posts_table.up.sql`
   - `db/migrations/000002_create_posts_table.down.sql`

3. **生成されたファイルを確認**

   ```bash
   cat db/migrations/000002_create_posts_table.up.sql
   ```

4. **マイグレーションを実行**

   ```bash
   make migrate-up
   ```

## 便利なコマンド

### マイグレーション状態の確認

現在のデータベースに適用されているマイグレーションを確認:

```bash
make migrate-status
```

### マイグレーションファイルの検証

マイグレーションファイルに問題がないか検証:

```bash
make migrate-lint
```

### データベーススキーマの確認

現在のデータベースのスキーマを表示:

```bash
make schema-inspect
```

## スキーマ定義の書き方

### 基本的なカラムタイプ

```hcl
column "example_varchar" {
  type = varchar(255)
  null = false
}

column "example_int" {
  type = integer
  null = true
}

column "example_text" {
  type = text
  null = false
}

column "example_bool" {
  type = boolean
  null = false
  default = false
}

column "example_timestamp" {
  type = timestamp
  null = false
  default = sql("CURRENT_TIMESTAMP")
}

column "example_uuid" {
  type = uuid
  null = false
  default = sql("uuid_generate_v4()")
}
```

### インデックス

```hcl
index "idx_example_column" {
  columns = [column.example_column]
}

# 複合インデックス
index "idx_multiple_columns" {
  columns = [column.col1, column.col2]
}

# 条件付きインデックス
index "idx_conditional" {
  columns = [column.email]
  where   = "deleted_at IS NULL"
}

# ユニークインデックス
index "idx_unique_email" {
  unique  = true
  columns = [column.email]
}
```

### 外部キー

```hcl
foreign_key "fk_user_id" {
  columns     = [column.user_id]
  ref_columns = [table.users.column.id]
  on_delete   = CASCADE
  on_update   = NO_ACTION
}
```

### CHECK制約

```hcl
check "positive_price" {
  expr = "price > 0"
}
```

## ワークフロー例

### 新しいテーブルを追加

1. `db/schema.hcl` に新しいテーブル定義を追加
2. `make migrate-generate name=add_xxx_table`
3. 生成されたマイグレーションファイルを確認
4. `make migrate-up` で適用
5. 必要に応じて `make migrate-down` でロールバック

### 既存のテーブルにカラムを追加

1. `db/schema.hcl` の該当テーブルに新しいカラムを追加
2. `make migrate-generate name=add_xxx_column`
3. 生成されたマイグレーションファイルを確認
4. `make migrate-up` で適用

### インデックスを追加

1. `db/schema.hcl` の該当テーブルにインデックス定義を追加
2. `make migrate-generate name=add_xxx_index`
3. 生成されたマイグレーションファイルを確認
4. `make migrate-up` で適用

## トラブルシューティング

### Atlasがインストールされていない

```bash
make install-tools
```

### Dockerコンテナが起動していない

```bash
make docker-up
```

### スキーマ定義ファイルの構文エラー

```bash
atlas schema inspect --env schema
```

エラーメッセージを確認して、`db/schema.hcl` の構文を修正してください。

### 既存のデータベースとスキーマ定義が一致しない

現在のデータベーススキーマを `db/schema.hcl` に反映:

```bash
atlas schema inspect --env dev > db/schema.hcl.new
```

生成された `db/schema.hcl.new` を確認して、必要に応じて `db/schema.hcl` にマージしてください。

## 参考資料

- [Atlas公式ドキュメント](https://atlasgo.io/docs)
- [golang-migrate統合ガイド](https://atlasgo.io/guides/migration-tools/golang-migrate)
- [HCLスキーマリファレンス](https://atlasgo.io/atlas-schema/sql-resources)
