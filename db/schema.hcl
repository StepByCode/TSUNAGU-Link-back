// TSUNAGU-Link データベーススキーマ定義
// このファイルは、Atlasを使用してマイグレーションファイルを自動生成するための
// スキーマ定義です。既存のマイグレーション（000001_create_users_table.up.sql）
// の内容を反映しています。

schema "public" {
}

// usersテーブル
table "users" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }

  column "email" {
    null = false
    type = varchar(255)
  }

  column "name" {
    null = false
    type = varchar(255)
  }

  column "password" {
    null = false
    type = varchar(255)
  }

  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }

  column "deleted_at" {
    null = true
    type = timestamp
  }

  primary_key {
    columns = [column.id]
  }

  index "users_email_key" {
    unique  = true
    columns = [column.email]
  }

  index "idx_users_email" {
    columns = [column.email]
    where   = "deleted_at IS NULL"
  }

  index "idx_users_deleted_at" {
    columns = [column.deleted_at]
  }
}
