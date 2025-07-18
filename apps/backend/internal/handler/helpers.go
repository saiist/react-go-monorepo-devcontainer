package handler

import (
	"net/http"
	"strconv"
	
	"github.com/react-go-monorepo/backend/internal/dto"
	"github.com/react-go-monorepo/backend/internal/model"
)

// parseUserQueryParams parses query parameters from the request
func parseUserQueryParams(r *http.Request) dto.UserQueryParams {
	query := r.URL.Query()
	
	params := dto.UserQueryParams{
		Page:  1,
		Limit: 20,
	}
	
	// Parse boolean parameters
	if v := query.Get("isActive"); v != "" {
		b := v == "true"
		params.IsActive = &b
	}
	
	if v := query.Get("isVerified"); v != "" {
		b := v == "true"
		params.IsVerified = &b
	}
	
	// Parse search
	params.Search = query.Get("search")
	
	// Parse pagination
	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			params.Limit = l
		}
	}
	
	return params
}

// validateCreateUserRequest validates the create user request
func validateCreateUserRequest(req *model.CreateUserRequest) error {
	if req.Email == "" {
		return &ValidationError{
			Message: "Email is required",
			Field:   "email",
		}
	}
	
	if req.Username == "" {
		return &ValidationError{
			Message: "Username is required",
			Field:   "username",
		}
	}
	
	if req.Password == "" {
		return &ValidationError{
			Message: "Password is required",
			Field:   "password",
		}
	}
	
	return nil
}