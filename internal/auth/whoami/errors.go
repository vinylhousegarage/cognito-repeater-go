package whoami

import "errors"

var (
	ErrFailedToCreateRequest            = errors.New("failed to create request")              // 500 Internal Server Error
	ErrFailedToCreateUserinfoRequest    = errors.New("failed to create userinfo request")     // 500 Internal Server Error
	ErrFailedToDecodeMetadata           = errors.New("failed to decode metadata")             // 502 Bad Gateway
	ErrFailedToFetchMetadata            = errors.New("failed to fetch metadata")              // 502 Bad Gateway
	ErrFailedToFetchUserinfoRequest     = errors.New("failed to fetch userinfo request")      // 502 Bad Gateway
	ErrFailedToParseUserinfoResponse    = errors.New("failed to parse userinfo response")     // 502 Bad Gateway
	ErrFailedToReadMetadataResponse     = errors.New("failed to read metadata response body") // 502 Bad Gateway
	ErrFailedToReadUserinfoResponse     = errors.New("failed to read userinfo response body") // 502 Bad Gateway
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid authorization header format")   // 400 Bad Request
	ErrMissingAuthorizationHeader       = errors.New("missing authorization header")          // 401 Unauthorized
	ErrMissingSubject                   = errors.New("missing subject (sub)")                 // 401 Unauthorized
	ErrMissingUserinfoEndpoint          = errors.New("missing userinfo_endpoint")             // 502 Bad Gateway
	ErrSubjectIsNotString               = errors.New("subject (sub) claim is not a string")   // 401 Unauthorized
	ErrUnexpectedMetadataStatusCode     = errors.New("unexpected response from metadata")     // 502 Bad Gateway
	ErrUnexpectedUserinfoStatusCode     = errors.New("unexpected response from userinfo")     // 502 Bad Gateway
)
