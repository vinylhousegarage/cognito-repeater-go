package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/errors"
	"cognito-repeater-go/internal/root"
)

func NewRouter(p MetadataURLProvider) http.Handler {
	mux := http.NewServeMux()

	auth.RegisterAuthRoutes(mux, p)
	errors.RegisterErrorRoutes(mux)
	root.RegisterRootRoutes(mux)

	return mux
}
