# React-Go Monorepo

Go + React モノレポプロジェクト - Contract-First API開発

## 技術スタック

- **Backend**: Go 1.23.5, Chi Router, GORM, PostgreSQL, Redis
- **Frontend**: React 19.1, TypeScript, Vite, TanStack Query, Zustand
- **共通**: OpenAPI 3.0, Turbo, Docker, pnpm workspaces
- **開発ツール**: Air (Go hot reload), DevContainer, GitHub Actions

## セットアップ

### 前提条件

- Docker Desktop
- VSCode
- Git

### 開発環境の起動

#### DevContainer（推奨）

1. リポジトリをクローン
```bash
git clone <repository-url>
cd react-go-monorepo
```

2. VSCodeで開く
```bash
code .
```

3. DevContainerで開く
- Command Palette (Cmd/Ctrl+Shift+P)
- "Dev Containers: Reopen in Container" を選択
- 初回起動時は自動的にセットアップが実行されます

4. 開発サーバーを起動
```bash
pnpm dev
```

#### ローカル環境

必要なツール:
- Node.js 22+
- Go 1.23+
- PostgreSQL 17
- Redis 8
- pnpm 9+

```bash
# 依存関係のインストール
pnpm install

# Docker サービスの起動
pnpm docker:up

# 型の生成
pnpm generate

# データベースマイグレーション
pnpm db:migrate

# 開発サーバーの起動
pnpm dev
```

アクセス:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/api/v1/docs

## 開発

### コマンド

```bash
# 開発サーバー
pnpm dev              # フロントエンドとバックエンドを同時起動
pnpm backend:dev      # バックエンドのみ起動
pnpm --filter frontend dev  # フロントエンドのみ起動

# ビルド
pnpm build           # 全てのアプリをビルド
pnpm backend:build   # バックエンドのみビルド

# テスト
pnpm test            # 全てのテストを実行
pnpm backend:test    # バックエンドのテスト

# 型生成（重要）
pnpm generate        # OpenAPIから型を自動生成

# データベース
pnpm db:migrate      # マイグレーション実行
pnpm db:migrate:down # 最後のマイグレーションをロールバック
pnpm db:migrate:create NAME=<name>  # 新しいマイグレーション作成

# Docker
pnpm docker:up       # PostgreSQLとRedisを起動
pnpm docker:down     # サービスを停止

# コード品質
pnpm format          # Prettierでフォーマット
pnpm lint            # ESLintとgolangci-lintを実行
pnpm type-check      # TypeScriptの型チェック

# その他
pnpm health-check    # APIヘルスチェック
pnpm clean          # ビルド成果物を削除
```

### プロジェクト構造

```
.
├── apps/
│   ├── frontend/           # Reactアプリケーション
│   │   ├── src/
│   │   │   ├── api/       # 自動生成されたAPIクライアント
│   │   │   ├── components/
│   │   │   ├── hooks/
│   │   │   ├── pages/
│   │   │   └── stores/    # Zustand stores
│   │   └── vite.config.ts
│   └── backend/            # Go APIサーバー
│       ├── cmd/server/    # エントリポイント
│       ├── internal/      # 内部パッケージ
│       │   ├── api/       # 自動生成されたAPI型
│       │   ├── config/    # 設定管理
│       │   ├── db/        # データベース接続
│       │   └── handler/   # HTTPハンドラー
│       └── migrations/    # DBマイグレーション
├── api/
│   └── openapi.yaml       # API定義（契約）
├── packages/              # 共有パッケージ（将来用）
└── .devcontainer/         # DevContainer設定
```

### Contract-First API開発

このプロジェクトはOpenAPI仕様を中心とした開発フローを採用しています：

1. **API設計**: `/api/openapi.yaml`を編集
2. **型生成**: `pnpm generate`を実行
3. **実装**: 生成された型を使用して実装

**重要**: APIを変更する際は必ず以下の順序で作業してください：
1. `openapi.yaml`を更新
2. `pnpm generate`を実行
3. フロントエンド・バックエンドの実装を更新

### Git コミット規約

Conventional Commitsに従います:

- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: コードスタイル
- `refactor`: リファクタリング
- `test`: テスト
- `chore`: その他

例:
```bash
git commit -m "feat: ユーザー認証機能を追加"
git commit -m "fix: ログイン時のエラーハンドリングを修正"
```

## デプロイ

### Docker イメージのビルド

```bash
# Backend
docker build -t react-go-backend:latest -f apps/backend/Dockerfile apps/backend

# Frontend  
docker build -t react-go-frontend:latest -f apps/frontend/Dockerfile apps/frontend
```

### 環境変数

`.env.example` を参考に、各環境用の `.env` ファイルを作成してください。

主な環境変数:
- `DATABASE_URL`: PostgreSQL接続文字列（sslmode=disableを含む）
- `JWT_SECRET`: JWT署名用のシークレットキー
- `VITE_API_URL`: フロントエンドからのAPI接続先
- `FRONTEND_URL`: CORS設定用のフロントエンドURL

## トラブルシューティング

### pnpm devでバックエンドが起動しない
```bash
# Airがインストールされているか確認
which air

# インストールされていない場合
go install github.com/air-verse/air@latest
```

### データベース接続エラー
```bash
# Dockerサービスが起動しているか確認
docker ps

# 起動していない場合
pnpm docker:up
```

### 型生成エラー
```bash
# OpenAPI仕様のバリデーション
npx @redocly/cli lint api/openapi.yaml
```

## ライセンス

MIT
