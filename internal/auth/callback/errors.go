package callback

import "errors"

var (
	ErrMissingCode        = errors.New("missing code")
	ErrMissingState       = errors.New("missing state")
	ErrMissingStateCookie = errors.New("missing oauth_state cookie")
	ErrInvalidState       = errors.New("invalid state")
	ErrEmptyTokenEndpoint = errors.New("token_endpoint is empty")
)
