# スキーママイグレーション実施手順書

このドキュメントは、`db/schema.hcl`の更新からGit Pushまでの一連のワークフローを示します。

## 前提条件

- Dockerが起動していること
- Atlasがインストールされていること（`make install-tools`でインストール可能）
- 作業ブランチにチェックアウト済みであること

## ステップバイステップ手順

### 📋 Step 1: 現在の状態を確認

```bash
# 現在のブランチを確認
git branch --show-current

# マイグレーション状態を確認
make migrate-status

# 現在のDBスキーマを確認（オプション）
make schema-inspect
```

**チェックポイント:**
- [ ] 正しいブランチで作業していることを確認
- [ ] 既存のマイグレーションがすべて適用済みであることを確認

---

### ✏️ Step 2: スキーマ定義を編集

`db/schema.hcl`を開いて、必要な変更を加えます。

**変更例:**
- 新しいテーブルの追加
- 既存テーブルへのカラム追加
- インデックスの追加
- 外部キー制約の追加

```bash
# エディタでスキーマを編集
vim db/schema.hcl
# または
code db/schema.hcl
```

**チェックポイント:**
- [ ] スキーマ定義の構文が正しいことを確認
- [ ] カラムのnull制約、デフォルト値が適切に設定されていることを確認
- [ ] 外部キー制約がある場合、参照先のテーブル・カラムが正しいことを確認

---

### 🔧 Step 3: マイグレーションファイルを生成

変更内容を説明する適切な名前でマイグレーションファイルを生成します。

```bash
# 推奨: 名前を明示的に指定
make migrate-generate name=add_posts_table

# または: 対話的に名前を入力
make migrate-generate

# または: 完全自動（後でリネーム推奨）
make migrate-generate-auto
```

**命名規則の例:**
- `add_posts_table` - 新しいテーブルを追加
- `add_email_to_users` - usersテーブルにemailカラムを追加
- `add_index_user_email` - user_emailインデックスを追加
- `add_foreign_key_posts_user` - posts.user_idに外部キーを追加

**チェックポイント:**
- [ ] 以下のファイルが生成されたことを確認
  - `db/migrations/XXXXXX_[名前].up.sql`
  - `db/migrations/XXXXXX_[名前].down.sql`

---

### 👀 Step 4: 生成されたSQLファイルをレビュー

生成されたマイグレーションファイルの内容を必ず確認します。

```bash
# 生成されたファイルを確認
ls -la db/migrations/

# 最新のマイグレーションファイルを表示
tail -2 db/migrations/*.sql | head -20

# または、特定のファイルを直接確認
cat db/migrations/XXXXXX_add_posts_table.up.sql
cat db/migrations/XXXXXX_add_posts_table.down.sql
```

**レビューポイント:**
- [ ] `.up.sql`: 意図した変更が正しく記述されているか
- [ ] `.down.sql`: ロールバック処理が正しく記述されているか
- [ ] データの損失を引き起こす危険な操作（DROP, DELETE）が含まれていないか確認
- [ ] インデックス作成時に`CONCURRENTLY`オプションが必要かどうか検討（本番環境でのロックを避けるため）

---

### 🧪 Step 5: ローカル環境でマイグレーションをテスト

本番適用前に、必ずローカル環境でテストします。

```bash
# マイグレーションを適用
make migrate-up

# マイグレーション状態を確認
make migrate-status

# 適用後のスキーマを確認
make schema-inspect

# 必要に応じてアプリケーションを起動してテスト
make run
# または
make dev
```

**チェックポイント:**
- [ ] マイグレーションがエラーなく適用された
- [ ] `migrate-status`で新しいマイグレーションが適用済みになっている
- [ ] アプリケーションが正常に起動する
- [ ] 関連するAPIエンドポイントが正常に動作する

---

### 🔄 Step 6: ロールバックをテスト（重要）

ロールバックが正しく動作することを確認します。

```bash
# 最新のマイグレーションをロールバック
make migrate-down

# 状態を確認
make migrate-status

# 再度適用
make migrate-up
```

**チェックポイント:**
- [ ] ロールバックがエラーなく実行された
- [ ] 再度適用してもエラーが発生しない
- [ ] データの不整合が発生していない

---

### 📝 Step 7: 変更をコミット

すべてのテストが成功したら、変更をコミットします。

```bash
# 変更されたファイルを確認
git status

# スキーマ定義とマイグレーションファイルをステージング
git add db/schema.hcl
git add db/migrations/

# コミット（わかりやすいメッセージで）
git commit -m "$(cat <<'EOF'
feat(db): Add posts table with user relationship

- Add posts table with id, title, content, user_id, created_at columns
- Add foreign key constraint from posts.user_id to users.id
- Add index on posts.user_id for query performance
EOF
)"
```

