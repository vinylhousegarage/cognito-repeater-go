package main

import (
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter/v2"

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

	app := router.NewRouter(routeDeps, http.DefaultClient, logger)
	adapter := v2adapter.NewV2(app)

	lambda.Start(adapter.ProxyWithContext)
}
