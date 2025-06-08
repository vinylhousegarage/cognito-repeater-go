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

	"go.uber.org/zap"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	d deps.RouteDependencies,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) {
	mux.HandleFunc("/callback", callback.NewCallbackHandler(d.CallbackProvider, cli, logger))
	mux.HandleFunc("/login", login.NewLoginHandler(d.LoginProvider, cli, logger))
	mux.HandleFunc("/logout", logout.NewLogoutHandler(d.LogoutProvider, cli, logger))
	mux.HandleFunc("/logout/redirect", logoutredirect.NewLogoutRedirectHandler(logger))
	mux.HandleFunc("/me", me.NewMeHandler(d.MeProvider, cli, logger))
	mux.HandleFunc("/revoke", revoke.NewRevokeHandler(d.RevokeProvider, cli, logger))
	mux.HandleFunc("/whoami", whoami.NewWhoamiHandler(d.WhoamiProvider, cli, logger))
}
