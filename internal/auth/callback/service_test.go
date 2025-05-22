package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
