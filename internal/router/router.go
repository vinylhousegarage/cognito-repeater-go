package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(cfg *config.Config, client httpclient.HTTPClient) http.Handler {
	mux := http.NewServeMux()

	handlerDeps := deps.HandlerDependencies{
		Config:     cfg,
		HTTPClient: client,
	}

	auth.RegisterAuthRoutes(mux, handlerDeps)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)

	return mux
}
