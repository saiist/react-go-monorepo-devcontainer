# ========== ビルドステージ ==========
FROM node:24.4-alpine AS builder

# pnpmのインストール
RUN corepack enable && corepack prepare pnpm@9.5.0 --activate

WORKDIR /build

# 依存関係のキャッシュ
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# ソースコードのコピーとビルド
COPY . .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN pnpm build

# ========== 実行ステージ ==========
FROM nginx:1.27.0-alpine

# Nginxの設定
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /build/dist /usr/share/nginx/html

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost || exit 1

EXPOSE 80
