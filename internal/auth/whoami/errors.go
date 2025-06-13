package whoami

import "errors"

var (
	ErrFailedToCreateRequest            = errors.New("failed to create request")
	ErrFailedToCreateUserinfoRequest    = errors.New("failed to create userinfo request")
	ErrFailedToDecodeMetadata           = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata            = errors.New("failed to fetch metadata")
	ErrFailedToFetchUserinfoRequest     = errors.New("failed to fetch userinfo request")
	ErrFailedToParseUserinfoResponse    = errors.New("failed to parse userinfo response")
	ErrFailedToReadMetadataResponse     = errors.New("failed to read metadata response body")
	ErrFailedToReadUserinfoResponse     = errors.New("failed to read userinfo response body")
	ErrInvalidAuthorizationHeaderFormat = errors.New("invalid authorization header format")
	ErrMissingAuthorizationHeader       = errors.New("missing authorization header")
	ErrMissingUserinfoEndpoint          = errors.New("missing userinfo_endpoint")
	ErrUnexpectedMetadataStatusCode     = errors.New("unexpected response from metadata")
	ErrUnexpectedUserinfoStatusCode     = errors.New("unexpected response from userinfo")
)
