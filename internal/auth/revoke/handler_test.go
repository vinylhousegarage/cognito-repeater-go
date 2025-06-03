package revoke

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

type mockRevokeHandlerProvider struct {
	ClientSecretStr     string
	MetadataURLStr      string
	UserPoolClientIDStr string
}

func (m *mockRevokeHandlerProvider) ClientSecretValue() string {
	return m.ClientSecretStr
}

func (m *mockRevokeHandlerProvider) MetadataURL() string {
	return m.MetadataURLStr
}

func (m *mockRevokeHandlerProvider) UserPoolClientIDValue() string {
	return m.UserPoolClientIDStr
}

func newMockRevokeHandlerProvider(metadataURL, clientSecret, clientID string) deps.RevokeHandlerProvider {
	return &mockRevokeHandlerProvider{
		MetadataURLStr:      metadataURL,
		ClientSecretStr:     clientSecret,
		UserPoolClientIDStr: clientID,
	}
}

const (
	mockMetadataURL = "https://mock.example.com/.well-known/openid-configuration"
	mockSecret      = "dummy-secret"
	mockClientID    = "dummy-client-id"
)

func TestNewRevokeHandler_Success(t *testing.T) {
	t.Parallel()

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			url := req.URL.String()
			switch {
			case strings.Contains(url, "/.well-known/openid-configuration"):
				body := io.NopCloser(strings.NewReader(`{"revocation_endpoint":"https://mock-revoke-url"}`))
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil

			case strings.Contains(url, "https://mock-revoke-url"):
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil

			default:
				return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("not found"))}, nil
			}
		},
	}

	form := strings.NewReader("token=dummy-token")
	req := httptest.NewRequest(http.MethodPost, "/revoke", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mockProvider := newMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	handler := NewRevokeHandler(mockProvider, client)
	handler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestNewRevokeHandler_MissingToken(t *testing.T) {
	t.Parallel()

	client := testhelpers.NewMockHTTPClientPanic()
	req := httptest.NewRequest(http.MethodPost, "/revoke", nil)
	w := httptest.NewRecorder()

	mockProvider := newMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	handler := NewRevokeHandler(mockProvider, client)
	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestNewRevokeHandler_RevokeURLFailure(t *testing.T) {
	t.Parallel()

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(strings.NewReader("invalid-json"))
			return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
		},
	}

	form := strings.NewReader("token=dummy-token")
	req := httptest.NewRequest(http.MethodPost, "/revoke", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mockProvider := newMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	handler := NewRevokeHandler(mockProvider, client)
	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get revoke endpoint")
}

func TestNewRevokeHandler_RevokeFailsWith500(t *testing.T) {
	t.Parallel()

	client := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			url := req.URL.String()
			if strings.Contains(url, "/.well-known/openid-configuration") {
				body := io.NopCloser(strings.NewReader(`{"revocation_endpoint":"https://mock-revoke-url"}`))
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
			}
			return &http.Response{StatusCode: http.StatusInternalServerError, Status: "500 Internal Server Error", Body: io.NopCloser(strings.NewReader(""))}, nil
		},
	}

	form := strings.NewReader("token=dummy-token")
	req := httptest.NewRequest(http.MethodPost, "/revoke", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mockProvider := newMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	handler := NewRevokeHandler(mockProvider, client)
	handler(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "revocation failed with status")
}
