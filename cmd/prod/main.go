package main

import (
	"cognito-repeater-go/internal/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	app := router.NewRouter()
	adapter := httpadapter.New(app)
	lambda.Start(func(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return adapter.Proxy(req)
	})
}
