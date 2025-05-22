package auth

import (
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

func NewAuthHandlers(cfg config.MetadataURLProvider, client httpclient.HTTPClient) AuthHandlers {
	return AuthHandlers{
		LoginURLProvider:    login.NewLoginClient(client),
		LogoutURLProvider:   logout.NewLogoutClient(client),
		MetadataURLProvider: cfg,
	}
}
