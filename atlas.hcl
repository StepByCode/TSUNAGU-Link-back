// Atlas configuration for TSUNAGU-Link-back

// 開発環境の設定
env "dev" {
  // 現在のデータベース状態（既存のマイグレーション適用済み）
  url = "postgres://tsunagu:tsunagu_password@localhost:5432/tsunagu_db?sslmode=disable"

  // マイグレーションディレクトリ（golang-migrate形式）
  migration {
    dir = "file://db/migrations?format=golang-migrate"
  }

  // スキーマソース（スキーマ定義ファイルから読み込み）
  src = "file://db/schema.hcl"
}

// 本番環境の設定（環境変数から読み込み）
env "prod" {
  // 環境変数から接続情報を取得
  url = getenv("DATABASE_URL")

  // 開発用データベース（スキーマ計算用）
  dev = "docker://postgres/16/dev?search_path=public"

  // マイグレーションディレクトリ（golang-migrate形式）
  migration {
    dir = "file://db/migrations?format=golang-migrate"
  }
}

// スキーマファイルを使用する場合の設定
env "schema" {
  // 既存のマイグレーションを基準にする
  url = "file://db/migrations?format=golang-migrate"

  // 開発用データベース（スキーマ計算用）
  dev = "docker://postgres/16/dev"

  // マイグレーションディレクトリ（golang-migrate形式）
  migration {
    dir = "file://db/migrations?format=golang-migrate"
  }

  // スキーマ定義ファイルから読み込み
  src = "file://db/schema.hcl"
}
