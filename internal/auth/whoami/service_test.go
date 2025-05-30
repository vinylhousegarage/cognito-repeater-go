package whoami

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserinfoURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"userinfo_endpoint": "https://example.com/oauth2/userinfo"}`))
	}))
	defer ts.Close()

	endpoint, err := GetUserinfoURL(ts.URL, http.DefaultClient)

	assert.NoError(t, err, "failed to fetch userinfo_endpoint")
	assert.Equal(t, "https://example.com/oauth2/userinfo", endpoint)
}

func TestGetUserinfoURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetUserinfoURL(ts.URL, http.DefaultClient)

	assert.Error(t, err, "unexpected status code")
}

func TestFetchUserinfo_Success(t *testing.T) {
	t.Parallel()

	const userinfoJSON = `{
		"sub": "abc123",
		"email": "user@example.com",
		"email_verified": true
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(userinfoJSON))
	}))
	defer server.Close()

	client := &http.Client{}
	data, err := FetchUserinfo(client, server.URL, "dummy-token")

	assert.NoError(t, err)
	assert.Equal(t, "abc123", data["sub"])
	assert.Equal(t, "user@example.com", data["email"])
	assert.Equal(t, true, data["email_verified"])
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchUserinfo_HTTPError(t *testing.T) {
	t.Parallel()

	brokenClient := &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("simulated network error")
		}),
	}

	_, err := FetchUserinfo(brokenClient, "https://example.com/oauth2/userinfo", "dummy-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to call userinfo endpoint")
}

func TestFetchUserinfo_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"keys": [`))
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := FetchUserinfo(client, server.URL, "dummy-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse userinfo response")
}

func TestFetchUserinfo_StatusCodeNotOK(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer server.Close()

	client := &http.Client{}
	_, err := FetchUserinfo(client, server.URL, "dummy-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userinfo endpoint returned status")
}
