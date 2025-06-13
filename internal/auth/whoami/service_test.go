package whoami

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestGetUserinfoURL_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"userinfo_endpoint":"https://example.com/oauth2/userInfo"}`))
	}))
	defer ts.Close()

	logger := zaptest.NewLogger(t)
	endpoint, err := GetUserinfoURL(ts.URL, http.DefaultClient, logger)

	assert.NoError(t, err, "failed to fetch userinfo_endpoint")
	assert.Equal(t, "https://example.com/oauth2/userInfo", endpoint)
}

func TestGetUserinfoURL_Errors(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	t.Run("failed-to-fetch-metadata", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		_, err := GetUserinfoURL("http://dummy", client, logger)
		assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
	})

	t.Run("failed-to-read-body", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       &testhelpers.ReadErrorCloser{},
				}, nil
			},
		}

		_, err := GetUserinfoURL("http://dummy", client, logger)
		assert.ErrorIs(t, err, ErrFailedToReadMetadataResponse)
	})

	t.Run("unexpected-status-code", func(t *testing.T) {
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

		_, err := GetUserinfoURL("http://dummy", client, logger)
		assert.ErrorIs(t, err, ErrUnexpectedMetadataStatusCode)
	})

	t.Run("invalid-json-response", func(t *testing.T) {
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

		_, err := GetUserinfoURL("http://dummy", client, logger)
		assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
	})

	t.Run("missing-userinfo_endpoint", func(t *testing.T) {
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

		_, err := GetUserinfoURL("http://dummy", client, logger)
		assert.ErrorIs(t, err, ErrMissingUserinfoEndpoint)
	})
}

func TestFetchUserinfo_Success(t *testing.T) {
	t.Parallel()

	const jsonBody = `{
		"sub": "abc-123-def-456"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		assert.True(t, strings.HasPrefix(auth, "Bearer "), "Authorization header must be set")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(jsonBody))
	}))
	defer server.Close()

	logger := zaptest.NewLogger(t)
	userInfo, err := FetchUserinfo(server.URL, httpclient.HTTPClient(http.DefaultClient), "dummy-token", logger)

	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "abc-123-def-456", userInfo.Sub)
}

func TestFetchUserinfo_Errors(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	t.Run("failed-to-create-request", func(t *testing.T) {
		t.Parallel()

		client := httpclient.HTTPClient(http.DefaultClient)
		_, err := FetchUserinfo("://invalid-url", client, "token", logger)

		assert.ErrorIs(t, err, ErrFailedToCreateUserinfoRequest)
	})

	t.Run("http-client-do-error", func(t *testing.T) {
		t.Parallel()

		brokenClient := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		_, err := FetchUserinfo("https://example.com", brokenClient, "token", logger)
		assert.ErrorIs(t, err, ErrFailedToFetchUserinfoRequest)
	})

	t.Run("read-body-error", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       &testhelpers.ReadErrorCloser{},
				}, nil
			},
		}

		_, err := FetchUserinfo("https://example.com", client, "token", logger)
		assert.ErrorIs(t, err, ErrFailedToReadUserinfoResponse)
	})

	t.Run("unexpected-userinfo-response", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader(`access denied`)),
				}, nil
			},
		}

		_, err := FetchUserinfo("https://example.com", client, "token", logger)
		assert.ErrorIs(t, err, ErrUnexpectedUserinfoStatusCode)
	})

	t.Run("invalid-json", func(t *testing.T) {
		t.Parallel()

		client := &testhelpers.MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{invalid}`)),
				}, nil
			},
		}

		_, err := FetchUserinfo("https://example.com", client, "token", logger)
		assert.ErrorIs(t, err, ErrFailedToParseUserinfoResponse)
	})
}
