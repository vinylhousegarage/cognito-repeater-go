package login

import "errors"

var (
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToDecodeMetadata    = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata     = errors.New("failed to fetch metadata")
	ErrInvalidMetadataNoEndpoint = errors.New("invalid metadata: token_endpoint is empty")
	ErrUnexpectedStatusCode      = errors.New("unexpected status code from metadata endpoint")
)
