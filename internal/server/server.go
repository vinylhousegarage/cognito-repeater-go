package server

import (
	"log"
	"net/http"

	"cognito-repeater-go/internal/router"
)

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: router.NewRouter(),
	}
}
