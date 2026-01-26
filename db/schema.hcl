// TSUNAGU-Link データベーススキーマ定義
// このファイルは、Atlasを使用してマイグレーションファイルを自動生成するための
// スキーマ定義です。既存のマイグレーション（db/docs/future_schema_draft.sql）
// の内容を反映しています。

schema "public" {
}
# スキーマの定義
schema "master_data" {}
schema "users" {}
schema "profiles" {}
schema "attendance_logs" {}
# 1. UUID機能（拡張）
# ※ PostgreSQLの拡張機能はAtlasの設定で管理するかSQLで事前実行が必要です
# 2. master_data テーブル群
table "grades" {
  schema = schema.master_data
  column "id" {
    type = serial
  }
  column "label" {
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
table "locations" {
  schema = schema.master_data
  column "id" {
    type = serial
  }
  column "name" {
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
# 3. ユーザーテーブル
table "users" {
  schema = schema.users
  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "email" {
    type = text
  }
  column "password" {
    type = text
  }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  index "users_email_key" {
    unique  = true
    columns = [column.email]
  }
}
# 4. プロフィールテーブル（外部キー制約の例）
table "profiles" {
  schema = schema.profiles
  column "user_id" {
    type = uuid
  }
  column "name" {
    type = text
  }
  column "grade_id" {
    type = int
  }
  primary_key {
    columns = [column.user_id]
  }
  # これが外部キー制約 (Foreign Key)
  foreign_key "profiles_user_id_fkey" {
    columns     = [column.user_id]         # 自分のテーブルのカラム
    ref_columns = [table.users.column.id] # 相手のテーブルのカラム
    on_delete   = CASCADE                 # ユーザー消えたらプロフィールも消す
  }
  foreign_key "profiles_grade_id_fkey" {
    columns     = [column.grade_id]
    ref_columns = [table.grades.column.id]
  }
}