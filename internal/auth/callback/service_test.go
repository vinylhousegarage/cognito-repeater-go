package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestValidateCallbackRequestValidInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=xyz789", nil)

	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "xyz789",
	})

	code, err := ValidateCallbackRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, "abc123", code)
}

func TestValidateCallbackRequestMissingCode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)

	assert.EqualError(t, err, "missing code")
}

func TestValidateCallbackRequestMissingState(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)

	assert.EqualError(t, err, "missing state")
}

func TestValidateCallbackRequestMissingCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)

	_, err := ValidateCallbackRequest(req)

	assert.EqualError(t, err, "missing oauth_state cookie")
}

func TestValidateCallbackRequestStateMismatch(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "wrong"})

	_, err := ValidateCallbackRequest(req)

	assert.EqualError(t, err, "invalid state")
}

func TestGetCallbackURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token_endpoint":"https://example.com/oauth2/token"}`))
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	endpoint, err := GetCallbackURL(ts.URL, http.DefaultClient, mockLogger)

	assert.NoError(t, err, "failed to fetch token endpoint")
	assert.Equal(t, "https://example.com/oauth2/token", endpoint)
}

func TestGetCallbackURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, mockLogger)

	assert.EqualError(t, err, "unexpected status code: 500")
}

func TestGetCallbackURLMalformedJSON(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token_endpoint":`))
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err, "expected JSON decode error")
}

func TestGetCallbackURLMissingTokenEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issuer": "https://example.com"}`))
	}))
	defer ts.Close()

	mockLogger := zap.NewNop()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient, mockLogger)

	assert.Error(t, err, "expected missing token_endpoint error")
}
