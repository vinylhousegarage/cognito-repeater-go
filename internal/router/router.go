package router

import (
	"net/http"

	"cognito-repeater-go/internal/errors"
	"cognito-repeater-go/internal/root"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", root.PingHandler)

	errors.RegisterErrorRoutes(mux)

	return mux
}
