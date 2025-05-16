package router

import (
	"net/http"

	"cognito-repeater-go/internal/errors"
	"cognito-repeater-go/internal/root"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	errors.RegisterErrorRoutes(mux)
	root.RegisterRootRoutes(mux)

	return mux
}
