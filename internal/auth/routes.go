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
	mux.HandleFunc("/callback", callback.NewCallbackHandler(d.CallbackProvider, cli))
	mux.HandleFunc("/login", login.NewLoginHandler(d.LoginProvider, cli))
	mux.HandleFunc("/logout", logout.NewLogoutHandler(d.LogoutProvider, cli))
	mux.HandleFunc("/logout/redirect", logoutredirect.NewLogoutRedirectHandler)
	mux.HandleFunc("/me", me.NewMeHandler(d.MeProvider, cli))
	mux.HandleFunc("/revoke", revoke.NewRevokeHandler(d.RevokeProvider, cli))
	mux.HandleFunc("/whoami", whoami.NewWhoamiHandler(d.WhoamiProvider, cli))
}
