package me

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/testhelpers"
)

func TestNewMeHandler_MissingAuthorizationHeader(t *testing.T) {
	cfg := &config.Config{
		Region:           "ap-northeast-1",
		ClientSecret:     "client-secret",
		LogoutURI:        "https://example.com/logout",
		RedirectURI:      "https://localhost/callback",
		Scope:            "openid",
		UserPoolClientID: "abc123clientidxyz4567890123",
		UserPoolID:       "ap-northeast-1_Abc123XYZ",
	}

	handler := NewMeHandler(cfg, &testhelpers.MockHTTPClient{})

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
