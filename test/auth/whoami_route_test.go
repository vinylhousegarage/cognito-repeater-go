package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestWhoamiRouteIntegration(t *testing.T) {
	t.Parallel()

	r := router.NewRouter(
		testhelpers.NewMockRouteDependencies(),
		&testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				rec := httptest.NewRecorder()

				switch req.URL.Path {
				case "/.well-known/openid-configuration":
					rec.WriteHeader(http.StatusOK)
					if _, err := rec.WriteString(`{"userinfo_endpoint": "http://example.com/oauth2/userinfo"}`); err != nil {
						return nil, err
					}
				case "/oauth2/userinfo":
					if req.Header.Get("Authorization") != "Bearer valid-token" {
						rec.WriteHeader(http.StatusUnauthorized)
						return rec.Result(), nil
					}
					rec.WriteHeader(http.StatusOK)
					if _, err := rec.WriteString(`{"sub": "abc123", "email": "user@example.com"}`); err != nil {
						return nil, err
					}
				default:
					rec.WriteHeader(http.StatusNotFound)
				}

				return rec.Result(), nil
			},
		},
	)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/whoami", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
