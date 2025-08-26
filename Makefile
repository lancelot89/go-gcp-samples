.PHONY: help init fmt lint test build clean tidy

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: ## 初期設定とワークスペース同期
	go work sync
	$(MAKE) tidy

fmt: ## コードフォーマット
	cd v1-cloud-run && go fmt ./...
	cd v2-firestore && go fmt ./...

lint: ## Lintチェック
	@if ! which golangci-lint > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	cd v1-cloud-run && golangci-lint run ./...
	cd v2-firestore && golangci-lint run ./...

test: ## 全モジュールのテスト実行
	cd v1-cloud-run && go test ./... -count=1 -race -v -cover
	cd v2-firestore && go test ./... -count=1 -race -v -cover

test-short: ## 短縮テスト（エミュレータ依存なし）
	cd v1-cloud-run && go test ./... -short -count=1 -race -v
	cd v2-firestore && go test ./... -short -count=1 -race -v

build: ## 全モジュールのビルド
	@echo "Building v1-cloud-run..."
	cd v1-cloud-run && go build -o bin/app ./...
	@echo "Building v2-firestore..."
	cd v2-firestore && go build -o bin/app ./...

clean: ## ビルド成果物をクリーン
	rm -rf v1-cloud-run/bin
	rm -rf v2-firestore/bin
	go clean -cache -testcache

tidy: ## 依存関係の整理
	cd v1-cloud-run && go mod tidy
	cd v2-firestore && go mod tidy

check: fmt lint test ## フォーマット、Lint、テストをまとめて実行

ci: ## CI用タスク（キャッシュなし）
	$(MAKE) fmt
	$(MAKE) lint
	$(MAKE) test