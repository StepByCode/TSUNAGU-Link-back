// TSUNAGU-Link データベーススキーマ定義
// このファイルは、Atlasを使用してマイグレーションファイルを自動生成するための
// スキーマ定義です。

# =============================================================================
# スキーマ定義
# =============================================================================
schema "public" {}
schema "master_data" {}
schema "users" {}
schema "profiles" {}
schema "attendance_logs" {}

# =============================================================================
# Enum定義
# =============================================================================
enum "attendance_type" {
  schema = schema.attendance_logs
  values = ["check-in", "check-out"]
}

# =============================================================================
# 1. master_data テーブル群
# =============================================================================

# 学年マスタ
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

# 拠点マスタ
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

# 組織マスタ
table "organizations" {
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

# =============================================================================
# 2. users テーブル（認証基盤）
# =============================================================================
table "users" {
  schema = schema.users
  column "id" {
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "zitadel_id" {
    type = uuid
    null = true
  }
  column "email" {
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
  index "users_zitadel_id_key" {
    unique  = true
    columns = [column.zitadel_id]
  }
}

# =============================================================================
# 3. profiles テーブル（生徒情報）
# =============================================================================
table "profiles" {
  schema = schema.profiles
  column "user_id" {
    type = uuid
  }
  column "display_id" {
    type = text
  }
  column "name" {
    type = text
  }
  column "nfc_serial" {
    type = text
    null = true
  }
  column "grade_id" {
    type = int
  }
  column "org_id" {
    type = int
  }
  primary_key {
    columns = [column.user_id]
  }
  index "profiles_display_id_key" {
    unique  = true
    columns = [column.display_id]
  }
  index "profiles_nfc_serial_key" {
    unique  = true
    columns = [column.nfc_serial]
  }
  foreign_key "profiles_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "profiles_grade_id_fkey" {
    columns     = [column.grade_id]
    ref_columns = [table.grades.column.id]
  }
  foreign_key "profiles_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organizations.column.id]
  }
}

# =============================================================================
# 4. attendance_logs テーブル（打刻データ）
# =============================================================================
table "attendance_logs" {
  schema = schema.attendance_logs
  column "id" {
    type = bigserial
  }
  column "user_id" {
    type = uuid
  }
  column "location_id" {
    type = int
  }
  column "type" {
    type = enum.attendance_type
  }
  column "timestamp" {
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  index "attendance_logs_user_id_idx" {
    columns = [column.user_id]
  }
  index "attendance_logs_timestamp_idx" {
    columns = [column.timestamp]
  }
  foreign_key "attendance_logs_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = CASCADE
  }
  foreign_key "attendance_logs_location_id_fkey" {
    columns     = [column.location_id]
    ref_columns = [table.locations.column.id]
  }
}
