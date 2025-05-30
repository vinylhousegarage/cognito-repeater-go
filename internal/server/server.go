package server

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/router"
)

func NewServer(d deps.RouteDependencies, cli httpclient.HTTPClient) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(d, cli),
	}
}
