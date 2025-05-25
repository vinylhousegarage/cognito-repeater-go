package me

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestExtractAccessToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{"valid header", "Bearer abc.def.ghi", "abc.def.ghi", false},
		{"missing header", "", "", true},
		{"invalid scheme", "Basic abc", "", true},
		{"no token", "Bearer", "", true},
	}

	for _, c := range cases {
		req := httptest.NewRequest("GET", "/me", nil)
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}

		token, err := extractAccessToken(req)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error but got none", c.name)
		}
		if !c.wantErr && token != c.wantToken {
			t.Errorf("%s: expected token %q but got %q", c.name, c.wantToken, token)
		}
	}
}

func TestGetJWKSURI_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"jwks_uri":"https://example.com/dummy/.well-known/jwks.json"}`))
	}))
	defer ts.Close()

	endpoint, err := GetJWKSURI(http.DefaultClient, ts.URL)

	assert.NoError(t, err, "failed to fetch jwks_uri")
	assert.Equal(t, "https://example.com/dummy/.well-known/jwks.json", endpoint)
}

func TestGetJWKSURI_HTTPClientError(t *testing.T) {
	t.Parallel()

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("simulated network failure")
		},
	}

	_, err := GetJWKSURI(mockClient, "https://example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch metadata")
}

func TestGetJWKSURI_StatusCodeError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := GetJWKSURI(client, server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestGetJWKSURI_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{ "jwks_uri": `))
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := GetJWKSURI(client, server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode metadata")
}

func TestGetJWKSURI_EmptyJWKSURI(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{ "jwks_uri": "" }`))
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := GetJWKSURI(client, server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwks_uri is empty")
}

func TestGetJWKSURI_InvalidURL(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	_, err := GetJWKSURI(client, "http://[::1]:namedport")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}
