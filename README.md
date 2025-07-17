# My App

Go + React モノレポプロジェクト

## 技術スタック

- **Backend**: Go 1.22.5, Chi Router, GORM, PostgreSQL
- **Frontend**: React 18.3, TypeScript, Vite, TanStack Query
- **共通**: OpenAPI 3.0, Docker, pnpm workspaces

## セットアップ

### 前提条件

- Docker Desktop
- VSCode
- Git

### 開発環境の起動

1. リポジトリをクローン
```bash
git clone <repository-url>
cd my-app
```

2. VSCodeで開く
```bash
code .
```

3. DevContainerで開く
- Command Palette (Cmd/Ctrl+Shift+P)
- "Dev Containers: Reopen in Container" を選択

4. 開発サーバーを起動
```bash
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
pnpm dev

# ビルド
pnpm build

# テスト
pnpm test

# 型生成
pnpm generate

# フォーマット
pnpm format

# リント
pnpm lint
```

### プロジェクト構造

```
.
├── apps/
│   ├── frontend/     # Reactアプリケーション
│   └── backend/      # Go APIサーバー
├── api/
│   └── openapi.yaml  # API定義
└── packages/         # 共有パッケージ（将来用）
```

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
docker build -t my-app-backend:latest -f apps/backend/Dockerfile apps/backend

# Frontend
docker build -t my-app-frontend:latest -f apps/frontend/Dockerfile apps/frontend
```

### 環境変数

`.env.example` を参考に、各環境用の `.env` ファイルを作成してください。

## ライセンス

MIT
