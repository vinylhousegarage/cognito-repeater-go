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

func RegisterAuthRoutes(
	mux *http.ServeMux,
	callbackCfg callback.CallbackHandlerProvider,
	loginCfg login.LoginHandlerProvider,
	logoutCfg logout.LogoutHandlerProvider,
	meCfg me.MeHandlerProvider,
	h httpclient.HTTPClient,
) {
	mux.HandleFunc("/callback", callback.CallbackHandler(callbackCfg, h))
	mux.HandleFunc("/login", login.LoginHandler(loginCfg, h))
	mux.HandleFunc("/logout", logout.LogoutHandler(logoutCfg, h))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
	mux.HandleFunc("/me", me.MeHandler(meCfg, h))
}
