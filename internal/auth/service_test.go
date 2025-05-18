package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMetadataURL struct {
	URL string
}

func (m *mockMetadataURL) MetadataURL() string {
	return m.URL
}

func Test_fetchMetadata_SendsRequestToCorrectURL(t *testing.T) {
	var requestedPath string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	mock := &mockMetadataURL{URL: ts.URL}

	_, err := fetchMetadata(mock)
	assert.NoError(t, err)
	assert.Equal(t, "/",
		requestedPath, "fetchMetadata should request root path")
}

func Test_fetchMetadata_ParsesResponseIntoMap(t *testing.T) {
	expected := map[string]interface{}{
		"issuer":               "https://example.com",
		"end_session_endpoint": "https://example.com/logout",
	}
	mockJSON, _ := json.Marshal(expected)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockJSON)
	}))
	defer ts.Close()

	mock := &mockMetadataURL{URL: ts.URL}
	result, err := fetchMetadata(mock)

	assert.NoError(t, err)
	assert.Equal(t, expected["issuer"], result["issuer"])
	assert.Equal(t, expected["end_session_endpoint"], result["end_session_endpoint"])
}

func Test_interfaceToString_ReturnsOnlyStringFields(t *testing.T) {
	input := map[string]interface{}{
		"issuer":               "https://example.com",
		"end_session_endpoint": "https://example.com/logout",
		"expires_in":           3600,
		"enabled":              true,
		"algorithms":           []string{"RS256"},
	}

	expected := map[string]string{
		"issuer":               "https://example.com",
		"end_session_endpoint": "https://example.com/logout",
	}

	result := interfaceToString(input)

	assert.Equal(t, expected, result)
}

func Test_getLogoutEndpoint_ReturnsValue_WhenExists(t *testing.T) {
	input := map[string]string{
		"end_session_endpoint": "https://example.com/logout",
	}

	result, err := getLogoutEndpoint(input)

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/logout", result)
}
