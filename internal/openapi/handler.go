package openapi

import (
	"net/http"
	"os"
	"path/filepath"

	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

func NewOpenAPIHandler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("docs", "swagger.json")

		absPath, err := filepath.Abs(path)
		if err != nil {
			response.WriteErrorResponse(w, ErrFailedToReadOpenAPISpec, logger)
			return
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			response.WriteErrorResponse(w, ErrFailedToReadOpenAPISpec, logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if n, err := w.Write(data); err != nil {
			logger.Error("failed to write response", zap.Int("bytes_written", n), zap.Error(err))
		}
	}
}
