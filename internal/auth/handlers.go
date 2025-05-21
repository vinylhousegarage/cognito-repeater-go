package auth

import (
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/config"
	"net/http"
)

func NewAuthHandlers(cfg config.MetadataURLProvider, client *http.Client) AuthHandlers {
	return AuthHandlers{
		LogoutURLProvider:   logout.NewLogoutService(client),
		MetadataURLProvider: cfg,
	}
}
