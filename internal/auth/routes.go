package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/logoutredirect"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
}
