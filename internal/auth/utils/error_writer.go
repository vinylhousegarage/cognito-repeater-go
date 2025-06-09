package utils

import (
	"net/http"

	"go.uber.org/zap"
)

func WritePlainError(w http.ResponseWriter, status int, err error, logger *zap.Logger) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	logger.Error(http.StatusText(status), zap.Error(err))
	http.Error(w, http.StatusText(status), status)
}
