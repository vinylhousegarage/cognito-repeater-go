package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(deps auth.RouteDependencies, cli httpclient.HTTPClient) http.Handler {
	mux := http.NewServeMux()
	auth.RegisterAuthRoutes(mux, deps, cli)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)
	return mux
}
