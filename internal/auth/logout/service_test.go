package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestGetLogoutURL_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": "https://example.com/logout"}`))
	}))
	defer ts.Close()

	endpoint, err := GetLogoutURL(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", endpoint)
}

func TestGetLogoutURL_RequestCreationError(t *testing.T) {
	t.Parallel()

	_, err := GetLogoutURL(":", http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToCreateRequest)
}

func TestGetLogoutURL_HTTPClientError(t *testing.T) {
	t.Parallel()

	_, err := GetLogoutURL("http://invalid.host.local", http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
}

func TestGetLogoutURL_StatusCodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetLogoutURL(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnexpectedMetadataStatusCode)
}

func TestGetLogoutURL_DecodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer ts.Close()

	_, err := GetLogoutURL(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
}

func TestGetLogoutURL_MissingEndSessionEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": ""}`))
	}))
	defer ts.Close()

	_, err := GetLogoutURL(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingEndSessionEndpoint)
}
