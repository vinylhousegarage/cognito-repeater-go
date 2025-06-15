package router

import (
	"net/http"

	"cognito-repeater-go/internal/auth"
	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/health"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/openapi"
	"cognito-repeater-go/internal/root"
	"cognito-repeater-go/internal/simulatederrors"

	"github.com/getkin/kin-openapi/openapi3"

	"go.uber.org/zap"
)

func NewRouter(d deps.RouteDependencies, cli httpclient.HTTPClient, spec *openapi3.T, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()
	auth.RegisterAuthRoutes(mux, d, cli, logger)
	health.RegisterHealthRoutes(mux, logger)
	openapi.RegisterOpenAPIRoutes(mux, spec, logger)
	root.RegisterRootRoutes(mux, logger)
	simulatederrors.RegisterErrorRoutes(mux, logger)
	return mux
}
