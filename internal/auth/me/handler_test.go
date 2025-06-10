package me

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestNewMeHandler_MissingAuthorizationHeader(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		Region:           "ap-northeast-1",
		ClientSecret:     "client-secret",
		LogoutURI:        "https://example.com/logout",
		RedirectURI:      "https://localhost/callback",
		Scope:            "openid",
		UserPoolClientID: "abc123clientidxyz4567890123",
		UserPoolID:       "ap-northeast-1_Abc123XYZ",
	}

	handler := NewMeHandler(cfg, &testhelpers.MockHTTPClient{}, testhelpers.MockLogger)

	req := httptest.NewRequest(http.MethodPost, "/me", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "expected status 400 Bad Request")

	expectedBody := utils.ErrMissingToken.Error()
	assert.Equal(t, expectedBody, strings.TrimSpace(rr.Body.String()), "unexpected response body")
}
