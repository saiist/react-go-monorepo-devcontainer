package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/react-go-monorepo/backend/internal/logger"
	"github.com/react-go-monorepo/backend/internal/model"
	"github.com/react-go-monorepo/backend/internal/service"
)

// UserHandler はユーザー関連のHTTPハンドラー
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers はユーザー一覧を取得します
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	
	// クエリパラメータの解析
	params := parseUserQueryParams(r)
	
	// ログ: API呼び出し開始
	log.Info("fetching_users",
		"filters", map[string]interface{}{
			"is_active":   params.IsActive,
			"is_verified": params.IsVerified,
			"search":      params.Search,
			"page":        params.Page,
			"limit":       params.Limit,
		},
	)
	
	// サービス層の呼び出し
	users, total, err := h.userService.GetUsers(ctx, params)
	if err != nil {
		log.WithError(err).Error("failed_to_fetch_users")
		respondError(w, r, err)
		return
	}
	
	// 成功ログ
	log.Info("users_fetched_successfully",
		"count", len(users),
		"total", total,
		"page", params.Page,
	)
	
	// レスポンス
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  users,
		"total": total,
		"page":  params.Page,
		"limit": params.Limit,
	})
}

// CreateUser は新しいユーザーを作成します（管理者のみ）
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid_request_body",
			"error", err.Error(),
		)
		respondError(w, r, &ValidationError{
			Message: "Invalid request body",
			Field:   "body",
		})
		return
	}
	
	// リクエストログ（センシティブ情報はマスク）
	log.Info("creating_user",
		"email", logger.MaskEmail(req.Email),
		"username", req.Username,
		"full_name", req.FullName,
	)
	
	// バリデーション
	if err := validateCreateUserRequest(&req); err != nil {
		log.Warn("user_validation_failed",
			"error", err.Error(),
			"email", logger.MaskEmail(req.Email),
		)
		respondError(w, r, err)
		return
	}
	
	// ユーザー作成
	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		// エラーの種類によってログレベルを変更
		if errors.Is(err, service.ErrUserAlreadyExists) {
			log.Warn("user_already_exists",
				"email", logger.MaskEmail(req.Email),
				"username", req.Username,
			)
		} else {
			log.WithError(err).Error("user_creation_failed",
				"email", logger.MaskEmail(req.Email),
			)
		}
		respondError(w, r, err)
		return
	}
	
	// 成功ログ
	log.Info("user_created_successfully",
		"user_id", user.ID,
		"email", logger.MaskEmail(user.Email),
		"username", user.Username,
	)
	
	respondJSON(w, http.StatusCreated, user)
}

// GetUserByID はIDによってユーザーを取得します
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	
	userID := chi.URLParam(r, "userId")
	if _, err := uuid.Parse(userID); err != nil {
		log.Warn("invalid_user_id",
			"user_id", userID,
			"error", err.Error(),
		)
		respondError(w, r, &ValidationError{
			Message: "Invalid user ID format",
			Field:   "userId",
		})
		return
	}
	
	log.Debug("fetching_user_by_id", "user_id", userID)
	
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Info("user_not_found", "user_id", userID)
		} else {
			log.WithError(err).Error("failed_to_fetch_user", "user_id", userID)
		}
		respondError(w, r, err)
		return
	}
	
	// 個人情報へのアクセスログ（監査目的）
	log.Info("user_accessed",
		"user_id", user.ID,
		"accessed_by", getCurrentUserID(ctx), // 認証ミドルウェアから取得
		"action", "view_user_details",
	)
	
	respondJSON(w, http.StatusOK, user)
}

// UpdateUser はユーザー情報を更新します
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	
	userID := chi.URLParam(r, "userId")
	
	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid_update_request", "error", err.Error())
		respondError(w, r, &ValidationError{
			Message: "Invalid request body",
			Field:   "body",
		})
		return
	}
	
	// 更新内容のログ（変更フィールドのみ）
	updates := make(map[string]interface{})
	if req.Email != nil {
		updates["email"] = logger.MaskEmail(*req.Email)
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsVerified != nil {
		updates["is_verified"] = *req.IsVerified
	}
	
	log.Info("updating_user",
		"user_id", userID,
		"updates", updates,
		"updated_by", getCurrentUserID(ctx),
	)
	
	user, err := h.userService.UpdateUser(ctx, userID, &req)
	if err != nil {
		log.WithError(err).Error("user_update_failed",
			"user_id", userID,
		)
		respondError(w, r, err)
		return
	}
	
	// 監査ログ: 重要な変更を記録
	if req.IsActive != nil || req.IsVerified != nil {
		log.Warn("user_status_changed",
			"user_id", userID,
			"updated_by", getCurrentUserID(ctx),
			"changes", updates,
			"timestamp", time.Now().UTC(),
		)
	}
	
	log.Info("user_updated_successfully",
		"user_id", user.ID,
		"updated_fields", len(updates),
	)
	
	respondJSON(w, http.StatusOK, user)
}

// DeleteUser はユーザーを削除します（管理者のみ）
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	
	userID := chi.URLParam(r, "userId")
	
	// 削除前の確認ログ
	log.Warn("user_deletion_requested",
		"user_id", userID,
		"requested_by", getCurrentUserID(ctx),
		"timestamp", time.Now().UTC(),
	)
	
	// ユーザー情報を取得（監査目的）
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		log.WithError(err).Error("failed_to_fetch_user_for_deletion", "user_id", userID)
		respondError(w, r, err)
		return
	}
	
	// 削除実行
	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		log.WithError(err).Error("user_deletion_failed",
			"user_id", userID,
		)
		respondError(w, r, err)
		return
	}
	
	// 削除成功の監査ログ（重要操作）
	log.Warn("user_deleted",
		"user_id", userID,
		"user_email", logger.MaskEmail(user.Email),
		"user_username", user.Username,
		"deleted_by", getCurrentUserID(ctx),
		"timestamp", time.Now().UTC(),
	)
	
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, r *http.Request, err error) {
	traceID := r.Header.Get("X-Cloud-Trace-Context")
	
	var status int
	var errorResponse model.Error
	
	switch e := err.(type) {
	case *ValidationError:
		status = http.StatusBadRequest
		errorResponse = model.Error{
			Code:    "VALIDATION_ERROR",
			Message: e.Error(),
			Details: map[string]interface{}{
				"field": e.Field,
			},
		}
	case *NotFoundError:
		status = http.StatusNotFound
		errorResponse = model.Error{
			Code:    "NOT_FOUND",
			Message: e.Error(),
		}
	case *ConflictError:
		status = http.StatusConflict
		errorResponse = model.Error{
			Code:    "CONFLICT",
			Message: e.Error(),
		}
	default:
		status = http.StatusInternalServerError
		errorResponse = model.Error{
			Code:    "INTERNAL_ERROR",
			Message: "An internal error occurred",
		}
	}
	
	// エラーレスポンスにメタデータを追加
	errorResponse.Timestamp = time.Now().UTC()
	errorResponse.TraceID = traceID
	errorResponse.Path = r.URL.Path
	
	respondJSON(w, status, errorResponse)
}

func getCurrentUserID(ctx context.Context) string {
	// 認証ミドルウェアから現在のユーザーIDを取得
	// 実装は認証システムに依存
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return "anonymous"
}