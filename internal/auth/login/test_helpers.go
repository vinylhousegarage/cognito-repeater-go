package login

import (
	"net/http"
	"net/http/httptest"

	"cognito-repeater-go/test/testhelpers"
)

type mockLoginHandlerProvider struct {
	ClientID    string
	RedirectURI string
	Scope       string
	Metadata    string
}

func (m *mockLoginHandlerProvider) MetadataURL() string           { return m.Metadata }
func (m *mockLoginHandlerProvider) UserPoolClientIDValue() string { return m.ClientID }
func (m *mockLoginHandlerProvider) RedirectURIValue() string      { return m.RedirectURI }
func (m *mockLoginHandlerProvider) ScopeValue() string            { return m.Scope }

func newMockLoginHandlerProvider() *mockLoginHandlerProvider {
	return &mockLoginHandlerProvider{
		ClientID:    "example-client-id",
		RedirectURI: "https://example.com/callback",
		Scope:       "openid",
		Metadata:    "https://mock.example.com/.well-known/openid-configuration",
	}
}

func newMockHTTPClient() *testhelpers.MockHTTPClient {
	return &testhelpers.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.Write([]byte(`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`))
			return rec.Result(), nil
		},
	}
}
