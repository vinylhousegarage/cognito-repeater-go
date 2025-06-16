package revoke

import (
	"net/http"

	"cognito-repeater-go/internal/apperror"
)

var (
	ErrFailedToCreateRequest        = apperror.New(http.StatusInternalServerError, "failed to create request")        // 500 Internal Server Error
	ErrFailedToCreateRevokeRequest  = apperror.New(http.StatusInternalServerError, "failed to create revoke request") // 500 Internal Server Error
	ErrFailedToDecodeMetadata       = apperror.New(http.StatusBadGateway, "failed to decode metadata")                // 502 Bad Gateway
	ErrFailedToFetchMetadata        = apperror.New(http.StatusBadGateway, "failed to fetch metadata")                 // 502 Bad Gateway
	ErrFailedToSendRevokeRequest    = apperror.New(http.StatusBadGateway, "failed to send revoke request")            // 502 Bad Gateway
	ErrMissingClientCredentials     = apperror.New(http.StatusBadRequest, "missing client credentials")               // 400 Bad Request
	ErrMissingRevocationEndpoint    = apperror.New(http.StatusBadGateway, "missing revocation_endpoint")              // 502 Bad Gateway
	ErrUnexpectedRevokeStatusCode   = apperror.New(http.StatusBadGateway, "unexpected response from revoke")          // 502 Bad Gateway
	ErrUnexpectedMetadataStatusCode = apperror.New(http.StatusBadGateway, "unexpected response from metadata")        // 502 Bad Gateway
)
