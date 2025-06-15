package openapi

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

const msg = "Failed to read OpenAPI spec"

var (
	ErrFailedToReadOpenAPISpec = apperror.New(http.StatusInternalServerError, msg) // 500 Internal Server Error
)
