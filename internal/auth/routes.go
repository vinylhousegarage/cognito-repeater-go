package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
	"cognito-repeater-go/internal/config"
)

func RegisterAuthRoutes(mux *http.ServeMux, provider logout.LogoutEndpointProvider, p config.MetadataURLProvider) {
	mux.HandleFunc("/logout", logout.LogoutHandler(provider, p))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
}
