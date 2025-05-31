package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhoamiRouteIntegration(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			switch req.URL.Path {
			case "/.well-known/openid-configuration":
				rec.WriteHeader(http.StatusOK)
				rec.WriteString(`{"userinfo_endpoint": "` + ts.URL + `/oauth2/userinfo"}`)
			case "/oauth2/userinfo":
				if req.Header.Get("Authorization") != "Bearer valid-token" {
					rec.WriteHeader(http.StatusUnauthorized)
					return rec.Result(), nil
				}
				rec.WriteHeader(http.StatusOK)
				rec.WriteString(`{"sub": "abc123", "email": "user@example.com"}`)
			default:
				rec.WriteHeader(http.StatusNotFound)
			}
			return rec.Result(), nil
		},
	}

	mux.Handle("/whoami", WhoamiHandler(
		&testhelpers.MockMetadataURL{URL: ts.URL + "/.well-known/openid-configuration"},
		client,
	))

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/whoami", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var userinfo map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&userinfo)
	require.NoError(t, err)

	assert.Equal(t, "abc123", userinfo["sub"])
	assert.Equal(t, "user@example.com", userinfo["email"])
}
