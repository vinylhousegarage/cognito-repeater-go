package me

import (
	"cognito-repeater-go/internal/apperror"
	"net/http"
)

var (
	ErrExponentTooLarge             = apperror.New(http.StatusInternalServerError, "exponent too large to fit in int") // 500 Internal Server Error
	ErrFailedToCreateJWKSetRequest  = apperror.New(http.StatusInternalServerError, "failed to create jwks request")    // 500 Internal Server Error
	ErrFailedToCreateRequest        = apperror.New(http.StatusInternalServerError, "failed to create request")         // 500 Internal Server Error
	ErrFailedToDecodeJWKSet         = apperror.New(http.StatusBadGateway, "failed to decode jwks JSON")                // 502 Bad Gateway
	ErrFailedToDecodeJWTHeader      = apperror.New(http.StatusBadRequest, "failed to decode JWT header")               // 400 Bad Request
	ErrFailedToDecodeMetadata       = apperror.New(http.StatusBadGateway, "failed to decode metadata")                 // 502 Bad Gateway
	ErrFailedToFetchJWKSet          = apperror.New(http.StatusBadGateway, "failed to fetch jwks")                      // 502 Bad Gateway
	ErrFailedToFetchMetadata        = apperror.New(http.StatusBadGateway, "failed to fetch metadata")                  // 502 Bad Gateway
	ErrFailedToParseJWTHeader       = apperror.New(http.StatusBadRequest, "failed to parse JWT header")                // 400 Bad Request
	ErrFailedToWriteUserResponse    = apperror.New(http.StatusInternalServerError, "failed to write user response")    // 500 Internal Server Error
	ErrInvalidAudience              = apperror.New(http.StatusUnauthorized, "invalid audience")                        // 401 Unauthorized
	ErrInvalidBase64URL             = apperror.New(http.StatusBadGateway, "invalid base64url encoding")                // 502 Bad Gateway
	ErrInvalidE                     = apperror.New(http.StatusInternalServerError, "invalid e (exponent)")             // 500 Internal Server Error
	ErrInvalidExponent              = apperror.New(http.StatusInternalServerError, "invalid exponent value")           // 500 Internal Server Error
	ErrInvalidIssuer                = apperror.New(http.StatusUnauthorized, "unexpected issuer")                       // 401 Unauthorized
	ErrInvalidJWTFormat             = apperror.New(http.StatusBadRequest, "invalid JWT format")                        // 400 Bad Request
	ErrInvalidN                     = apperror.New(http.StatusInternalServerError, "invalid n (modulus)")              // 500 Internal Server Error
	ErrInvalidSigningAlg            = apperror.New(http.StatusBadRequest, "unexpected signing method")                 // 400 Bad Request
	ErrJWKSetNil                    = apperror.New(http.StatusInternalServerError, "JWKSet is nil")                    // 500 Internal Server Error
	ErrJWTParseFailed               = apperror.New(http.StatusInternalServerError, "failed to parse JWT")              // 500 Internal Server Error
	ErrKIDNotFoundInJWKSet          = apperror.New(http.StatusUnauthorized, "kid not found in JWK set")                // 401 Unauthorized
	ErrMissingAudience              = apperror.New(http.StatusUnauthorized, "audience claim is missing")               // 401 Unauthorized
	ErrMissingJWKSURI               = apperror.New(http.StatusBadGateway, "missing jwks_uri")                          // 502 Bad Gateway
	ErrMissingKID                   = apperror.New(http.StatusBadRequest, "kid not found in JWT header")               // 400 Bad Request
	ErrMissingSubject               = apperror.New(http.StatusUnauthorized, "missing subject (sub)")                   // 401 Unauthorized
	ErrTokenExpired                 = apperror.New(http.StatusUnauthorized, "token is expired")                        // 401 Unauthorized
	ErrUnexpectedMetadataStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from metadata")         // 502 Bad Gateway
)
