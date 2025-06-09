package callback

import "errors"

var (
	ErrInvalidState              = errors.New("invalid state")
	ErrMissingCode               = errors.New("missing code")
	ErrMissingState              = errors.New("missing state")
	ErrMissingStateCookie        = errors.New("missing oauth_state cookie")
)
