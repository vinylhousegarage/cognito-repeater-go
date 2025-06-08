package login

import "errors"

var (
	ErrCreateRequest         = errors.New("failed to create request")
	ErrDecodeMetadata        = errors.New("failed to decode metadata")
	ErrFetchMetadata         = errors.New("failed to fetch metadata")
	ErrMissingAuthorization  = errors.New("authorization_endpoint is empty in metadata")
	ErrUnexpectedStatusCode  = errors.New("unexpected status code from metadata endpoint")
)
