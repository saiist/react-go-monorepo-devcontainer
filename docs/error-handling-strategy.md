# エラーハンドリング戦略

## 概要

このプロジェクトでは、APIとフロントエンドで責務を分離したエラーハンドリング戦略を採用しています。

## 設計原則

### 1. APIレスポンス（OpenAPI仕様）
- **エラーコード**: 技術的なエラー種別を示す（例: `VALIDATION_ERROR`, `CONFLICT`）
- **メッセージ**: 開発者向けの技術的メッセージ（英語）
- **詳細情報**: エラーの詳細データ（フィールド名、値など）

### 2. フロントエンド
- **表示メッセージ**: エンドユーザー向けの日本語メッセージ
- **エラーコードベースの管理**: `constants/errorMessages.ts`で一元管理

## エラーレスポンスの構造

```json
{
  "code": "CONFLICT",
  "message": "User with email user@example.com already exists",
  "details": {
    "field": "email",
    "value": "user@example.com"
  },
  "timestamp": "2025-01-20T10:30:00Z",
  "traceId": "b3d4f9e8-1234-5678-9abc-def012345678",
  "path": "/api/v1/auth/register"
}
```

## メリット

1. **責務の分離**
   - API: エラーの技術的な情報を提供
   - UI: ユーザーフレンドリーなメッセージを表示

2. **保守性**
   - メッセージの変更がAPIの変更を必要としない
   - 表示文言の A/B テストが容易

3. **多言語対応**
   - フロントエンドで言語切り替えが可能
   - APIは言語に依存しない

4. **デバッグ効率**
   - `traceId`によるログ追跡
   - `timestamp`による時系列分析
   - `path`によるエンドポイント特定

## 実装ガイドライン

### バックエンド
```go
// エラーレスポンスの作成
func NewErrorResponse(code api.ErrorCode, message string, path string) api.Error {
    timestamp := time.Now()
    traceId := uuid.New().String()
    
    return api.Error{
        Code:      code,
        Message:   message,
        Timestamp: &timestamp,
        TraceId:   &traceId,
        Path:      &path,
    }
}
```

### フロントエンド
```typescript
// エラーハンドリング
import { getErrorMessage } from '@/utils/errorHandler';

try {
  await api.register(userData);
} catch (error) {
  // エラーコードから日本語メッセージを取得
  const message = getErrorMessage(error.body);
  toast.error(message); // "既に登録されています"
}
```

## エラーコード一覧

| コード | 説明 | 使用場面 |
|--------|------|----------|
| `VALIDATION_ERROR` | 入力値の検証エラー | フォーム入力が不正 |
| `UNAUTHORIZED` | 認証が必要 | ログインが必要な操作 |
| `FORBIDDEN` | アクセス権限なし | 権限不足 |
| `NOT_FOUND` | リソースが存在しない | 存在しないデータへのアクセス |
| `CONFLICT` | リソースの競合 | 重複登録など |
| `RATE_LIMIT_EXCEEDED` | レート制限超過 | API呼び出し過多 |
| `INTERNAL_ERROR` | サーバー内部エラー | 予期しないエラー |

## OpenAPI仕様の活用

OpenAPI仕様には以下の情報を含めています：

1. **エラーコードの説明** (`x-enum-descriptions`)
2. **レスポンス例** (`examples`)
3. **エラー構造の定義** (`Error` スキーマ)

これにより、API利用者は：
- エラーレスポンスの形式を理解できる
- 各エラーコードの意味を把握できる
- 実際のレスポンス例を参照できる

## まとめ

この戦略により、APIの汎用性を保ちながら、フロントエンドで柔軟なエラーメッセージ管理が可能になります。開発効率とユーザー体験の両方を向上させる設計となっています。