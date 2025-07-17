#!/bin/bash
set -euo pipefail

# カラー出力の定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Git設定
log_info "Git設定中..."
git config --global user.name "${GIT_USER_NAME:-Developer}"
git config --global user.email "${GIT_USER_EMAIL:-dev@example.com}"
git config --global core.editor "code --wait"
git config --global pull.rebase true
git config --global init.defaultBranch main
git config --global core.autocrlf input

# SSH設定（必要に応じて）
if [ -d "/home/vscode/.ssh" ]; then
    chmod 700 /home/vscode/.ssh
    if [ -f "/home/vscode/.ssh/id_rsa" ]; then
        chmod 600 /home/vscode/.ssh/id_rsa
    fi
fi

# Node.jsグローバルパッケージのインストール
log_info "🔧 Node.jsグローバルパッケージをインストール中..."
npm install -g turbo npm-check-updates || {
    log_error "Node.jsグローバルパッケージのインストールに失敗しました"
    exit 1
}


# Go開発ツールのインストール
echo "🔧 Go開発ツールをインストール中..."
go install github.com/air-verse/air
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
go install github.com/golangci/golangci-lint/cmd/golangci-lint
go install github.com/golang-migrate/migrate/v4/cmd/migrate
go install github.com/swaggo/swag/cmd/swag

# Go依存関係のダウンロード
if [ -f "apps/backend/go.mod" ]; then
    log_info "🔧 Go依存関係をダウンロード中..."
    cd apps/backend
    go mod tidy
    cd ../..
fi

# pnpm-storeディレクトリの権限設定
if [ -d "/workspace/.pnpm-store" ] && [ "$(stat -c %U /workspace/.pnpm-store)" = "root" ]; then
    log_info "pnpm-storeディレクトリの権限を修正中..."
    sudo chown -R vscode:vscode /workspace/.pnpm-store
fi

# 依存関係のインストール
echo "📦 依存関係をインストール中..."
pnpm install --frozen-lockfile

# 型の生成
echo "🔧 APIの型を生成中..."
if ! pnpm run generate; then
    log_warning "API型の生成でエラーが発生しました。手動で 'pnpm run generate' を実行してください。"
fi

# 環境変数ファイルの作成
if [ ! -f .env ]; then
    cp .env.example .env
    log_info "✅ .env ファイルを作成しました"
else
    log_warning ".env ファイルは既に存在します"
fi

# データベースのマイグレーション（PostgreSQLが起動している場合）
if command -v migrate &> /dev/null && [ -d "apps/backend/migrations" ]; then
    log_info "🗄️ データベースマイグレーションを確認中..."

    # PostgreSQLの起動を待機（最大30秒）
    for i in {1..30}; do
        if nc -z postgres 5432; then
            break
        fi
        log_info "PostgreSQLの起動を待機中... ($i/30)"
        sleep 1
    done

    if nc -z postgres 5432; then
        cd apps/backend
        migrate -path ./migrations -database "${DATABASE_URL:-postgres://postgres:postgres@postgres:5432/myapp_dev?sslmode=disable}" up
        cd ../..
        log_info "✅ マイグレーション完了"
    else
        log_warning "PostgreSQLが起動していないため、マイグレーションをスキップしました"
    fi
fi



# 開発用の初期データ投入（オプション）
if [ -f "scripts/seed-dev-data.sh" ]; then
    log_info "🌱 開発用データを投入中..."
    bash scripts/seed-dev-data.sh
fi


# ヘルスチェック
log_info "🏥 環境のヘルスチェック中..."
node --version
go version
pnpm --version

# サービスの接続確認
if nc -z postgres 5432; then
    log_info "✅ PostgreSQLに接続できます"
else
    log_warning "PostgreSQLに接続できません"
fi

if nc -z redis 6379; then
    log_info "✅ Redisに接続できます"
else
    log_warning "Redisに接続できません"
fi

log_info "✅ セットアップ完了！"
echo ""
echo "開発を始めるには:"
echo "  pnpm dev        - 開発サーバーを起動"
echo ""
echo "利用可能なコマンド:"
echo "  pnpm dev        - 開発サーバーを起動"
echo "  pnpm build      - プロダクションビルド"
echo "  pnpm test       - テストを実行"
echo "  pnpm lint       - リントを実行"
echo "  pnpm format     - コードをフォーマット"
echo "  pnpm generate   - API型を生成"
echo "  pnpm db:migrate - DBマイグレーションを実行"
echo ""
echo "VSCode推奨:"
echo "  Cmd/Ctrl+Shift+P → 'Developer: Reload Window' で拡張機能を有効化"
