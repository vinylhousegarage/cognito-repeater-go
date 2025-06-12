package login

import "errors"

var (
	ErrFailedToCreateRequest        = errors.New("failed to create request")
	ErrFailedToDecodeMetadata       = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata        = errors.New("failed to fetch metadata")
	ErrFailedToParseLoginURL        = errors.New("failed to parse login endpoint URL")
	ErrMissingAuthorizationEndpoint = errors.New("missing authorization_endpoint")
	ErrUnexpectedStatusCode         = errors.New("unexpected response from metadata")
)
