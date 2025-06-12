package utils

import "errors"

var (
	ErrFailedToParseForm = errors.New("failed to parse form")       // 400 Bad Request
	ErrMissingToken      = errors.New("token is missing from body") // 400 Bad Request
)
