package logout

import "errors"

var (
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToDecodeMetadata    = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata     = errors.New("failed to fetch metadata")
	ErrMissingEndSessionEndpoint = errors.New("missing end_session_endpoint")
	ErrUnexpectedStatusCode      = errors.New("unexpected response from metadata")
)
