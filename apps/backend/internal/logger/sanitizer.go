package logger

import (
	"encoding/json"
	"strings"
)

// センシティブなフィールド名のリスト
var sensitiveFields = []string{
	"password",
	"token",
	"secret",
	"api_key",
	"apikey",
	"access_token",
	"refresh_token",
	"authorization",
	"credit_card",
	"ssn",
	"tax_id",
}

// SanitizeValue はセンシティブな情報をマスクします
func SanitizeValue(key string, value interface{}) interface{} {
	// キー名をチェック
	lowerKey := strings.ToLower(key)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(lowerKey, sensitive) {
			return "***REDACTED***"
		}
	}

	// 値が文字列の場合、特定のパターンをチェック
	if strVal, ok := value.(string); ok {
		// JWTトークンのパターン
		if strings.HasPrefix(strVal, "Bearer ") || strings.Count(strVal, ".") == 2 {
			return "***JWT_REDACTED***"
		}
		// メールアドレスの部分マスク
		if strings.Contains(strVal, "@") && strings.Contains(key, "email") {
			parts := strings.Split(strVal, "@")
			if len(parts) == 2 && len(parts[0]) > 2 {
				masked := parts[0][:2] + strings.Repeat("*", len(parts[0])-2) + "@" + parts[1]
				return masked
			}
		}
	}

	return value
}

// SanitizeJSON はJSONデータのセンシティブな情報をマスクします
func SanitizeJSON(data []byte) json.RawMessage {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		// パースできない場合は空のJSONを返す
		return json.RawMessage("{}")
	}

	sanitized := sanitizeValue(obj)
	result, _ := json.Marshal(sanitized)
	return json.RawMessage(result)
}

// 再帰的にマップや配列をサニタイズ
func sanitizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for k, v := range val {
			sanitized[k] = SanitizeValue(k, sanitizeValue(v))
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(val))
		for i, item := range val {
			sanitized[i] = sanitizeValue(item)
		}
		return sanitized
	default:
		return val
	}
}

// MaskEmail はメールアドレスを部分的にマスクします
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	
	local := parts[0]
	domain := parts[1]
	
	if len(local) <= 2 {
		return strings.Repeat("*", len(local)) + "@" + domain
	}
	
	// 最初の2文字と最後の1文字を表示
	masked := local[:2] + strings.Repeat("*", len(local)-3) + local[len(local)-1:]
	return masked + "@" + domain
}