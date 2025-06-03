package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/health"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(d deps.RouteDependencies, cli httpclient.HTTPClient) http.Handler {
	mux := http.NewServeMux()
	auth.RegisterAuthRoutes(mux, d, cli)
	health.RegisterHealthRoutes(mux)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)
	return mux
}
