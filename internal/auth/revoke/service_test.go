package revoke

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestGetRevokeURLSuccess(t *testing.T) {
	t.Parallel()

	body := `{"revocation_endpoint": "https://example.com/revoke"}`
	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}

	url, err := GetRevokeURL("https://mock-metadata-url", mockClient)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/revoke", url)
}

func TestGetRevokeURLNon200Status(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	_, err := GetRevokeURL("https://mock-url", mockClient)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestGetRevokeURLInvalidJSON(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("not-json")),
			}, nil
		},
	}

	_, err := GetRevokeURL("https://mock-url", mockClient)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode metadata")
}

func TestGetRevokeURLEmptyEndpoint(t *testing.T) {
	t.Parallel()

	body := `{"revocation_endpoint": ""}`
	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}

	_, err := GetRevokeURL("https://mock-url", mockClient)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "revocation_endpoint is missing")
}

func TestSendRevokeRequest_Success(t *testing.T) {
	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

			bodyBytes, err := io.ReadAll(req.Body)
			assert.NoError(t, err)
			bodyStr := string(bodyBytes)

			assert.Contains(t, bodyStr, "token=mock-refresh-token")
			assert.Contains(t, bodyStr, "token_type_hint=refresh_token")

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	revokeURL := "https://example.com/oauth2/revoke"
	refreshToken := "mock-refresh-token"

	resp, err := SendRevokeRequest(revokeURL, mockClient, refreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
