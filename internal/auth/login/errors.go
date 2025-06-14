package login

import (
	"cognito-repeater-go/internal/apperror"
	"net/http"
)

var (
	ErrFailedToCreateRequest        = apperror.New(http.StatusInternalServerError, "failed to create request") // 500 Internal Server Error
	ErrFailedToDecodeMetadata       = apperror.New(http.StatusBadGateway, "failed to decode metadata")         // 502 Bad Gateway
	ErrFailedToFetchMetadata        = apperror.New(http.StatusBadGateway, "failed to fetch metadata")          // 502 Bad Gateway
	ErrMissingAuthorizationEndpoint = apperror.New(http.StatusBadGateway, "missing authorization_endpoint")    // 502 Bad Gateway
	ErrUnexpectedMetadataStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from metadata") // 502 Bad Gateway
)
