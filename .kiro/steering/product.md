# Product Overview

TSUNAGU Linkのバックエンドサービス - ユーザー認証とユーザー管理を提供するRESTful APIプラットフォーム

## Core Capabilities

- **ユーザー認証**: JWTベースの認証システム（将来的にZitadel認証基盤へ移行予定）
- **ユーザー管理**: ユーザーのCRUD操作とライフサイクル管理
- **API駆動設計**: OpenAPI 3.0仕様に基づく標準化されたREST API
- **データベース統合**: PostgreSQLを使用したリレーショナルデータ管理
- **開発者体験**: ホットリロード、マイグレーション、Dockerサポートによる高速な開発サイクル

## Target Use Cases

- **フロントエンドアプリケーション連携**: React/Vue/Next.jsなどのSPAからのAPI呼び出し
- **モバイルアプリケーション**: iOSおよびAndroidアプリのバックエンド
- **サードパーティ統合**: OpenAPI仕様を通じた外部サービス連携

## Value Proposition

- **型安全性**: Go言語の静的型付けによる堅牢性
- **パフォーマンス**: Goランタイムとコンパイル言語の高速性
- **保守性**: レイヤードアーキテクチャによる明確な責任分離
- **スケーラビリティ**: PostgreSQLとDockerによる水平スケーリング対応
- **開発生産性**: Makefileタスクランナー、Air（ホットリロード）、マイグレーションツールによる効率的な開発体験

## Future Roadmap

- Zitadel認証基盤への統合（エンタープライズ認証対応）
- より多くのAPIエンドポイントの拡充
- 単体テスト・統合テストの充実
- CI/CDパイプラインの構築

---
_Focus on patterns and purpose, not exhaustive feature lists_
