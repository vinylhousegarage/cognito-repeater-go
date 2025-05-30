package testhelpers

import (
	"net/http"

	"cognito-repeater-go/internal/config"
)

var MockCfg config.CognitoMetadataProvider = &config.Config{
	Region:           "ap-northeast-1",
	ClientSecret:     "client-secret",
	LogoutURI:        "https://example.com/logout",
	RedirectURI:      "https://localhost/callback",
	Scope:            "openid",
	UserPoolClientID: "abc123clientidxyz4567890123",
	UserPoolID:       "ap-northeast-1_Abc123XYZ",
}

type MockMetadataURL struct { URL string }
func (m *MockMetadataURL) MetadataURL() string { return m.URL }

type MockHTTPClient struct { DoFunc func(req *http.Request) (*http.Response, error) }
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) { return m.DoFunc(req) }

func NewMockHTTPClientOK() httpclient.HTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		},
	}
}

func NewMockHTTPClientPanic() httpclient.HTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			panic("unexpected HTTP call")
		},
	}
}

type MockCallbackProvider struct{}
func (m MockCallbackProvider) ClientSecretValue() string     { return "mock-client-secret" }
func (m MockCallbackProvider) MetadataURL() string           { return "https://example.com/.well-known/openid-configuration" }
func (m MockCallbackProvider) RedirectURIValue() string      { return "https://example.com/callback" }
func (m MockCallbackProvider) UserPoolClientIDValue() string { return "mock-client-id" }

type MockLoginProvider struct{}
func (m MockLoginProvider) ClientSecretValue() string     { return "mock-client-secret" }
func (m MockLoginProvider) MetadataURL() string           { return "https://example.com/.well-known/openid-configuration" }
func (m MockLoginProvider) RedirectURIValue() string      { return "https://example.com/callback" }
func (m MockLoginProvider) UserPoolClientIDValue() string { return "mock-client-id" }

type MockLogoutProvider struct{}
func (m MockLogoutProvider) MetadataURL() string { return "https://example.com/.well-known/openid-configuration" }

type MockMeProvider struct{}
func (m MockMeProvider) Audience() string    { return "mock-audience" }
func (m MockMeProvider) GetJWKSURI() string  { return "https://example.com/client-id/.well-known/jwks.json" }
func (m MockMeProvider) Issuer() string      { return "mock-issuer" }
func (m MockMeProvider) MetadataURL() string { return "https://example.com/.well-known/openid-configuration" }
