package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestRevokeRoute_Integration(t *testing.T) {
	t.Parallel()

	var ts *httptest.Server

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, `{"revocation_endpoint": "%s/oauth2/revoke"}`, ts.URL)
			if err != nil {
				t.Errorf("failed to write openid config: %v", err)
			}
		case "/oauth2/revoke":
			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			if r.FormValue("token") != "dummy-token" {
				http.Error(w, "invalid token", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	})

	ts = httptest.NewServer(handler)
	defer ts.Close()

	mockDeps := testhelpers.NewMockRouteDependencies()

	r := router.NewRouter(mockDeps, http.DefaultClient)

	form := url.Values{}
	form.Set("token", "dummy-token")

	req := httptest.NewRequest(http.MethodPost, "/revoke", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	resp := rec.Result()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
