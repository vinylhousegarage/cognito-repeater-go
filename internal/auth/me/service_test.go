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

	endpoint, err := GetJWKSURI(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.NoError(t, err, "failed to fetch jwks_uri")
	assert.Equal(t, "https://example.com/dummy/.well-known/jwks.json", endpoint)
}

func TestGetJWKSURI_RequestCreationError(t *testing.T) {
	t.Parallel()

	_, err := GetJWKSURI(":", http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToCreateRequest)
}

func TestGetJWKSURI_HTTPClientError(t *testing.T) {
	t.Parallel()

	_, err := GetJWKSURI("http://invalid.host.local", http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToFetchMetadata)
}

func TestGetJWKSURI_StatusCodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetJWKSURI(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnexpectedMetadataStatusCode)
}

func TestGetJWKSURI_DecodeError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer ts.Close()

	_, err := GetJWKSURI(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToDecodeMetadata)
}

func TestGetJWKSURI_ErrMissingJWKSURI(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"end_session_endpoint": ""}`))
	}))
	defer ts.Close()

	_, err := GetJWKSURI(ts.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingJWKSURI)
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

	jwks, err := FetchJWKSet(server.URL, http.DefaultClient, testhelpers.MockLogger)

	assert.NoError(t, err)
	assert.Len(t, jwks.Keys, 1)
	assert.Equal(t, "example-kid", jwks.Keys[0].Kid)
}

func TestFetchJWKSet_RequestCreationError(t *testing.T) {
	t.Parallel()

	_, err := FetchJWKSet(":", http.DefaultClient, testhelpers.MockLogger)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToCreateJWKSetRequest)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchJWKSet_HTTPClientError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		}),
	}

	_, err := FetchJWKSet("https://example.com/jwks", brokenClient, testhelpers.MockLogger)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToFetchJWKS)
}

func TestFetchJWKSet_DecodeError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys": [`))
	}))
	defer server.Close()

	_, err := FetchJWKSet(server.URL, http.DefaultClient, testhelpers.MockLogger)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToDecodeJWKSet)
}
