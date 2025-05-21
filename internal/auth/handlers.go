package auth

import (
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/config"
	"net/http"
)

func NewAuthHandlers(cfg config.MetadataURLProvider, client *http.Client) AuthHandlers {
	return AuthHandlers{
		LoginURLProvider:    login.NewLoginClient(client),
		LogoutURLProvider:   logout.NewLogoutClient(client),
		MetadataURLProvider: cfg,
	}
}
