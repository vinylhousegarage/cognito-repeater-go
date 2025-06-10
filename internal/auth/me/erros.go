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
	ErrInvalidAudience           = errors.New("invalid audience")
	ErrInvalidIssuer             = errors.New("unexpected issuer")
	ErrInvalidSigningAlg         = errors.New("unexpected signing method")
	ErrJWTParseFailed            = errors.New("failed to parse JWT")
	ErrMissingAudience           = errors.New("audience claim is missing")
	ErrMissingJWKSURI            = errors.New("missing jwks_uri in metadata response")
	ErrMissingSubject            = errors.New("missing subject (sub)")
	ErrTokenExpired              = errors.New("token is expired")
	ErrUnexpectedStatusCode      = errors.New("unexpected status code from metadata endpoint")
)
