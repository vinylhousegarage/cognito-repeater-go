package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/callback"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
)

func RegisterAuthRoutes(mux *http.ServeMux, d deps.HandlerDependencies) {
	mux.HandleFunc("/callback", callback.CallbackHandler(d))
	mux.HandleFunc("/login", login.LoginHandler(d))
	mux.HandleFunc("/logout", logout.LogoutHandler(d))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
}
