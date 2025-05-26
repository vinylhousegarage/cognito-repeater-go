package auth_test

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/router"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestMeHandler_Integration_Success(t *testing.T) {
	t.Parallel()

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","kid":"example-kid"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"test-sub"}`))
	dummyToken := header + "." + payload + ".signature"

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()

			switch {
			case strings.Contains(req.URL.String(), "/.well-known/openid-configuration"):
				rec.WriteHeader(http.StatusOK)
				if _, err := rec.WriteString(`{"jwks_uri": "https://example.com/jwks"}`); err != nil {
					log.Printf("failed to write mock jwks_uri response: %v", err)
				}

			case strings.Contains(req.URL.String(), "/jwks"):
				rec.WriteHeader(http.StatusOK)
				if _, err := rec.WriteString(`{
					"keys": [{
						"kid": "example-kid",
						"kty": "RSA",
						"alg": "RS256",
						"use": "sig",
						"n": "example-n",
						"e": "AQAB"
					}]
				}`); err != nil {
					log.Printf("failed to write mock jwks response: %v", err)
				}

			default:
				rec.WriteHeader(http.StatusNotFound)
			}

			return rec.Result(), nil
		},
	}

	cfg := &config.Config{
		Region:           "ap-northeast-1",
		UserPoolID:       "test-pool",
		UserPoolClientID: "test-client",
	}

	router := router.NewRouter(cfg, client)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+dummyToken)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "test-sub", body["sub"])
}
