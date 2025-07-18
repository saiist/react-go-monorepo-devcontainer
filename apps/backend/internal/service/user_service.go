package service

import (
	"errors"
)

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserService handles user business logic
type UserService struct {
	// repository would be injected here
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}