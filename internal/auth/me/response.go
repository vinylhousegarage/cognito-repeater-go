package me

// UserResponse is the response returned by the /me endpoint
type UserResponse struct {
	Sub string `json:"sub" example:"12345678-aaaa-bbbb-cccc-1234567890ab"`
}
