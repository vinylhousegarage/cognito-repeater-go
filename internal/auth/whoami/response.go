package whoami

// UserInfoResponse represents the response from the /whoami endpoint.
// This is typically returned after verifying the access_token against Cognito's userinfo endpoint.
type UserInfoResponse struct {
	Sub string `json:"sub" example:"12345678-aaaa-bbbb-cccc-1234567890ab"`
}
