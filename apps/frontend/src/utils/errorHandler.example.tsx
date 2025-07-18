/**
 * エラーハンドリングの実装例
 * 
 * このファイルは実際のコンポーネントではなく、
 * エラーハンドリングの実装パターンを示すためのサンプルです。
 */

import { useState } from 'react';
import type { Error as ApiError } from '@/api/generated/types.gen';
import { 
  getErrorMessage, 
  getFieldErrorMessage, 
  formatErrorForLogging,
  isAuthError,
  isRetryableError 
} from './errorHandler';

/**
 * ログインフォームでのエラーハンドリング例
 */
export function LoginFormExample() {
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleLogin = async (email: string, password: string) => {
    try {
      // APIコール（実際の実装では生成されたクライアントを使用）
      // const response = await api.login({ body: { email, password } });
    } catch (error: any) {
      const apiError: ApiError = error.body;
      
      // エラーログの出力
      console.error('Login failed:', formatErrorForLogging(apiError));
      
      // 認証エラーの場合の処理
      if (isAuthError(apiError)) {
        // トースト通知などで表示
        const message = getErrorMessage(apiError);
        // toast.error(message);
      }
      
      // バリデーションエラーの場合、フィールド毎のエラーを表示
      if (apiError.code === 'VALIDATION_ERROR') {
        const emailError = getFieldErrorMessage(apiError, 'email');
        const passwordError = getFieldErrorMessage(apiError, 'password');
        
        setErrors({
          email: emailError || '',
          password: passwordError || ''
        });
      }
    }
  };

  return (
    <form>
      <input type="email" />
      {errors.email && <span className="error">{errors.email}</span>}
      
      <input type="password" />
      {errors.password && <span className="error">{errors.password}</span>}
    </form>
  );
}

/**
 * データ取得でのエラーハンドリング例（リトライ機能付き）
 */
export function DataFetchExample() {
  const [retryCount, setRetryCount] = useState(0);
  const MAX_RETRIES = 3;

  const fetchData = async () => {
    try {
      // APIコール
      // const data = await api.getUsers();
      // return data;
    } catch (error: any) {
      const apiError: ApiError = error.body;
      
      // リトライ可能なエラーの場合
      if (isRetryableError(apiError) && retryCount < MAX_RETRIES) {
        setRetryCount(prev => prev + 1);
        
        // 指数バックオフでリトライ
        const delay = Math.pow(2, retryCount) * 1000;
        await new Promise(resolve => setTimeout(resolve, delay));
        
        return fetchData();
      }
      
      // リトライ不可またはリトライ上限に達した場合
      throw error;
    }
  };
}

/**
 * グローバルエラーハンドラーの例
 */
export function setupGlobalErrorHandler() {
  // APIクライアントのインターセプターでエラーハンドリング
  // api.interceptors.response.use(
  //   response => response,
  //   error => {
  //     const apiError: ApiError = error.body;
  //     
  //     // 認証エラーの場合は自動的にログイン画面へ
  //     if (isAuthError(apiError)) {
  //       window.location.href = '/login';
  //       return;
  //     }
  //     
  //     // その他のエラーは汎用メッセージを表示
  //     const message = getErrorMessage(apiError);
  //     // toast.error(message);
  //     
  //     return Promise.reject(error);
  //   }
  // );
}

/**
 * React Query / TanStack Queryでの使用例
 */
export function useApiWithErrorHandling() {
  // const { data, error, isLoading } = useQuery({
  //   queryKey: ['users'],
  //   queryFn: () => api.getUsers(),
  //   retry: (failureCount, error) => {
  //     // リトライ可能なエラーの場合のみリトライ
  //     const apiError: ApiError = error?.body;
  //     return failureCount < 3 && isRetryableError(apiError);
  //   },
  //   onError: (error: any) => {
  //     const apiError: ApiError = error.body;
  //     const message = getErrorMessage(apiError);
  //     // toast.error(message);
  //   }
  // });
}