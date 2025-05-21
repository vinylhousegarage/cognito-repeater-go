package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
	"cognito-repeater-go/internal/config"
)

type AuthHandlers struct {
	LoginURLProvider    login.LoginURLProvider
	LogoutURLProvider   logout.LogoutURLProvider
	MetadataURLProvider config.MetadataURLProvider
}

func RegisterAuthRoutes(mux *http.ServeMux, h AuthHandlers) {
	mux.HandleFunc("/login", login.LoginHandler(h.LoginURLProvider, h.MetadataURLProvider))
	mux.HandleFunc("/logout", logout.LogoutHandler(h.LogoutURLProvider, h.MetadataURLProvider))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
}
