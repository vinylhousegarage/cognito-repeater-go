package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(
	cfg config.CognitoMetadataProvider,
	callbackCfg callback.CallbackHandlerProvider,
	logoutCfg logout.LogoutEndpointProvider,
	client httpclient.HTTPClient,
) http.Handler {
	mux := http.NewServeMux()

	d := deps.HandlerDependencies{
		Config:     cfg,
		HTTPClient: client,
	}

	auth.RegisterAuthRoutes(mux, d, callbackCfg, logoutCfg, client)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)

	return mux
}
