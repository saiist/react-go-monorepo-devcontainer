.PHONY: help dev dev-frontend dev-backend build test clean docker-up docker-down generate db-migrate health-check

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo "  make dev          - 開発サーバーを起動"
	@echo "  make build        - プロダクションビルド"
	@echo "  make test         - テストを実行"
	@echo "  make generate     - API型を生成"
	@echo "  make db-migrate   - DBマイグレーションを実行"
	@echo "  make docker-up    - Docker環境を起動"
	@echo "  make docker-down  - Docker環境を停止"
	@echo "  make clean        - ビルド成果物を削除"

# 開発
dev: docker-up
	@echo "🚀 開発サーバーを起動中..."
	@make -j 2 dev-frontend dev-backend

dev-frontend:
	pnpm --filter frontend dev

dev-backend:
	cd apps/backend && air

# ビルド
build:
	pnpm build

# テスト
test:
	pnpm test

# 型生成
generate:
	pnpm generate

# DBマイグレーションを実行
db-migrate:
	pnpm db:migrate

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

# health check API テスト
dev-health-check:
	 curl -X GET http://localhost:8080/api/v1/health
