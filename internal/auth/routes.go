package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/logout_redirect"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/logout/redirect", logout_redirect.LogoutRedirectHandler)
}
