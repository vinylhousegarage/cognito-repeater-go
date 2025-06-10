package me

import "errors"

var (
	ErrFailedToCreateJWKSRequest = errors.New("failed to create jwks request")
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToDecodeJWKS        = errors.New("failed to decode jwks JSON")
	ErrFailedToDecodeMetadata    = errors.New("failed to decode metadata")
	ErrFailedToFetchJWKS         = errors.New("failed to fetch jwks")
	ErrFailedToFetchMetadata     = errors.New("failed to fetch metadata")
	ErrFailedToParseLoginURL     = errors.New("failed to parse login endpoint URL")
	ErrMissingJWKSURI            = errors.New("missing jwks_uri in metadata response")
	ErrUnexpectedStatusCode      = errors.New("unexpected status code from metadata endpoint")
)
