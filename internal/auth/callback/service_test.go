package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCallbackRequestValidInput(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc123&state=xyz789", nil)

	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "xyz789",
	})

	code, err := ValidateCallbackRequest(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedCode := "abc123"
	if code != expectedCode {
		t.Errorf("expected code %q, got %q", expectedCode, code)
	}
}

func TestValidateCallbackRequestMissingCode(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)
	if err == nil || err.Error() != "missing code" {
		t.Errorf("expected error 'missing code', got %v", err)
	}
}

func TestValidateCallbackRequestMissingState(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)
	if err == nil || err.Error() != "missing state" {
		t.Errorf("expected error 'missing state', got %v", err)
	}
}

func TestValidateCallbackRequestMissingCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)

	_, err := ValidateCallbackRequest(req)
	if err == nil || err.Error() != "missing oauth_state cookie" {
		t.Errorf("expected error 'missing oauth_state cookie', got %v", err)
	}
}

func TestValidateCallbackRequestStateMismatch(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "wrong"})

	_, err := ValidateCallbackRequest(req)
	if err == nil || err.Error() != "invalid state" {
		t.Errorf("expected error 'invalid state', got %v", err)
	}
}

func TestGetCallbackURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token_endpoint":"https://example.com/oauth2/token"}`))
	}))
	defer ts.Close()

	endpoint, err := GetCallbackURL(ts.URL, http.DefaultClient)

	assert.NoError(t, err, "failed to fetch token endpoint")
	assert.Equal(t, "https://example.com/oauth2/token", endpoint)
}

func TestGetCallbackURLStatusCode500(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient)

	assert.Error(t, err, "unexpected status code")
}

func TestGetCallbackURLMalformedJSON(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token_endpoint":`))
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient)

	assert.Error(t, err, "expected JSON decode error")
}

func TestGetCallbackURLMissingTokenEndpoint(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"issuer": "https://example.com"}`))
	}))
	defer ts.Close()

	_, err := GetCallbackURL(ts.URL, http.DefaultClient)

	assert.Error(t, err, "expected missing token_endpoint error")
}
