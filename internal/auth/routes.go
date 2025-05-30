package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/callback"
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
	"cognito-repeater-go/internal/auth/me"
	"cognito-repeater-go/internal/httpclient"
)

type RouteDependencies struct {
	CallbackProvider callback.CallbackHandlerProvider
	LoginProvider    login.LoginHandlerProvider
	LogoutProvider   logout.LogoutHandlerProvider
	MeProvider       me.MeHandlerProvider
}

func RegisterAuthRoutes(
	mux *http.ServeMux,
	deps RouteDependencies,
	cli httpclient.HTTPClient,
) {
	mux.HandleFunc("/callback", callback.CallbackHandler(deps.CallbackProvider, cli))
	mux.HandleFunc("/login", login.LoginHandler(deps.LoginProvider, cli))
	mux.HandleFunc("/logout", logout.LogoutHandler(deps.LogoutProvider, cli))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
	mux.HandleFunc("/me", me.MeHandler(deps.MeProvider, cli))
}
