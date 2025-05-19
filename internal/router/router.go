package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	auth.RegisterAuthRoutes(mux)
	root.RegisterRootRoutes(mux)
	simulatederrors.RegisterErrorRoutes(mux)

	return mux
}
