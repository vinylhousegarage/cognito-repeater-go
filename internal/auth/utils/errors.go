package utils

import "errors"

var (
	ErrFailedToParseForm = errors.New("failed to parse form")
	ErrMissingToken      = errors.New("token is missing from body")
)
