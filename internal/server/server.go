package server

import (
	"net/http"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"
)

func NewServer(cfg *config.Config, client *http.Client) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(cfg, client),
	}
}
