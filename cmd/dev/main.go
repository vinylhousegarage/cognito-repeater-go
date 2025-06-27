package main

import (
	"fmt"
	"net/http"
	"os"

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
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "logger.Sync() error: %v\n", err)
		}
	}()

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

	srv := server.NewServer(routeDeps, http.DefaultClient, logger)

	logger.Info("server is starting", zap.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("server failed to start", zap.Error(err), zap.String("address", srv.Addr))
	}
}
