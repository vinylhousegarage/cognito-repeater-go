package whoami

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDeps(doFunc func(*http.Request) (*http.Response, error)) (deps.WhoamiHandlerProvider, httpclient.HTTPClient) {
	return testhelpers.MockCfg, &testhelpers.MockHTTPClient{DoFunc: doFunc}
}

func TestWhoamiHandler_Success(t *testing.T) {
	t.Parallel()

	cfg, cli := newTestDeps(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.String(), "openid-configuration") {
			body := `{"userinfo_endpoint":"https://mock/userinfo"}`
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			if _, err := rec.WriteString(body); err != nil {
				return nil, fmt.Errorf("failed to write response body: %w", err)
			}
			return rec.Result(), nil
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			body := `{"sub":"abc123","email":"user@example.com"}`
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			if _, err := rec.WriteString(body); err != nil {
				return nil, fmt.Errorf("failed to write response body: %w", err)
			}
			return rec.Result(), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", req.URL.String())
	})

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer mocktoken")
	w := httptest.NewRecorder()

	handler := WhoamiHandler(cfg, cli)
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", body["sub"])
	assert.Equal(t, "user@example.com", body["email"])
}

func TestWhoamiHandler_MissingAuthorization(t *testing.T) {
	t.Parallel()

	cfg, cli := newTestDeps(func(req *http.Request) (*http.Response, error) {
		return nil, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	w := httptest.NewRecorder()

	handler := WhoamiHandler(cfg, cli)
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "authorization header is missing")
}

func TestWhoamiHandler_MetadataFetchError(t *testing.T) {
	t.Parallel()

	cfg, cli := newTestDeps(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("metadata fetch failed")
	})

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	handler := WhoamiHandler(cfg, cli)
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "failed to get userinfo endpoint")
}

func TestWhoamiHandler_UserinfoFetchUnauthorized(t *testing.T) {
	t.Parallel()

	cfg, cli := newTestDeps(func(req *http.Request) (*http.Response, error) {
		if strings.Contains(req.URL.String(), "openid-configuration") {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			if _, err := rec.WriteString(`{"userinfo_endpoint":"https://mock/userinfo"}`); err != nil {
				return nil, fmt.Errorf("failed to write metadata: %w", err)
			}
			return rec.Result(), nil
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusUnauthorized)
			if _, err := rec.WriteString(`{"error":"unauthorized"}`); err != nil {
				return nil, fmt.Errorf("failed to write userinfo error: %w", err)
			}
			return rec.Result(), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", req.URL.String())
	})

	req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	handler := WhoamiHandler(cfg, cli)
	handler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "failed to fetch userinfo")
}
