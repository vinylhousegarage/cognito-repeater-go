package urlresolver

import (
	"errors"
	"net/http"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap/zaptest"
)

func TestFetchMetadata_Success(t *testing.T) {
	t.Parallel()

	mockJSON := `{
		"authorization_endpoint": "https://example.com/authorize",
		"token_endpoint": "https://example.com/token",
		"end_session_endpoint": "https://example.com/logout",
		"revocation_endpoint": "https://example.com/revoke",
		"jwks_uri": "https://example.com/jwks"
	}`

	mockClient := testhelpers.NewMockHTTPClientWithMetadata(mockJSON)
	mockLogger := zaptest.NewLogger(t)

	meta, err := FetchMetadata("https://mock-url.com/.well-known/openid-configuration", mockClient, mockLogger)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/token", meta.TokenEndpoint)
	assert.Equal(t, "https://example.com/jwks", meta.JWKSURI)
}

func TestFetchMetadata_HTTPError(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	mockLogger := zaptest.NewLogger(t)
	_, err := FetchMetadata("https://mock-url.com", mockClient, mockLogger)

	assert.ErrorContains(t, err, "failed to fetch metadata")
}

func TestFetchMetadata_BadJSON(t *testing.T) {
	t.Parallel()

	badJSON := `{ invalid json`
	mockClient := testhelpers.NewMockHTTPClientWithMetadata(badJSON)

	mockLogger := zaptest.NewLogger(t)
	_, err := FetchMetadata("https://mock-url.com", mockClient, mockLogger)

	assert.ErrorContains(t, err, "failed to decode metadata")
}
