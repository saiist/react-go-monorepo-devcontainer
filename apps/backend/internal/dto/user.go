package dto

// UserQueryParams represents query parameters for user listing
type UserQueryParams struct {
	IsActive   *bool
	IsVerified *bool
	Search     string
	Page       int
	Limit      int
}