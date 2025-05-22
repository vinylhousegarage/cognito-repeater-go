package auth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestLoginRouteIsRegisteredInProductionRouter(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"authorization_endpoint":"https://example.com/oauth2/authorize"}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		},
	}

	router := router.NewRouter(testhelpers.MockCfg, mockClient)

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "oauth2/authorize")
}
