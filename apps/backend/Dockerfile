# ========== ビルドステージ ==========
FROM golang:1.24.5-alpine AS builder

# ビルドに必要なパッケージ
RUN apk add --no-cache git make ca-certificates

WORKDIR /build

# 依存関係のキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピーとビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION:-unknown}" \
    -o server cmd/server/main.go

# ========== 実行ステージ ==========
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# バイナリのコピー
COPY --from=builder /build/server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ヘルスチェック用
EXPOSE 8080

# 非rootユーザーで実行
USER nonroot:nonroot

ENTRYPOINT ["./server"]
