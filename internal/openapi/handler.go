package openapi

import (
	"net/http"

	"cognito-repeater-go/internal/response"

	"github.com/getkin/kin-openapi/openapi3"

	"go.uber.org/zap"
)

func NewOpenAPIHandler(spec *openapi3.T, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := spec.MarshalJSON()
		if err != nil {
			response.WriteErrorResponse(w, ErrFailedToMarshalOpenAPISpec, logger)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, err := w.Write(data)
		if err != nil {
			logger.Error("failed to write OpenAPI response", zap.Int("bytes_written", n), zap.Error(err))
		}
	}
}
