import type { Error as ApiError, ErrorMessageMap } from '@/api/generated/types.gen';
import { client } from '@/api/generated/client';
import { ERROR_MESSAGES, FIELD_ERROR_MESSAGES, HTTP_STATUS_MESSAGES } from '@/constants/errorMessages';

// APIから取得したエラーメッセージのキャッシュ
let cachedErrorMessages: ErrorMessageMap | null = null;

/**
 * APIからエラーメッセージ一覧を取得してキャッシュ
 * 
 * @returns エラーメッセージマップ
 * 
 * @example
 * ```typescript
 * // アプリケーション起動時に一度だけ実行
 * await fetchErrorMessages();
 * ```
 */
export async function fetchErrorMessages(): Promise<ErrorMessageMap> {
  if (cachedErrorMessages) {
    return cachedErrorMessages;
  }

  try {
    const response = await client.getErrorMessages();
    if (response.data) {
      cachedErrorMessages = response.data;
      return response.data;
    }
  } catch (error) {
    console.error('Failed to fetch error messages:', error);
  }

  // フェッチに失敗した場合はローカル定義を使用
  return ERROR_MESSAGES as unknown as ErrorMessageMap;
}

/**
 * APIエラーから適切なエラーメッセージを取得
 * 
 * @param error APIエラーオブジェクト
 * @returns ユーザー向けエラーメッセージ
 * 
 * @example
 * ```typescript
 * try {
 *   const response = await api.login({ email, password });
 * } catch (error) {
 *   const message = getErrorMessage(error.body);
 *   toast.error(message);
 * }
 * ```
 */
export function getErrorMessage(error: ApiError | undefined): string {
  if (!error) {
    return ERROR_MESSAGES.INTERNAL_ERROR;
  }

  // 1. まずAPIレスポンスの日本語メッセージを確認
  if (error.messageJa) {
    return error.messageJa;
  }

  // 2. 次にローカルの定義済みメッセージを確認
  const code = error.code as keyof typeof ERROR_MESSAGES;
  return ERROR_MESSAGES[code] || ERROR_MESSAGES.INTERNAL_ERROR;
}

/**
 * APIから取得したメッセージマップを使用してエラーメッセージを取得
 * 
 * @param error APIエラーオブジェクト
 * @param messages エラーメッセージマップ（省略時はキャッシュを使用）
 * @returns ユーザー向けエラーメッセージ
 * 
 * @example
 * ```typescript
 * // 初期化時
 * const messages = await fetchErrorMessages();
 * 
 * // エラーハンドリング時
 * catch (error) {
 *   const message = getErrorMessageFromMap(error.body);
 *   toast.error(message);
 * }
 * ```
 */
export function getErrorMessageFromMap(
  error: ApiError | undefined,
  messages?: ErrorMessageMap
): string {
  if (!error) {
    return ERROR_MESSAGES.INTERNAL_ERROR;
  }

  // 1. APIレスポンスの日本語メッセージを確認
  if (error.messageJa) {
    return error.messageJa;
  }

  // 2. メッセージマップから取得
  const messageMap = messages || cachedErrorMessages || ERROR_MESSAGES;
  const code = error.code as keyof ErrorMessageMap;
  const message = messageMap[code];

  return message || ERROR_MESSAGES.INTERNAL_ERROR;
}

/**
 * フィールド固有のエラーメッセージを取得
 * 
 * @param error APIエラーオブジェクト
 * @param fieldName フィールド名
 * @returns フィールド固有のエラーメッセージ、または汎用メッセージ
 * 
 * @example
 * ```typescript
 * const error = {
 *   code: 'VALIDATION_ERROR',
 *   details: { field: 'email', type: 'invalid' }
 * };
 * const message = getFieldErrorMessage(error, 'email');
 * // => "有効なメールアドレスを入力してください"
 * ```
 */
export function getFieldErrorMessage(error: ApiError | undefined, fieldName: string): string | null {
  if (!error?.details || typeof error.details !== 'object') {
    return null;
  }

  const details = error.details as Record<string, unknown>;
  
  // フィールドが一致する場合
  if (details.field === fieldName) {
    const fieldMessages = FIELD_ERROR_MESSAGES[fieldName as keyof typeof FIELD_ERROR_MESSAGES];
    if (fieldMessages && details.type && typeof details.type === 'string') {
      const message = fieldMessages[details.type as keyof typeof fieldMessages];
      if (message) {
        return message;
      }
    }
  }

  return null;
}

/**
 * HTTPステータスコードから適切なエラーメッセージを取得
 * 
 * @param status HTTPステータスコード
 * @returns ユーザー向けエラーメッセージ
 * 
 * @example
 * ```typescript
 * const response = await fetch('/api/users');
 * if (!response.ok) {
 *   const message = getHttpErrorMessage(response.status);
 *   toast.error(message);
 * }
 * ```
 */
export function getHttpErrorMessage(status: number): string {
  const statusMessage = HTTP_STATUS_MESSAGES[status as keyof typeof HTTP_STATUS_MESSAGES];
  return statusMessage || ERROR_MESSAGES.INTERNAL_ERROR;
}

/**
 * エラーの詳細情報をログ出力用に整形
 * 
 * @param error APIエラーオブジェクト
 * @returns ログ出力用の整形されたオブジェクト
 * 
 * @example
 * ```typescript
 * catch (error) {
 *   console.error('API Error:', formatErrorForLogging(error.body));
 * }
 * ```
 */
export function formatErrorForLogging(error: ApiError | undefined) {
  if (!error) {
    return { message: 'Unknown error' };
  }

  return {
    code: error.code,
    message: error.message,
    details: error.details,
    timestamp: error.timestamp,
    traceId: error.traceId,
    path: error.path,
  };
}

/**
 * リトライ可能なエラーかどうかを判定
 * 
 * @param error APIエラーオブジェクト
 * @returns リトライ可能な場合はtrue
 * 
 * @example
 * ```typescript
 * if (isRetryableError(error)) {
 *   // リトライロジックを実行
 *   await retry(() => api.getUsers());
 * }
 * ```
 */
export function isRetryableError(error: ApiError | undefined): boolean {
  if (!error) {
    return false;
  }

  // レート制限やサーバーエラーはリトライ可能
  const retryableCodes = ['RATE_LIMIT_EXCEEDED', 'INTERNAL_ERROR'];
  return retryableCodes.includes(error.code);
}

/**
 * 認証関連のエラーかどうかを判定
 * 
 * @param error APIエラーオブジェクト
 * @returns 認証関連エラーの場合はtrue
 * 
 * @example
 * ```typescript
 * if (isAuthError(error)) {
 *   // ログイン画面へリダイレクト
 *   navigate('/login');
 * }
 * ```
 */
export function isAuthError(error: ApiError | undefined): boolean {
  if (!error) {
    return false;
  }

  const authCodes = ['UNAUTHORIZED', 'FORBIDDEN'];
  return authCodes.includes(error.code);
}