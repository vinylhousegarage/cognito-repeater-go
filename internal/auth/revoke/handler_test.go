package revoke

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRevokeHandlerProvider struct {
	clientID     string
	clientSecret string
}

func (m *mockRevokeHandlerProvider) MetadataURL() string {
	return "https://mock-domain/.well-known/openid-configuration"
}

func (m *mockRevokeHandlerProvider) ClientSecretValue() string {
	return m.clientSecret
}

func (m *mockRevokeHandlerProvider) UserPoolClientIDValue() string {
	return m.clientID
}

func TestNewRevokeHandler_Success(t *testing.T) {
	t.Parallel()

	mockMetadataResp := `{"revocation_endpoint":"https://mock-domain/oauth2/revoke"}`

	mockRevokeResp := &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/.well-known/openid-configuration") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(mockMetadataResp)),
				}, nil
			}
			if req.Method == http.MethodPost && strings.Contains(req.URL.Path, "/oauth2/revoke") {
				if u, p, ok := req.BasicAuth(); !ok || u != "mock-client-id" || p != "mock-secret" {
					require.True(t, ok, "expected basic auth credentials")
					require.Equal(t, "mock-client-id", u)
					require.Equal(t, "mock-secret", p)
				}
				defer func() {
					if err := req.Body.Close(); err != nil {
						t.Logf("failed to close request body: %v", err)
					}
				}()
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)

				assert.Contains(t, string(body), "token=mock-refresh-token")
				return mockRevokeResp, nil
			}
			require.FailNow(t, "unexpected request", req.URL.String())
			return nil, nil
		},
	}

	form := "token=mock-refresh-token"
	req := httptest.NewRequest(http.MethodPost, "/revoke", bytes.NewBufferString(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := NewRevokeHandler(&mockRevokeHandlerProvider{
		clientID:     "mock-client-id",
		clientSecret: "mock-secret",
	}, client, testhelpers.MockLogger)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func buildRequest() (*http.Request, *httptest.ResponseRecorder) {
	form := "token=mock-refresh-token"
	req := httptest.NewRequest(http.MethodPost, "/revoke", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	return req, rr
}

func TestNewRevokeHandler_Errors(t *testing.T) {
	t.Parallel()

	t.Run("missing-client-credentials", func(t *testing.T) {
		t.Parallel()

		mockProvider := &mockRevokeHandlerProvider{clientID: "", clientSecret: ""}
		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"revocation_endpoint":"https://mock-domain/oauth2/revoke"}`)),
					}, nil
				}
				return nil, nil
			},
		}

		req, rr := buildRequest()

		handler := NewRevokeHandler(mockProvider, client, testhelpers.MockLogger)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), ErrMissingClientCredentials.Error())
	})

	t.Run("failed-to-fetch-metadata", func(t *testing.T) {
		t.Parallel()

		mockProvider := &mockRevokeHandlerProvider{}
		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, ErrFailedToFetchMetadata
			},
		}

		req, rr := buildRequest()

		handler := NewRevokeHandler(mockProvider, client, testhelpers.MockLogger)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadGateway, rr.Code)
		assert.Contains(t, rr.Body.String(), ErrFailedToFetchMetadata.Error())
	})

	t.Run("unexpected-error-from-service", func(t *testing.T) {
		t.Parallel()

		mockProvider := &mockRevokeHandlerProvider{}
		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("unexpected service error")
			},
		}

		req, rr := buildRequest()

		handler := NewRevokeHandler(mockProvider, client, testhelpers.MockLogger)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "unexpected service error")
	})
}
