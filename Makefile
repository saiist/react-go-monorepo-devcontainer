.PHONY: help dev build test clean docker-up docker-down generate

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo "  make dev          - 開発サーバーを起動"
	@echo "  make build        - プロダクションビルド"
	@echo "  make test         - テストを実行"
	@echo "  make generate     - API型を生成"
	@echo "  make docker-up    - Docker環境を起動"
	@echo "  make docker-down  - Docker環境を停止"
	@echo "  make clean        - ビルド成果物を削除"

# 開発
dev: docker-up
	pnpm dev

# ビルド
build:
	pnpm build

# テスト
test:
	pnpm test

# 型生成
generate:
	pnpm generate

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# クリーン
clean:
	pnpm clean
	rm -rf .turbo
	rm -rf apps/backend/tmp
	rm -rf apps/backend/coverage.out
