package me

// UserResponse is the response returned by the /me endpoint
type UserResponse struct {
    Sub string `json:"sub" example:"abc-123"`
}
