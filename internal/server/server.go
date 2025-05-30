package server

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/router"
)

func NewServer(deps auth.RouteDependencies, cli httpclient.HTTPClient) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(deps, cli),
	}
}
