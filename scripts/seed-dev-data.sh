#!/bin/bash

set -e

# 色付きログ出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 環境変数の読み込み
ENV_FILE=".env"
if [ -f "$ENV_FILE" ]; then
    log_info "環境変数を読み込んでいます: $ENV_FILE"
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    log_error ".env ファイルが見つかりません: $(pwd)/$ENV_FILE"
    exit 1
fi

# データベース接続情報の解析
if [ -z "$DATABASE_URL" ]; then
    log_error "DATABASE_URL が設定されていません"
    log_error "環境変数ファイル: $(pwd)/$ENV_FILE"
    exit 1
fi

# DATABASE_URLの値を表示（パスワードは一部マスク）
MASKED_URL=$(echo "$DATABASE_URL" | sed -E 's/(:\/\/[^:]+:)[^@]+(@)/\1****\2/')
log_info "DATABASE_URL: $MASKED_URL"

# PostgreSQL接続情報を解析
# 形式: postgres://user:password@host:port/dbname?sslmode=disable
# または: postgresql://user:password@host:port/dbname?sslmode=disable
if [[ $DATABASE_URL =~ postgres(ql)?://([^:]+):([^@]+)@([^:]+):([0-9]+)/([^?]+) ]]; then
    DB_USER="${BASH_REMATCH[2]}"
    DB_PASSWORD="${BASH_REMATCH[3]}"
    DB_HOST="${BASH_REMATCH[4]}"
    DB_PORT="${BASH_REMATCH[5]}"
    DB_NAME="${BASH_REMATCH[6]}"
else
    log_error "DATABASE_URL の形式が不正です"
    log_error "期待される形式: postgres[ql]://user:password@host:port/dbname"
    log_error "実際の値: $DATABASE_URL"
    log_error "環境変数ファイル: $(pwd)/$ENV_FILE"
    exit 1
fi

# psqlコマンドのエクスポート
export PGPASSWORD=$DB_PASSWORD

log_info "開発用シードデータを投入します..."
log_info "データベース: $DB_NAME@$DB_HOST:$DB_PORT"

# 既存データの確認
USER_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null || echo "0")

if [ "$USER_COUNT" -gt "0" ]; then
    log_warning "既にユーザーデータが存在します (${USER_COUNT}件)"
    read -p "既存データを削除して新しいデータを投入しますか？ (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "処理を中止しました"
        exit 0
    fi
    
    log_info "既存データを削除しています..."
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
TRUNCATE TABLE users CASCADE;
EOF
fi

# 開発用シードデータの投入
log_info "シードデータを投入しています..."

psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<'EOF'
-- テストユーザーの作成
-- パスワードは全て 'password123' をbcryptでハッシュ化したもの

INSERT INTO users (email, username, password_hash, full_name, avatar_url, is_active, is_verified) VALUES
    -- 管理者ユーザー
    ('admin@example.com', 'admin', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', '管理者', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', true, true),
    
    -- 一般ユーザー（アクティブ・検証済み）
    ('john.doe@example.com', 'johndoe', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'John Doe', 'https://api.dicebear.com/7.x/avataaars/svg?seed=john', true, true),
    ('jane.smith@example.com', 'janesmith', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Jane Smith', 'https://api.dicebear.com/7.x/avataaars/svg?seed=jane', true, true),
    ('alice.johnson@example.com', 'alicej', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Alice Johnson', 'https://api.dicebear.com/7.x/avataaars/svg?seed=alice', true, true),
    
    -- 未検証ユーザー
    ('bob.wilson@example.com', 'bobwilson', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Bob Wilson', 'https://api.dicebear.com/7.x/avataaars/svg?seed=bob', true, false),
    ('charlie.brown@example.com', 'charlieb', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Charlie Brown', 'https://api.dicebear.com/7.x/avataaars/svg?seed=charlie', true, false),
    
    -- 非アクティブユーザー
    ('david.lee@example.com', 'davidlee', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'David Lee', 'https://api.dicebear.com/7.x/avataaars/svg?seed=david', false, true),
    ('emma.davis@example.com', 'emmad', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Emma Davis', 'https://api.dicebear.com/7.x/avataaars/svg?seed=emma', false, true),
    
    -- 日本語名のユーザー
    ('taro.yamada@example.jp', 'taroyamada', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', '山田 太郎', 'https://api.dicebear.com/7.x/avataaars/svg?seed=taro', true, true),
    ('hanako.tanaka@example.jp', 'hanakot', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', '田中 花子', 'https://api.dicebear.com/7.x/avataaars/svg?seed=hanako', true, true),
    
    -- テスト用ユーザー（様々なケース）
    ('test.user1@example.com', 'testuser1', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Test User 1', NULL, true, true),
    ('test.user2@example.com', 'testuser2', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', 'Test User 2', 'https://api.dicebear.com/7.x/avataaars/svg?seed=test2', true, false),
    ('test.user3@example.com', 'testuser3', '$2a$10$GhgC4HMIvptPMdbXA33/PeZ/JfHDg53B4NniltrhRsnXowCWJ//sS', NULL, NULL, true, true);

-- 作成日時を分散させる（過去30日間）
UPDATE users SET created_at = CURRENT_TIMESTAMP - (random() * INTERVAL '30 days') WHERE email LIKE 'test%';
UPDATE users SET created_at = CURRENT_TIMESTAMP - INTERVAL '7 days' WHERE email = 'admin@example.com';
UPDATE users SET created_at = CURRENT_TIMESTAMP - INTERVAL '14 days' WHERE email LIKE 'john%' OR email LIKE 'jane%';
UPDATE users SET created_at = CURRENT_TIMESTAMP - INTERVAL '21 days' WHERE email LIKE 'alice%' OR email LIKE 'bob%';

-- 更新日時も調整
UPDATE users SET updated_at = created_at + (random() * INTERVAL '7 days') WHERE is_active = true;
UPDATE users SET updated_at = created_at WHERE is_active = false;

EOF

# 結果の確認
FINAL_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;")
log_success "シードデータの投入が完了しました！"
log_info "作成されたユーザー数: ${FINAL_COUNT}"

# サンプルユーザー情報の表示
log_info ""
log_info "=== サンプルユーザー情報 ==="
log_info "すべてのユーザーのパスワード: password123"
log_info ""
log_info "管理者:"
log_info "  - Email: admin@example.com"
log_info "  - Username: admin"
log_info ""
log_info "一般ユーザー（検証済み）:"
log_info "  - Email: john.doe@example.com"
log_info "  - Username: johndoe"
log_info ""
log_info "  - Email: jane.smith@example.com"
log_info "  - Username: janesmith"
log_info ""
log_info "テストユーザー:"
log_info "  - Email: test.user1@example.com"
log_info "  - Username: testuser1"

# パスワードハッシュに関する情報
log_info ""
log_info "=== パスワードハッシュ情報 ==="
log_info "すべてのユーザーのパスワード 'password123' は bcrypt でハッシュ化されています。"
log_info ""
log_info "別のパスワードを使用する場合:"
log_info "go run apps/backend/cmd/hash-password/main.go <新しいパスワード>"