# Claude Code Action セットアップガイド

このガイドでは、TSUNAGU-Link-backプロジェクトにおけるClaude Code Actionの設定と使い方について説明します。

## 概要

Claude Code Actionは、AnthropicのClaude AIを使用してプルリクエストを自動的にレビューするGitHub Actionsワークフローです。

### 主な機能

- 🔍 **自動コードレビュー**: PRの作成・更新時に自動的にレビューを実行
- 🤖 **@claudeメンション**: PRのコメントで`@claude`をメンションしてレビューを依頼
- 🛡️ **セキュリティチェック**: OWASP Top 10などのセキュリティ脆弱性を検出
- 📊 **品質チェック**: Goのベストプラクティス、パフォーマンス、テストカバレッジを確認
- 🗄️ **データベースレビュー**: Atlasマイグレーションファイルの妥当性を検証
- 📝 **仕様整合性**: Kiro-style SDD仕様との整合性を確認

## セットアップ手順

### 前提条件

- リポジトリの管理者権限
- Anthropic APIキー（[Anthropic Console](https://console.anthropic.com/)で取得）

### 1. Claude Code CLIでセットアップ（推奨）

最も簡単な方法は、Claude Code CLIを使用することです：

```bash
# Claude Code CLIを起動
claude

# GitHub Appインストールコマンドを実行
/install-github-app
```

このコマンドは以下を自動的に行います：
- GitHub Appのインストール
- 必要な権限の設定
- シークレットの設定ガイド

### 2. Anthropic APIキーの設定

1. GitHubリポジトリの設定を開く
2. **Settings** → **Secrets and variables** → **Actions**
3. **New repository secret**をクリック
4. 以下のシークレットを追加：

| Name | Value |
|------|-------|
| `ANTHROPIC_API_KEY` | AnthropicコンソールのAPIキー |

### 3. ワークフローの確認

ワークフローファイルは`.github/workflows/claude-review.yml`に配置されています。

## 使い方

### 自動レビュー

PRを作成または更新すると、Claude Code Actionが自動的にトリガーされます：

1. PRを作成
2. GitHub Actionsが自動的に実行される
3. Claudeがコードをレビュー
4. レビュー結果がPRのコメントとして投稿される

### 手動レビュー依頼

PRのコメントで`@claude`をメンションすることで、特定の質問やレビューを依頼できます：

```markdown
@claude このPRのセキュリティ面をレビューしてください
```

```markdown
@claude このマイグレーションファイルは正しいですか？
```

### レビュー内容

Claudeは以下の観点でレビューを行います：

#### 1. コード品質
- Goのベストプラクティスとイディオム
- コードの可読性と保守性
- エラーハンドリング
- リソース管理（defer, Close等）

#### 2. セキュリティ
- SQLインジェクション、XSS等のOWASP Top 10脆弱性
- 認証・認可の実装
- 機密情報の取り扱い
- 入力値の検証

#### 3. パフォーマンス
- 不要なメモリアロケーション
- N+1クエリ問題
- ゴルーチンリーク
- データベースインデックス

#### 4. テスト
- テストカバレッジ
- エッジケースのテスト
- モックの適切な使用

#### 5. データベース
- Atlasマイグレーションの正確性
- schema.hclとの整合性
- ロールバック可能性
- 外部キー制約とインデックス

#### 6. 仕様駆動開発（SDD）
- `.kiro/specs/`配下の仕様との整合性
- 要件・設計・タスクとの対応
- ステアリング文書との一貫性

## トラブルシューティング

### ワークフローが実行されない

1. **権限の確認**:
   - Settings → Actions → General → Workflow permissions
   - "Read and write permissions"が有効になっているか確認

2. **シークレットの確認**:
   - `ANTHROPIC_API_KEY`が正しく設定されているか確認

3. **ワークフローファイルの確認**:
   - `.github/workflows/claude-review.yml`が存在するか確認
   - YAMLの構文エラーがないか確認

### APIキーのエラー

```
Error: Invalid API key
```

- Anthropic APIキーが正しく設定されているか確認
- APIキーの有効期限が切れていないか確認
- Anthropicアカウントのクレジットが残っているか確認

### レビューが投稿されない

1. **PRの権限確認**:
   ```yaml
   permissions:
     pull-requests: write
   ```
   がワークフローに含まれているか確認

2. **ログの確認**:
   - GitHub Actions → Workflowsタブ
   - 失敗したワークフローのログを確認

## カスタマイズ

### レビュープロンプトの変更

`.github/workflows/claude-review.yml`の`prompt`セクションを編集することで、レビューの観点をカスタマイズできます：

```yaml
- name: Claude Code Review
  uses: anthropics/claude-code-action@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
    prompt: |
      あなたのカスタムプロンプトをここに記述
```

### トリガー条件の変更

特定のファイルパターンやブランチでのみ実行するように設定できます：

```yaml
on:
  pull_request:
    types: [opened, synchronize]
    paths:
      - 'internal/**'
      - 'cmd/**'
    branches:
      - main
      - develop
```

### モデルの変更

デフォルトではSonnetモデルが使用されますが、Opusモデルに変更することも可能です：

```yaml
- name: Claude Code Review
  uses: anthropics/claude-code-action@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
    claude_args: '{"model": "claude-opus-4-5-20251101"}'
    prompt: |
      ...
```

## ベストプラクティス

### 1. 適切な粒度でPRを作成

- 小さく、レビューしやすいPRを心がける
- 1つのPRで1つの機能や修正に絞る

### 2. コンテキストを提供

PRの説明欄に以下を記載：
- 変更の目的
- 関連するissueやタスク
- テスト方法
- 懸念点や確認してほしい箇所

### 3. Claudeのフィードバックを活用

- Claudeの指摘を真摯に受け止める
- 改善提案を積極的に取り入れる
- 不明点は`@claude`で質問する

### 4. 既存のCIと併用

Claude Code Actionは既存のCI（テスト、Lint、ビルド）を補完するものです：
- CIで機械的なチェック（構文、テスト実行）
- Claudeで論理的なレビュー（設計、セキュリティ）

## 参考リンク

- [Claude Code公式ドキュメント](https://code.claude.com/docs/ja/github-actions)
- [Claude Code Action GitHubリポジトリ](https://github.com/anthropics/claude-code-action)
- [Anthropic API Documentation](https://docs.anthropic.com/)
- [GitHub Actions Documentation](https://docs.github.com/ja/actions)

## よくある質問（FAQ）

### Q: レビューの実行にコストはかかりますか？

A: はい。Claude Code ActionはAnthropic APIを使用するため、API利用料金が発生します。料金の詳細は[Anthropicの料金ページ](https://www.anthropic.com/pricing)を参照してください。

### Q: プライベートリポジトリでも使用できますか？

A: はい。プライベートリポジトリでも使用可能です。ただし、GitHub Actionsの実行時間に制限がある場合があります。

### Q: 他のチームメンバーもClaudeに質問できますか？

A: はい。PRにアクセス権のあるすべてのユーザーが`@claude`メンションで質問できます。

### Q: レビュー結果を日本語で受け取れますか？

A: はい。このプロジェクトのワークフローでは、プロンプトで日本語でのレビューを指定しているため、レビュー結果は日本語で提供されます。

### Q: 特定のファイルだけレビューしてもらうことはできますか？

A: はい。`@claude`メンションで特定のファイルを指定してレビューを依頼できます：
```markdown
@claude internal/handler/user.go のセキュリティをチェックしてください
```
