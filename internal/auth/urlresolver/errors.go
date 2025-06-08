package urlresolver

import "errors"

var (
	ErrFailedToCreateRequest     = errors.New("failed to create request")
	ErrFailedToDecodeMetadata    = errors.New("failed to decode metadata")
	ErrFailedToFetchMetadata     = errors.New("failed to fetch metadata")
)
