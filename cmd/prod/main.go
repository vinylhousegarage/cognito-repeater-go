package main

import (
	"fmt"
	"net/http"
	"os"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

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

	app := router.NewRouter(routeDeps, http.DefaultClient, logger)
	adapter := httpadapter.NewV2(app)

	lambda.Start(adapter.ProxyWithContext)
}
