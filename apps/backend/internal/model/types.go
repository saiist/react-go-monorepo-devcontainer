package model

import "time"

// User represents a user model
type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	FullName   string    `json:"fullName,omitempty"`
	AvatarURL  string    `json:"avatarUrl,omitempty"`
	IsActive   bool      `json:"isActive"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FullName  string `json:"fullName,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email      *string `json:"email,omitempty"`
	Username   *string `json:"username,omitempty"`
	FullName   *string `json:"fullName,omitempty"`
	AvatarURL  *string `json:"avatarUrl,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
	IsVerified *bool   `json:"isVerified,omitempty"`
}

// Error represents an API error response
type Error struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	TraceID   string      `json:"traceId"`
	Path      string      `json:"path"`
}