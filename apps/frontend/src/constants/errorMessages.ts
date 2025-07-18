/**
 * APIエラーコードに対応する日本語メッセージ
 * 
 * @example
 * ```typescript
 * import { ERROR_MESSAGES } from '@/constants/errorMessages';
 * import type { ErrorCode } from '@/api/generated/types.gen';
 * 
 * const handleError = (code: ErrorCode) => {
 *   const message = ERROR_MESSAGES[code] || ERROR_MESSAGES.INTERNAL_ERROR;
 *   console.error(message);
 * };
 * ```
 */
export const ERROR_MESSAGES = {
  VALIDATION_ERROR: '入力内容に誤りがあります',
  UNAUTHORIZED: 'ログインが必要です',
  FORBIDDEN: 'アクセス権限がありません',
  NOT_FOUND: 'データが見つかりません',
  CONFLICT: '既に登録されています',
  RATE_LIMIT_EXCEEDED: 'アクセス制限中です。しばらくお待ちください',
  INTERNAL_ERROR: 'システムエラーが発生しました',
} as const;

/**
 * フィールド固有のエラーメッセージ
 * 
 * @example
 * ```typescript
 * const error = {
 *   code: 'VALIDATION_ERROR',
 *   details: { field: 'email', type: 'invalid' }
 * };
 * 
 * const message = FIELD_ERROR_MESSAGES.email.invalid;
 * ```
 */
export const FIELD_ERROR_MESSAGES = {
  email: {
    required: 'メールアドレスは必須です',
    invalid: '有効なメールアドレスを入力してください',
    exists: 'このメールアドレスは既に使用されています',
  },
  username: {
    required: 'ユーザー名は必須です',
    minLength: 'ユーザー名は3文字以上で入力してください',
    maxLength: 'ユーザー名は100文字以内で入力してください',
    exists: 'このユーザー名は既に使用されています',
  },
  password: {
    required: 'パスワードは必須です',
    minLength: 'パスワードは8文字以上で入力してください',
    maxLength: 'パスワードは128文字以内で入力してください',
    pattern: 'パスワードは大文字、小文字、数字、特殊文字を少なくとも1つずつ含む必要があります',
    mismatch: 'パスワードが一致しません',
  },
  fullName: {
    maxLength: '氏名は255文字以内で入力してください',
  },
  avatarUrl: {
    invalid: '有効なURLを入力してください',
    maxLength: 'URLは500文字以内で入力してください',
  },
} as const;

/**
 * HTTPステータスコードに対応するデフォルトメッセージ
 */
export const HTTP_STATUS_MESSAGES = {
  400: 'リクエストが正しくありません',
  401: '認証が必要です',
  403: 'アクセスが拒否されました',
  404: 'ページが見つかりません',
  409: 'データの競合が発生しました',
  429: 'リクエストが多すぎます',
  500: 'サーバーエラーが発生しました',
  502: 'サーバーが応答していません',
  503: 'サービスが一時的に利用できません',
} as const;