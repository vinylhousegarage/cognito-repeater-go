package whoami

import (
	"cognito-repeater-go/internal/apperror"
	"net/http"
)

var (
	ErrFailedToCreateRequest            = apperror.New(http.StatusInternalServerError, "failed to create request")          // 500 Internal Server Error
	ErrFailedToCreateUserinfoRequest    = apperror.New(http.StatusInternalServerError, "failed to create userinfo request") // 500 Internal Server Error
	ErrFailedToDecodeMetadata           = apperror.New(http.StatusBadGateway, "failed to decode metadata")                  // 502 Bad Gateway
	ErrFailedToFetchMetadata            = apperror.New(http.StatusBadGateway, "failed to fetch metadata")                   // 502 Bad Gateway
	ErrFailedToFetchUserinfoRequest     = apperror.New(http.StatusBadGateway, "failed to fetch userinfo request")           // 502 Bad Gateway
	ErrFailedToParseUserinfoResponse    = apperror.New(http.StatusBadGateway, "failed to parse userinfo response")          // 502 Bad Gateway
	ErrFailedToReadMetadataResponse     = apperror.New(http.StatusBadGateway, "failed to read metadata response body")      // 502 Bad Gateway
	ErrFailedToReadUserinfoResponse     = apperror.New(http.StatusBadGateway, "failed to read userinfo response body")      // 502 Bad Gateway
	ErrInvalidAuthorizationHeaderFormat = apperror.New(http.StatusBadRequest, "invalid authorization header format")        // 400 Bad Request
	ErrMissingAuthorizationHeader       = apperror.New(http.StatusUnauthorized, "missing authorization header")             // 401 Unauthorized
	ErrMissingSubject                   = apperror.New(http.StatusUnauthorized, "missing subject (sub)")                    // 401 Unauthorized
	ErrMissingUserinfoEndpoint          = apperror.New(http.StatusBadGateway, "missing userinfo_endpoint")                  // 502 Bad Gateway
	ErrUnexpectedMetadataStatusCode     = apperror.New(http.StatusBadGateway, "unexpected response from metadata")          // 502 Bad Gateway
	ErrUnexpectedUserinfoStatusCode     = apperror.New(http.StatusBadGateway, "unexpected response from userinfo")          // 502 Bad Gateway
)
