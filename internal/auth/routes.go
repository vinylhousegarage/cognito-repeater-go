package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/logout_redirect"
	"cognito-repeater-go/internal/config"
)

func RegisterAuthRoutes(mux *http.ServeMux, p config.MetadataURLProvider) {
	mux.HandleFunc("/logout/redirect", LogoutRedirectHandler)
}
