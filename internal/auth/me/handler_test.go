package me

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/test/testhelpers"
)

type mockConfig struct{}

func (m *mockConfig) MetadataURL() string { return "https://example.com/.well-known/jwks.json" }
func (m *mockConfig) Issuer() string      { return "test-issuer" }
func (m *mockConfig) Audience() string    { return "test-audience" }

func TestMeHandler_MissingAuthorizationHeader(t *testing.T) {
	handler := MeHandler(deps.HandlerDependencies{
		Config:     &mockConfig{},
		HTTPClient: &testhelpers.MockHTTPClient{},
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "missing or malformed access token") {
		t.Errorf("unexpected response body: %s", rr.Body.String())
	}
}
