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
