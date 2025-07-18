package handler

import "fmt"

// ValidationError はバリデーションエラーを表します
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError はリソースが見つからないエラーを表します
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
}

// ConflictError は競合エラーを表します
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// UnauthorizedError は認証エラーを表します
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// ForbiddenError は認可エラーを表します
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}