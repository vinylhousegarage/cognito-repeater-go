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

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
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

	logger.Info("server is starting", zap.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("server failed to start", zap.Error(err), zap.String("address", srv.Addr))
		os.Exit(1)
	}
}
