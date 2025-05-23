package callback

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cognito-repeater-go/test/testhelpers"

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

var _ callbackURLProvider = (*callbackClient)(nil)

func TestGetCallbackURLReturnsExpectedEndpoint(t *testing.T) {
	t.Parallel()

	const responseBody = `{"token_endpoint":"https://example.com/oauth2/token"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}))
	defer ts.Close()

	mock := &testhelpers.MockMetadataURL{URL: ts.URL}

	mockClient := &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := responseBody
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := NewCallbackClient(mockClient)
	endpoint, err := client.GetCallbackURL(mock)

	assert.NoError(t, err, "failed to fetch token endpoint")
	assert.Equal(t, "https://example.com/oauth2/token", endpoint)
}
