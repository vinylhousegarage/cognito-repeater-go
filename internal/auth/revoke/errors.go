package revoke

import "errors"

var (
	ErrFailedToCreateRequest       = errors.New("failed to create request")
	ErrFailedToCreateRevokeRequest = errors.New("failed to create revoke request")
	ErrFailedToDecodeMetadata      = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata       = errors.New("failed to fetch metadata")
	ErrFailedToSendRevokeRequest   = errors.New("failed to send revoke request")
	ErrMissingClientCredentials    = errors.New("missing client credentials")
	ErrMissingRevocationEndpoint   = errors.New("missing revocation_endpoint")
	ErrUnexpectedStatusCode        = errors.New("unexpected response from metadata")
)
