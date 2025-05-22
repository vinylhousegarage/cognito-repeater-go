package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(cfg *config.Config, client *http.Client) http.Handler {
	mux := http.NewServeMux()

	authHandlers := auth.NewAuthHandlers(cfg, client)

	auth.RegisterAuthRoutes(mux, authHandlers)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)

	return mux
}
