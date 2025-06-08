package server

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/router"

	"go.uber.org/zap"
)

func NewServer(d deps.RouteDependencies, cli httpclient.HTTPClient, logger *zap.Logger) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(d, cli, logger),
	}
}
