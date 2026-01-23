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
- 以下のいずれか：
  - **Claude Pro/Maxサブスクリプション**（推奨：サブスクリプション内で動作）
  - **Anthropic APIキー**（従量課金、[Anthropic Console](https://console.anthropic.com/)で取得）

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

### 2. 認証方法の選択

#### 方法A: OAuth認証（Claude Pro/Maxサブスクリプションを使用）

**✨ サブスクリプションのクレジットで動作するため、追加料金不要！**

1. **ローカルでOAuthトークンを生成**:
   ```bash
   claude setup-token
   ```
   ※Claude Pro/Maxユーザーのみ利用可能

2. **GitHubシークレットに追加**:
   - GitHub Settings → Secrets and variables → Actions
   - **New repository secret**をクリック
   - 以下を設定：

   | Name | Value |
   |------|-------|
   | `CLAUDE_CODE_OAUTH_TOKEN` | 生成されたOAuthトークン |

#### 方法B: API Key認証（従量課金）

1. **Anthropic APIキーを取得**:
   - [Anthropic Console](https://console.anthropic.com/)でAPIキーを生成

2. **GitHubシークレットに追加**:
   - GitHub Settings → Secrets and variables → Actions
   - **New repository secret**をクリック
   - 以下を設定：

   | Name | Value |
   |------|-------|
   | `ANTHROPIC_API_KEY` | AnthropicコンソールのAPIキー |

#### 認証方法の比較

| 項目 | OAuth認証 | API Key認証 |
|------|---------|-----------|
| **必要なもの** | Claude Pro/Max | Anthropic APIキー |
| **料金** | サブスクリプション内 | 従量課金 |
| **セットアップ** | ローカルでトークン生成 | コンソールでキー取得 |
| **推奨ユーザー** | Pro/Maxユーザー | 一般ユーザー |

### 3. ワークフローの確認

ワークフローファイルは`.github/workflows/claude-review.yml`に配置されています。

## 使い方

### ラベルによる自動レビュー（推奨）

PRに`claude-review`ラベルを追加すると、Claude Code Actionが自動的にトリガーされます：

1. PRを作成
2. PRに`claude-review`ラベルを追加
3. GitHub Actionsが自動的に実行される
4. Claudeがコードをレビュー
5. レビュー結果がPRのコメントとして投稿される

#### ラベルの作成方法

リポジトリに`claude-review`ラベルがない場合は、以下の手順で作成してください：

1. リポジトリの**Issues**タブまたは**Pull requests**タブを開く
2. **Labels**をクリック
3. **New label**をクリック
4. 以下を設定：
   - **Label name**: `claude-review`
   - **Description**: Claude AIによるコードレビューをリクエスト
   - **Color**: お好みの色（例：紫 `#7B68EE`）
5. **Create label**をクリック

### @claudeメンションによるレビュー依頼

PRのコメントで`@claude`をメンションすることで、特定の質問やレビューを依頼できます：

```markdown
@claude このPRのセキュリティ面をレビューしてください
```

```markdown
@claude このマイグレーションファイルは正しいですか？
```

この方法は、ラベルを追加せずに特定の観点でレビューを依頼したい場合に便利です。

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

### Claude Code プロセスエラー（exit code 1）

```
error: Claude Code process exited with code 1
```

このエラーは認証やAPI接続の問題で発生することが多いです。

**対処法：**

1. **認証トークンの再生成**:
   ```bash
   # OAuth認証の場合
   claude setup-token

   # 生成されたトークンをGitHubシークレットに再設定
   ```

2. **シークレット名の確認**:
   - `CLAUDE_CODE_OAUTH_TOKEN` または `ANTHROPIC_API_KEY` の名前が正確か確認
   - 余分なスペースや改行が含まれていないか確認

3. **API制限の確認**:
   - Anthropicアカウントのクレジット残高を確認
   - レート制限に達していないか確認
   - [Anthropic Console](https://console.anthropic.com/)でAPIステータスを確認

4. **ワークフローログの確認**:
   - Actions → 該当のワークフロー実行 → "Claude Code Review"ステップ
   - `show_full_output: true` が有効になっているため、詳細なエラー情報を確認可能

5. **PRサイズの確認**:
   - 変更ファイル数が多すぎる場合は分割を検討
   - 大きなバイナリファイルが含まれていないか確認

### APIキーのエラー

```
Error: Invalid API key
```

- Anthropic APIキーが正しく設定されているか確認
- APIキーの有効期限が切れていないか確認
- Anthropicアカウントのクレジットが残っているか確認
- API keyが`sk-ant-`で始まっているか確認（OAuth tokenと混同していないか）

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

### 1. レビュータイミングの管理

- **準備ができたら`claude-review`ラベルを追加**: PRが下書き状態から準備完了になったタイミングでラベルを追加
- **CIが通過してからレビュー**: テスト、Lint、ビルドが成功してからClaudeレビューを実行すると効率的
- **コスト管理**: ラベルベースなので、必要な時だけレビューを実行でき、API利用料を抑えられる
- **複数回のレビュー**: ラベルを外して再度追加することで、修正後に再レビューを依頼可能

### 2. 適切な粒度でPRを作成

- 小さく、レビューしやすいPRを心がける
- 1つのPRで1つの機能や修正に絞る

### 3. コンテキストを提供

PRの説明欄に以下を記載：
- 変更の目的
- 関連するissueやタスク
- テスト方法
- 懸念点や確認してほしい箇所

### 4. Claudeのフィードバックを活用

- Claudeの指摘を真摯に受け止める
- 改善提案を積極的に取り入れる
- 不明点は`@claude`で質問する

### 5. 既存のCIと併用

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

### Q: `claude-review`ラベルを追加してもワークフローが実行されません。

A: 以下を確認してください：
1. ラベル名が正確に`claude-review`であること（大文字小文字を区別）
2. 認証シークレット（`CLAUDE_CODE_OAUTH_TOKEN`または`ANTHROPIC_API_KEY`）が設定されていること
3. GitHub Actionsの権限が有効になっていること（Settings → Actions → General → Workflow permissions）
4. ワークフローファイル`.github/workflows/claude-review.yml`がmainブランチにマージされていること

### Q: ラベルを使わずに自動的にレビューすることはできますか？

A: はい。ワークフローファイルのトリガーを`types: [opened, synchronize]`に変更することで、PRの作成・更新時に自動実行できます。ただし、API利用料が増加する可能性があるため、ラベルベースの方式を推奨しています。
