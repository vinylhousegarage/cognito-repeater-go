package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/errors"
	"cognito-repeater-go/internal/root"
)

func NewRouter(p config.MetadataURLProvider) http.Handler {
	mux := http.NewServeMux()

	auth.RegisterAuthRoutes(mux, p)
	errors.RegisterErrorRoutes(mux)
	root.RegisterRootRoutes(mux)

	return mux
}