**コミットメッセージの例:**
- `feat(db): Add posts table with user relationship`
- `feat(db): Add email column to users table`
- `feat(db): Add index on user email for faster lookups`
- `fix(db): Add missing NOT NULL constraint on posts.title`

**チェックポイント:**
- [ ] スキーマ定義ファイル（`db/schema.hcl`）がステージングされている
- [ ] マイグレーションファイル（`.up.sql`, `.down.sql`）がステージングされている
- [ ] コミットメッセージが変更内容を明確に説明している

---

### 🚀 Step 8: リモートリポジトリにプッシュ

```bash
# 現在のブランチにプッシュ
git push -u origin $(git branch --show-current)

# プッシュに失敗した場合（ネットワークエラー）、リトライ
# 2秒待ってリトライ
sleep 2 && git push -u origin $(git branch --show-current)
```

**チェックポイント:**
- [ ] プッシュが成功した
- [ ] ネットワークエラーが発生した場合は、2秒待ってリトライ（最大4回）

---

### ✅ Step 9: CI/CDとデプロイの確認

```bash
# GitHubでPRのステータスを確認
gh pr view

# CIの実行状況を確認
gh run list --limit 5
```

**チェックポイント:**
- [ ] CIパイプラインが正常に実行されている
- [ ] PRがマージ可能な状態になっている
- [ ] デプロイ後、本番環境でマイグレーションが自動実行される（Coolifyのentrypoint.sh経由）

---

## ⚠️ 注意事項

### 重要な原則

1. **`db/schema.hcl`が真実の源泉（Single Source of Truth）**
   - 常にスキーマ定義を先に編集してからマイグレーション生成
   - SQLファイルを手動で編集しない（Atlasの自動生成を信頼）

2. **生成されたSQLファイルは必ずレビュー**
   - Atlasは賢いが、意図しない変更が含まれる可能性がある
   - 特にDROP、DELETEを含む変更は慎重に確認

3. **ロールバックも必ずテスト**
   - `.down.sql`が正しく動作することを確認
   - 本番環境で問題が発生した際の復旧手段を確保

4. **本番環境への影響を考慮**
   - 大規模テーブルへのカラム追加は時間がかかる可能性
   - インデックス作成は`CONCURRENTLY`オプションを検討
   - 外部キー制約の追加は既存データの整合性を確認

5. **デプロイ = 自動マイグレーション実行**
   - Coolifyへのデプロイ時、`entrypoint.sh`が自動的にマイグレーションを実行
   - マイグレーションファイルの品質は極めて重要

### よくある問題と対処法

#### 問題: マイグレーション生成時にエラーが発生

```bash
# スキーマ定義の構文を検証
atlas schema inspect --env schema

# エラーメッセージを確認してdb/schema.hclを修正
```

#### 問題: マイグレーション適用時にエラーが発生

```bash
# ロールバック
make migrate-down

# スキーマ定義を修正
vim db/schema.hcl

# マイグレーションファイルを削除
rm db/migrations/XXXXXX_*

# 再生成
make migrate-generate name=修正版の名前

# 再適用
make migrate-up
```

#### 問題: プッシュ時にネットワークエラー

```bash
# 2秒待機後にリトライ
sleep 2 && git push -u origin $(git branch --show-current)

# それでも失敗する場合は4秒待機
sleep 4 && git push -u origin $(git branch --show-current)

# 最大4回までリトライ（2s, 4s, 8s, 16sの間隔）
```

---

## 📚 関連ドキュメント

- [ATLAS_GUIDE.md](./ATLAS_GUIDE.md) - Atlasの詳細な使い方
- [CLAUDE.md](../CLAUDE.md) - AI-DLCとSpec-Driven Developmentのガイド
- [Makefile](../Makefile) - 利用可能なmakeコマンド一覧

---

## 🎯 クイックリファレンス

```bash
# 1. スキーマ編集
vim db/schema.hcl

# 2. マイグレーション生成
make migrate-generate name=your_migration_name

# 3. レビュー
cat db/migrations/XXXXXX_your_migration_name.up.sql
cat db/migrations/XXXXXX_your_migration_name.down.sql

# 4. テスト適用
make migrate-up

# 5. ロールバックテスト
make migrate-down && make migrate-up

# 6. コミット
git add db/schema.hcl db/migrations/
git commit -m "feat(db): Your descriptive message"

# 7. プッシュ
git push -u origin $(git branch --show-current)
```

---

最終更新: 2026-01-16
