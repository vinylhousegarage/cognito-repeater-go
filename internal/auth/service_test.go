package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Metadata struct {
	EndSessionEndpoint string `json:"end_session_endpoint"`
}

func Test_fetchMetadata_Success(t *testing.T) {
	mockMetadata := Metadata{
		EndSessionEndpoint: "https://mock.example.com/logout",
	}

	mockJSON, err := json.Marshal(mockMetadata)
	assert.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockJSON)
	}))
	defer ts.Close()

	result, err := fetchMetadata(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, mockMetadata.EndSessionEndpoint, result.EndSessionEndpoint)
}
