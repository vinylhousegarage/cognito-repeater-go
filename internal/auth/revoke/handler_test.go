package revoke

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/response"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap"
)

const (
	mockMetadataURL = "https://mock.example.com/.well-known/openid-configuration"
	mockSecret      = "dummy-secret"
	mockClientID    = "dummy-client-id"
)

func TestNewRevokeHandler_Success(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
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

	mockProvider := testhelpers.NewMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	mockLogger := zap.NewNop()

	handler := NewRevokeHandler(mockProvider, mockClient, mockLogger)
	handler(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestNewRevokeHandler_MissingToken(t *testing.T) {
	t.Parallel()

	mockClient := testhelpers.NewMockHTTPClientPanic()
	req := httptest.NewRequest(http.MethodPost, "/revoke", nil)
	w := httptest.NewRecorder()

	mockProvider := testhelpers.NewMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	mockLogger := zap.NewNop()

	handler := NewRevokeHandler(mockProvider, mockClient, mockLogger)
	handler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestNewRevokeHandler_InvalidRevokeEndpointResponse(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := io.NopCloser(strings.NewReader("invalid-json"))
			return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
		},
	}

	form := strings.NewReader("token=dummy-token")
	req := httptest.NewRequest(http.MethodPost, "/revoke", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mockProvider := testhelpers.NewMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	mockLogger := zap.NewNop()

	handler := NewRevokeHandler(mockProvider, mockClient, mockLogger)
	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "failed to resolve revocation endpoint", errResp.Error)
}

func TestNewRevokeHandler_ServerErrorDuringRevoke(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			url := req.URL.String()
			if strings.Contains(url, "/.well-known/openid-configuration") {
				body := io.NopCloser(strings.NewReader(`{"revocation_endpoint":"https://mock-revoke-url"}`))
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
			}
			return &http.Response{StatusCode: http.StatusBadGateway, Status: "502 Bad Gateway", Body: io.NopCloser(strings.NewReader("revocation failed"))}, nil
		},
	}

	form := strings.NewReader("token=dummy-token")
	req := httptest.NewRequest(http.MethodPost, "/revoke", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mockProvider := testhelpers.NewMockRevokeHandlerProvider(
		mockMetadataURL,
		mockSecret,
		mockClientID,
	)

	mockLogger := zap.NewNop()

	handler := NewRevokeHandler(mockProvider, mockClient, mockLogger)
	handler(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "revocation failed with status")
}
