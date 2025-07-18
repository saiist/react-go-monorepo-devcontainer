package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/react-go-monorepo/backend/internal/logger"
)

// LoggingMiddleware はHTTPリクエスト/レスポンスのログを記録します
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// リクエストIDとトレースIDを生成/取得
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			
			// GCP Cloud Traceのヘッダーを取得
			traceID := r.Header.Get("X-Cloud-Trace-Context")
			if traceID == "" {
				traceID = requestID // フォールバック
			}

			// リクエストボディを読み取る（必要に応じて）
			var bodyBytes []byte
			var bodyLog interface{} = nil
			
			if r.Body != nil && shouldLogBody(r) {
				bodyBytes, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				// センシティブな情報をマスクしてログ用に準備
				if len(bodyBytes) > 0 {
					bodyLog = logger.SanitizeJSON(bodyBytes)
				}
			}

			// レスポンスライターをラップ
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			// レスポンスボディをキャプチャするライター
			responseBody := &bytes.Buffer{}
			ww.Tee(responseBody)

			// リクエスト固有のロガーを作成してコンテキストに追加
			reqLogger := log.WithRequest(requestID, traceID, r.Method, r.URL.Path, r.RemoteAddr)
			ctx := reqLogger.ToContext(r.Context())
			
			// レスポンスヘッダーにリクエストIDを追加
			ww.Header().Set("X-Request-ID", requestID)
			
			// リクエストログ
			reqLogger.Info("http_request_started",
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
				"content_length", r.ContentLength,
				"query_params", sanitizeQueryParams(r.URL.Query()),
			)
			
			// ボディがある場合はログに含める
			if bodyLog != nil {
				reqLogger.Debug("http_request_body",
					"body", bodyLog,
				)
			}

			// 次のハンドラーを実行
			next.ServeHTTP(ww, r.WithContext(ctx))
			
			// レスポンスログ
			duration := time.Since(start)
			
			// ログレベルをステータスコードに基づいて決定
			logLevel := getLogLevel(ww.Status())
			
			fields := []interface{}{
				"status", ww.Status(),
				"bytes_written", ww.BytesWritten(),
				"duration_ms", duration.Milliseconds(),
				"duration_human", duration.String(),
			}
			
			// エラーレスポンスの場合、レスポンスボディもログに含める
			if ww.Status() >= 400 && responseBody.Len() > 0 {
				var errorResponse map[string]interface{}
				if err := json.Unmarshal(responseBody.Bytes(), &errorResponse); err == nil {
					fields = append(fields, "error_response", errorResponse)
				}
			}
			
			switch logLevel {
			case "error":
				reqLogger.Error("http_request_completed", fields...)
			case "warn":
				reqLogger.Warn("http_request_completed", fields...)
			default:
				reqLogger.Info("http_request_completed", fields...)
			}
		})
	}
}

// shouldLogBody はリクエストボディをログに記録すべきかを判断します
func shouldLogBody(r *http.Request) bool {
	// コンテンツタイプをチェック
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return false
	}
	
	// JSONまたはフォームデータの場合のみログに記録
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded")
}

// sanitizeQueryParams はクエリパラメータをサニタイズします
func sanitizeQueryParams(params map[string][]string) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for key, values := range params {
		if len(values) == 1 {
			sanitized[key] = logger.SanitizeValue(key, values[0])
		} else {
			sanitizedValues := make([]interface{}, len(values))
			for i, v := range values {
				sanitizedValues[i] = logger.SanitizeValue(key, v)
			}
			sanitized[key] = sanitizedValues
		}
	}
	return sanitized
}

// getLogLevel はステータスコードに基づいてログレベルを決定します
func getLogLevel(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warn"
	default:
		return "info"
	}
}

// HealthCheckLoggingMiddleware はヘルスチェックエンドポイント用の軽量ログミドルウェア
func HealthCheckLoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ヘルスチェックは最小限のログのみ
			if r.URL.Path == "/health" {
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r)
				
				// エラーの場合のみログ
				if ww.Status() >= 500 {
					log.Error("health_check_failed",
						"status", ww.Status(),
						"remote_addr", r.RemoteAddr,
					)
				}
				return
			}
			
			// その他のエンドポイントは通常のログミドルウェアを使用
			LoggingMiddleware(log)(next).ServeHTTP(w, r)
		})
	}
}