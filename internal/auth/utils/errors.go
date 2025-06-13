package utils

import "errors"

var (
	ErrFailedToParseForm = apperror.New(http.StatusBadRequest, "failed to parse form")       // 400 Bad Request
	ErrMissingToken      = apperror.New(http.StatusBadRequest, "token is missing from body") // 400 Bad Request
)
