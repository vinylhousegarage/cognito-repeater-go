package utils

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func ExtractFormValue(r *http.Request, logger *zap.Logger) (string, error) {
	if err := r.ParseForm(); err != nil {
		logger.Error("failed to parse form", zap.Error(err))
		return "", ErrFailedToParseForm
	}

	token := r.FormValue("token")
	if token == "" {
		return "", ErrMissingToken
	}

	return token, nil
}
