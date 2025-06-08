package urlresolver

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap/zaptest"
)

func newMockHTTPClientWithMetadata(json string) httpclient.HTTPClient {
	return &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(json)),
			}, nil
		},
	}
}

func TestFetchMetadata_Success(t *testing.T) {
	t.Parallel()

	mockJSON := `{
		"authorization_endpoint": "https://example.com/authorize",
		"token_endpoint": "https://example.com/token",
		"end_session_endpoint": "https://example.com/logout",
		"revocation_endpoint": "https://example.com/revoke",
		"jwks_uri": "https://example.com/jwks"
	}`

	mockClient := newMockHTTPClientWithMetadata(mockJSON)
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
	mockClient := newMockHTTPClientWithMetadata(badJSON)

	mockLogger := zaptest.NewLogger(t)
	_, err := FetchMetadata("https://mock-url.com", mockClient, mockLogger)

	assert.ErrorContains(t, err, "failed to decode metadata")
}

func TestFetchMetadata_UnexpectedStatusCode(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		},
	}

	mockLogger := zaptest.NewLogger(t)
	_, err := FetchMetadata("https://mock-url.com", mockClient, mockLogger)

	assert.ErrorContains(t, err, "unexpected status code")
}
