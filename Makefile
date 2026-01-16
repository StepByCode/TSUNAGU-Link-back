.PHONY: help build run dev test test-verbose test-watch clean docker-up docker-down docker-build migrate openapi-gen install-tools

APP_NAME=server
BIN_DIR=bin
CMD_DIR=cmd/server

help: ## ヘルプメッセージを表示
	@echo '使い方: make [ターゲット]'
	@echo ''
	@echo '利用可能なターゲット:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## 開発ツールをインストール
	@echo "開発ツールをインストール中..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install gotest.tools/gotestsum@latest
	@which atlas > /dev/null || (echo "Atlasをインストール中..." && curl -sSf https://atlasgo.sh | sh)

build: ## アプリケーションをビルド
	@echo "アプリケーションをビルド中..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)

run: build ## アプリケーションをビルドして実行
	@echo "アプリケーションを起動中..."
	./$(BIN_DIR)/$(APP_NAME)

dev: ## 開発モードで実行（ホットリロード有効）
	@which air > /dev/null || (echo "air をインストール中..." && go install github.com/air-verse/air@latest)
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/air

test: ## テストを実行（gotestsum使用）
	@which gotestsum > /dev/null || (echo "gotestsum をインストール中..." && go install gotest.tools/gotestsum@latest)
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/gotestsum --format testname -- -race -coverprofile=coverage.out ./...
	@echo "カバレッジレポートを生成中..."
	@go tool cover -html=coverage.out -o coverage.html

test-verbose: ## テストを詳細出力で実行
	@which gotestsum > /dev/null || (echo "gotestsum をインストール中..." && go install gotest.tools/gotestsum@latest)
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/gotestsum --format standard-verbose -- -race -coverprofile=coverage.out ./...

test-watch: ## テストをウォッチモードで実行
	@which gotestsum > /dev/null || (echo "gotestsum をインストール中..." && go install gotest.tools/gotestsum@latest)
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/gotestsum --watch -- -race ./...

clean: ## ビルド成果物を削除
	@echo "クリーンアップ中..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

docker-up: ## Dockerコンテナを起動
	docker-compose up -d

docker-down: ## Dockerコンテナを停止
	docker-compose down

docker-build: ## Dockerイメージをビルド
	docker-compose build

docker-logs: ## Dockerログを表示
	docker-compose logs -f

migrate-up: ## データベースマイグレーションを実行（up）
	atlas migrate apply --env dev
migrate-down: ## データベースマイグレーションを実行（down）
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/migrate -path db/migrations -database "postgres://tsunagu:tsunagu_password@localhost:5432/tsunagu_db?sslmode=disable" down

migrate-create: ## 新しいマイグレーションファイルを作成（使い方: make migrate-create name=create_users_table）
	@if [ -z "$(name)" ]; then echo "使い方: make migrate-create name=マイグレーション名"; exit 1; fi
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/migrate create -ext sql -dir db/migrations -seq $(name)

migrate-generate: ## スキーマ定義からマイグレーションを自動生成（使い方: make migrate-generate [name=add_new_table]）
	@which atlas > /dev/null || (echo "Atlasがインストールされていません。make install-tools を実行してください。" && exit 1)
	@if [ -z "$(name)" ]; then \
		echo "マイグレーション名を入力してください（例: add_posts_table）:"; \
		read migration_name; \
		if [ -z "$$migration_name" ]; then \
			migration_name="migration_$$(date +%Y%m%d_%H%M%S)"; \
			echo "マイグレーション名が未指定のため、デフォルト名を使用: $$migration_name"; \
		fi; \
		echo "スキーマ定義からマイグレーションを生成中..."; \
		atlas migrate diff $$migration_name \
			--dir "file://db/migrations?format=golang-migrate" \
			--to "file://db/schema.hcl" \
			--dev-url "docker://postgres/16/dev"; \
	else \
		echo "スキーマ定義からマイグレーションを生成中..."; \
		atlas migrate diff $(name) \
			--dir "file://db/migrations?format=golang-migrate" \
			--to "file://db/schema.hcl" \
			--dev-url "docker://postgres/16/dev"; \
	fi

migrate-generate-auto: ## スキーマ差分から自動的にマイグレーションを生成（タイムスタンプ名）
	@which atlas > /dev/null || (echo "Atlasがインストールされていません。make install-tools を実行してください。" && exit 1)
	@migration_name="migration_$$(date +%Y%m%d_%H%M%S)"; \
	echo "自動生成モード: $$migration_name"; \
	echo "スキーマ定義からマイグレーションを生成中..."; \
	atlas migrate diff $$migration_name --env schema && \
	echo "" && \
	echo "⚠️  マイグレーション名を変更することを推奨します:" && \
	echo "   生成されたファイルをリネームしてください（例: 000XXX_$$migration_name.up.sql → 000XXX_add_posts_table.up.sql）"

migrate-lint: ## マイグレーションファイルを検証
	@which atlas > /dev/null || (echo "Atlasがインストールされていません。make install-tools を実行してください。" && exit 1)
	@echo "マイグレーションファイルを検証中..."
	atlas migrate lint --env dev

migrate-status: ## マイグレーションの状態を確認
	@which atlas > /dev/null || (echo "Atlasがインストールされていません。make install-tools を実行してください。" && exit 1)
	@echo "マイグレーション状態を確認中..."
	atlas migrate status --env dev

schema-inspect: ## 現在のデータベーススキーマを表示
	@which atlas > /dev/null || (echo "Atlasがインストールされていません。make install-tools を実行してください。" && exit 1)
	@echo "データベーススキーマを取得中..."
	atlas schema inspect --env dev

openapi-gen: ## OpenAPI仕様からコードを生成
	@echo "OpenAPI仕様からコードを生成中..."
	@mkdir -p internal/api
	@GOPATH=$$(go env GOPATH); $$GOPATH/bin/oapi-codegen -package api -generate types,server,spec api/openapi.yaml > internal/api/api_generated.go

deps: ## 依存関係をダウンロード
	go mod download
	go mod tidy

fmt: ## コードをフォーマット
	go fmt ./...
	gofmt -s -w .

lint: ## リンターを実行（golangci-lint が必要）
	@which golangci-lint > /dev/null || (echo "golangci-lint をインストールしてください: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

setup: install-tools deps docker-up migrate-up ## プロジェクトをセットアップ（ツールインストール、依存関係、Docker起動、マイグレーション実行）
	@echo "プロジェクトのセットアップが完了しました！"

all: clean build test ## クリーン、ビルド、テストを実行
