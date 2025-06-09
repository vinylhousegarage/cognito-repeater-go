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

	assert.ErrorIs(t, err, ErrMissingCode)
}

func TestValidateCallbackRequestMissingState(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrMissingState)
}

func TestValidateCallbackRequestMissingCookie(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrMissingStateCookie)
}

func TestValidateCallbackRequestStateMismatch(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/callback?code=abc&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "wrong"})

	_, err := ValidateCallbackRequest(req)

	assert.ErrorIs(t, err, ErrInvalidState)
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

