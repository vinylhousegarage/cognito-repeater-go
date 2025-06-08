package urlresolver

import "errors"

var (
	// HTTP/Request errors
	ErrFailedToCreateRequest  = errors.New("failed to create request")
	ErrFailedToDecodeMetadata = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata  = errors.New("failed to fetch metadata")
	ErrUnexpectedStatusCode   = errors.New("unexpected status code from metadata endpoint")

	// Metadata validation errors
	ErrMissingAuthorizationEndpoint = errors.New("missing authorization_endpoint in metadata")
	ErrMissingEndSessionEndpoint    = errors.New("missing end_session_endpoint in metadata")
	ErrMissingJWKSURI               = errors.New("missing jwks_uri in metadata")
	ErrMissingRevocationEndpoint    = errors.New("missing revocation_endpoint in metadata")
	ErrMissingTokenEndpoint         = errors.New("missing token_endpoint in metadata")
)
