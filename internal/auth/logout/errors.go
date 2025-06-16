package logout

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

var (
	ErrFailedToCreateRequest        = apperror.New(http.StatusInternalServerError, "failed to create request") // 500 Internal Server Error
	ErrFailedToDecodeMetadata       = apperror.New(http.StatusBadGateway, "failed to decode metadata")         // 502 Bad Gateway
	ErrFailedToFetchMetadata        = apperror.New(http.StatusBadGateway, "failed to fetch metadata")          // 502 Bad Gateway
	ErrMissingEndSessionEndpoint    = apperror.New(http.StatusBadGateway, "missing end_session_endpoint")      // 502 Bad Gateway
	ErrUnexpectedMetadataStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from metadata") // 502 Bad Gateway
)
