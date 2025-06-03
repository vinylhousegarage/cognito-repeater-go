package main

import (
	"log"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	routeDeps := deps.RouteDependencies{
		CallbackProvider: cfg,
		LoginProvider:    cfg,
		LogoutProvider:   cfg,
		MeProvider:       cfg,
		RevokeProvider:   cfg,
		WhoamiProvider:   cfg,
	}

	srv := server.NewServer(routeDeps, http.DefaultClient)

	log.Println("Listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
