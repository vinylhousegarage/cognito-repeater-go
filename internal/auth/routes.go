package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
)

func RegisterAuthRoutes(mux *http.ServeMux, p logout.LogoutEndpointProvider) {
	mux.HandleFunc("/logout", logout.LogoutHandler(p))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
}
