// @title cognito-repeater-go
// @version 1.0
// @description A relay API for AWS Cognito login/logout/token verification
// @description Implemented in Go 1.24.3.
// @termsOfService https://example.com/terms/

// @contact.name vinylhousegarage
// @contact.url https://github.com/vinylhousegarage/cognito-repeater-go

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host cognito-repeater-go.com
// @BasePath /
// @schemes https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
