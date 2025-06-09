package login

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

var mockLogger = zap.NewNop()

func TestGetLoginURL_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
	}))
	defer ts.Close()

	endpoint, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/oauth2/authorize", endpoint)
}

func TestGetLoginURL_RequestCreationError(t *testing.T) {
	t.Parallel()

	_, err := GetLoginURL(":", http.DefaultClient, mockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToCreateRequest)
}

func TestGetLoginURL_HTTPClientError(t *testing.T) {
	t.Parallel()

	_, err := GetLoginURL("http://invalid.host.local", http.DefaultClient, mockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
}

func TestGetLoginURL_StatusCodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnexpectedStatusCode)
}

func TestGetLoginURL_DecodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer ts.Close()

	_, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
}

func TestGetLoginURL_MissingAuthorizationEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"authorization_endpoint": ""}`))
	}))
	defer ts.Close()

	_, err := GetLoginURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingAuthorizationEndpoint)
}
