package me

import "errors"

var (
	ErrExponentTooLarge          = errors.New("exponent too large to fit in int")  // 500 Internal Server Error
	ErrFailedToCreateJWKSRequest = errors.New("failed to create jwks request")     // 500 Internal Server Error
	ErrFailedToCreateRequest     = errors.New("failed to create request")          // 500 Internal Server Error
	ErrFailedToDecodeJWKS        = errors.New("failed to decode jwks JSON")        // 502 Bad Gateway
	ErrFailedToDecodeJWTHeader   = errors.New("failed to decode JWT header")       // 400 Bad Request
	ErrFailedToDecodeMetadata    = errors.New("failed to decode metadata")         // 502 Bad Gateway
	ErrFailedToFetchJWKS         = errors.New("failed to fetch jwks")              // 502 Bad Gateway
	ErrFailedToFetchMetadata     = errors.New("failed to fetch metadata")          // 502 Bad Gateway
	ErrFailedToParseJWTHeader    = errors.New("failed to parse JWT header")        // 400 Bad Request
	ErrInvalidAudience           = errors.New("invalid audience")                  // 401 Unauthorized
	ErrInvalidBase64URL          = errors.New("invalid base64url encoding")        // 502 Bad Gateway
	ErrInvalidE                  = errors.New("invalid e (exponent)")              // 500 Internal Server Error
	ErrInvalidExponent           = errors.New("invalid exponent value")            // 500 Internal Server Error
	ErrInvalidIssuer             = errors.New("unexpected issuer")                 // 401 Unauthorized
	ErrInvalidJWTFormat          = errors.New("invalid JWT format")                // 400 Bad Request
	ErrInvalidN                  = errors.New("invalid n (modulus)")               // 500 Internal Server Error
	ErrInvalidSigningAlg         = errors.New("unexpected signing method")         // 401 Unauthorized
	ErrJWKNotFound               = errors.New("JWK not found")                     // 401 Unauthorized
	ErrJWKSetNil                 = errors.New("JWKSet is nil")                     // 401 Unauthorized
	ErrJWTParseFailed            = errors.New("failed to parse JWT")               // 500 Internal Server Error
	ErrMissingAudience           = errors.New("audience claim is missing")         // 401 Unauthorized
	ErrMissingJWKSURI            = errors.New("missing jwks_uri")                  // 502 Bad Gateway
	ErrMissingKID                = errors.New("kid not found in JWT header")       // 400 Bad Request
	ErrMissingSubject            = errors.New("missing subject (sub)")             // 401 Unauthorized
	ErrTokenExpired              = errors.New("token is expired")                  // 401 Unauthorized
	ErrUnexpectedStatusCode      = errors.New("unexpected response from metadata") // 502 Bad Gateway
)
