package openapi

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

var (
	ErrFailedToReadOpenAPISpec     = apperror.New(http.StatusInternalServerError, "Failed to read OpenAPI spec")     // 500 Internal Server Error
	ErrFailedToLoadOpenAPISpec     = apperror.New(http.StatusInternalServerError, "Failed to load OpenAPI spec")     // 500 Internal Server Error
	ErrFailedToValidateOpenAPISpec = apperror.New(http.StatusInternalServerError, "Failed to validate OpenAPI spec") // 500 Internal Server Error
)
