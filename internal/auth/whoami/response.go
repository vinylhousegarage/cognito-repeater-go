package whoami

// UserInfoResponse represents the response from the /whoami endpoint.
// This is typically returned after verifying the access_token against Cognito's userinfo endpoint.
type UserInfoResponse struct {
	Sub               string `json:"sub" example:"abc-123-def-456"`
	Email             string `json:"email" example:"user@example.com"`
	EmailVerified     bool   `json:"email_verified" example:"true"`
	Name              string `json:"name,omitempty" example:"Taro Yamada"`
	PreferredUsername string `json:"preferred_username,omitempty" example:"taro123"`
}
