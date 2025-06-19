package utils

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

var (
	ErrFailedToParseForm                = apperror.New(http.StatusBadRequest, "failed to parse form")                // 400 Bad Request
	ErrInvalidAuthorizationHeaderFormat = apperror.New(http.StatusBadRequest, "invalid authorization header format") // 400 Bad Request
	ErrMissingAuthorizationHeader       = apperror.New(http.StatusUnauthorized, "missing authorization header")      // 401 Unauthorized
	ErrMissingToken                     = apperror.New(http.StatusBadRequest, "token is missing from body")          // 400 Bad Request
)
