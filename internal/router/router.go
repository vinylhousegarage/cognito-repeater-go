package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/errors"
	"cognito-repeater-go/internal/root"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	auth.RegisterAuthRoutes(mux)
	errors.RegisterErrorRoutes(mux)
	root.RegisterRootRoutes(mux)

	return mux
}
