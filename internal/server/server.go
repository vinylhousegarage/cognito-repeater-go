package server

import (
	"net/http"

	"cognito-repeater-go/internal/router"
)

func NewServer(p MetadataURLProvider) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(p),
	}
}
