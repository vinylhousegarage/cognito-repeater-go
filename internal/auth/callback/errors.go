package callback

import "errors"

var (
	ErrFailedToCreateRequest  = errors.New("failed to create request")
	ErrFailedToDecodeMetadata = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata  = errors.New("failed to fetch metadata")
	ErrMissingTokenEndpoint   = errors.New("missing token_endpoint")
	ErrInvalidState           = errors.New("invalid state")
	ErrMissingCode            = errors.New("missing code")
	ErrMissingState           = errors.New("missing state")
	ErrMissingStateCookie     = errors.New("missing oauth_state cookie")
	ErrUnexpectedStatusCode   = errors.New("unexpected response from metadata")
)
