package auth

import (
	"net/http"

	"cognito-repeater-go/internal/auth/callback"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/logoutredirect"
	"cognito-repeater-go/internal/auth/me"
	"cognito-repeater-go/internal/auth/revoke"
	"cognito-repeater-go/internal/auth/whoami"
	"cognito-repeater-go/internal/httpclient"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	d deps.RouteDependencies,
	cli httpclient.HTTPClient,
) {
	mux.HandleFunc("/callback", callback.CallbackHandler(d.CallbackProvider, cli))
	mux.HandleFunc("/login", login.LoginHandler(d.LoginProvider, cli))
	mux.HandleFunc("/logout", logout.LogoutHandler(d.LogoutProvider, cli))
	mux.HandleFunc("/logout/redirect", logoutredirect.LogoutRedirectHandler)
	mux.HandleFunc("/me", me.MeHandler(d.MeProvider, cli))
	mux.HandleFunc("/revoke", revoke.RevokeHandler(d.RevokeProvider, cli))
	mux.HandleFunc("/whoami", whoami.WhoamiHandler(d.WhoamiProvider, cli))
}
