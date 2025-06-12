package revoke

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestGetRevokeURL_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"revocation_endpoint":"https://example.com/oauth2/revoke"}`))
	}))
	defer ts.Close()

	endpoint, err := GetRevokeURL(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.NoError(t, err, "failed to fetch revocation_endpoint")
	assert.Equal(t, "https://example.com/oauth2/revoke", endpoint)
}

func TestGetRevokeURL_Errors(t *testing.T) {
	t.Parallel()

	t.Run("failed to fetch metadata", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		_, err := GetRevokeURL("http://dummy", client, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				body := io.NopCloser(bytes.NewBufferString(`unexpected`))
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       body,
				}, nil
			},
		}

		_, err := GetRevokeURL("http://dummy", client, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrUnexpectedStatusCode)
	})

	t.Run("invalid json response", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				body := io.NopCloser(bytes.NewBufferString(`invalid-json`))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil
			},
		}

		_, err := GetRevokeURL("http://dummy", client, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
	})

	t.Run("missing revocation_endpoint", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				body := io.NopCloser(bytes.NewBufferString(`{}`))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil
			},
		}

		_, err := GetRevokeURL("http://dummy", client, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrMissingRevocationEndpoint)
	})
}

type mockSuccessClient struct{}

func (m *mockSuccessClient) Do(req *http.Request) (*http.Response, error) {
	bodyBytes, _ := io.ReadAll(req.Body)
	bodyStr := string(bodyBytes)
	if !strings.Contains(bodyStr, "token=mock-refresh-token") {
		panic("missing refresh token in body")
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		panic("wrong content-type")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("ok")),
	}, nil
}

func TestSendRevokeRequest_Success(t *testing.T) {
	t.Parallel()

	client := &mockSuccessClient{}

	resp, err := SendRevokeRequest(
		"https://example.com/oauth2/revoke",
		client,
		"mock-refresh-token",
		"mock-client-id",
		"mock-client-secret",
		testhelpers.MockLogger,
	)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendRevokeRequest_Errors(t *testing.T) {
	t.Parallel()

	validURL := "https://example.com/revoke"
	validToken := "refresh-token"
	validClientID := "client-id"
	validSecret := "client-secret"

	t.Run("missing client credentials", func(t *testing.T) {
		t.Parallel()

		_, err := SendRevokeRequest(validURL, testhelpers.NewMockHTTPClientPanic(), validToken, "", "", testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrMissingClientCredentials)
	})

	t.Run("failed to create revoke request", func(t *testing.T) {
		t.Parallel()

		brokenURL := "http://[::1]:NamedPort"
		_, err := SendRevokeRequest(brokenURL, testhelpers.NewMockHTTPClientPanic(), validToken, validClientID, validSecret, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrFailedToCreateRevokeRequest)
	})

	t.Run("failed to send revoke request", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		_, err := SendRevokeRequest(validURL, client, validToken, validClientID, validSecret, testhelpers.MockLogger)
		assert.ErrorIs(t, err, ErrFailedToSendRevokeRequest)
	})
}
