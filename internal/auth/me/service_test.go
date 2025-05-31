package me

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
)

func TestGetJWKSURI_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"jwks_uri":"https://example.com/dummy/.well-known/jwks.json"}`))
	}))
	defer ts.Close()

	endpoint, err := GetJWKSURI(ts.URL, http.DefaultClient)

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

	_, err := GetJWKSURI("https://example.com", mockClient)
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
	_, err := GetJWKSURI(server.URL, client)
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
	_, err := GetJWKSURI(server.URL, client)
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
	_, err := GetJWKSURI(server.URL, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jwks_uri is empty")
}

func TestGetJWKSURI_InvalidURL(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	_, err := GetJWKSURI("http://[::1]:namedport", client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestFetchJWKSet_Success(t *testing.T) {
	t.Parallel()

	const jwksJSON = `{
		"keys": [
			{
				"kid": "example-kid",
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"n": "example-n",
				"e": "AQAB"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(jwksJSON))
	}))
	defer server.Close()

	client := &http.Client{}
	jwks, err := FetchJWKSet(server.URL, client)

	assert.NoError(t, err)
	assert.Len(t, jwks.Keys, 1)
	assert.Equal(t, "example-kid", jwks.Keys[0].Kid)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchJWKSet_HTTPError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}),
	}

	_, err := FetchJWKSet("https://example.com/jwks", brokenClient)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch jwks")
}

func TestFetchJWKSet_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys": [`))
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := FetchJWKSet(server.URL, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse jwks JSON")
}
