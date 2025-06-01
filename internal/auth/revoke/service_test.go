package revoke

import (
	"bytes"
	"io"
	"net/http"
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
	assert.Contains(t, err.Error(), "failed to parse metadata")
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
