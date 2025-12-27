# AI-DLC and Spec-Driven Development

Kiro-style Spec Driven Development implementation on AI-DLC (AI Development Life Cycle)

## Project Context

### Paths
- Steering: `.kiro/steering/`
- Specs: `.kiro/specs/`

### Steering vs Specification

**Steering** (`.kiro/steering/`) - Guide AI with project-wide rules and context
**Specs** (`.kiro/specs/`) - Formalize development process for individual features

### Active Specifications
- Check `.kiro/specs/` for active specifications
- Use `/kiro:spec-status [feature-name]` to check progress

## Development Guidelines
- Think in English, generate responses in Japanese. All Markdown content written to project files (e.g., requirements.md, design.md, tasks.md, research.md, validation reports) MUST be written in the target language configured for this specification (see spec.json.language).

## Minimal Workflow
- Phase 0 (optional): `/kiro:steering`, `/kiro:steering-custom`
- Phase 1 (Specification):
  - `/kiro:spec-init "description"`
  - `/kiro:spec-requirements {feature}`
  - `/kiro:validate-gap {feature}` (optional: for existing codebase)
  - `/kiro:spec-design {feature} [-y]`
  - `/kiro:validate-design {feature}` (optional: design review)
  - `/kiro:spec-tasks {feature} [-y]`
- Phase 2 (Implementation): `/kiro:spec-impl {feature} [tasks]`
  - `/kiro:validate-impl {feature}` (optional: after implementation)
- Progress check: `/kiro:spec-status {feature}` (use anytime)

## Development Rules
- 3-phase approval workflow: Requirements → Design → Tasks → Implementation
- Human review required each phase; use `-y` only for intentional fast-track
- Keep steering current and verify alignment with `/kiro:spec-status`
- Follow the user's instructions precisely, and within that scope act autonomously: gather the necessary context and complete the requested work end-to-end in this run, asking questions only when essential information is missing or the instructions are critically ambiguous.

## Steering Configuration
- Load entire `.kiro/steering/` as project memory
- Default files: `product.md`, `tech.md`, `structure.md`
- Custom files are supported (managed via `/kiro:steering-custom`)

## Subagent-based Quality Assurance

This project uses Claude Code's subagent functionality for specialized validation tasks.

### Available Subagents (`.claude/agents/`)

1. **spec-validator** - 要件・仕様検証
   - 仕様の完全性・明確性・実現可能性をチェック
   - ステアリング文書との整合性検証
   - SDD全フェーズ対応
   - 使用例: "spec-validatorで{機能}の要件を検証して"

2. **design-reviewer** - 技術設計レビュー
   - アーキテクチャの健全性評価
   - 要件トレーサビリティの確認
   - 実装リスクの特定
   - GO/NO-GO判定と根拠提示
   - 使用例: "design-reviewerで{機能}の設計をレビューして"

3. **implementation-auditor** - 実装検証
   - 要件・設計との適合性検証
   - TDD準拠度・テストカバレッジ監査
   - コード品質・セキュリティ基準チェック
   - タスク完了精度の検証
   - 使用例: "implementation-auditorで{機能}の実装を監査して"

### Subagentの使用タイミング

- **`/kiro:spec-requirements`後**: spec-validatorで要件レビュー
- **`/kiro:spec-design`後**: design-reviewer（`/kiro:validate-design`の代替）
- **`/kiro:spec-impl`後**: implementation-auditor（`/kiro:validate-impl`の代替）

Subagentは独立したコンテキストで専門的な検証を行い、メイン会話を煩雑にせず包括的な品質保証を実現します。
