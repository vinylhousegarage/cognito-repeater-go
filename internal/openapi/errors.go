package openapi

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

var (
	ErrFailedToLoadOpenAPISpec     = apperror.New(http.StatusInternalServerError, "Failed to load OpenAPI spec")     // 500 Internal Server Error
	ErrFailedToMarshalOpenAPISpec  = apperror.New(http.StatusInternalServerError, "failed to marshal OpenAPI spec")  // 500 Internal Server Error
	ErrFailedToReadOpenAPISpec     = apperror.New(http.StatusInternalServerError, "Failed to read OpenAPI spec")     // 500 Internal Server Error
	ErrFailedToValidateOpenAPISpec = apperror.New(http.StatusInternalServerError, "Failed to validate OpenAPI spec") // 500 Internal Server Error
)
