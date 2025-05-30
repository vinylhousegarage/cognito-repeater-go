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

type MockCallbackHandlerProvider struct{}
func (m MockCallbackHandlerProvider) ClientSecretValue() string     { return "mock-client-secret" }
func (m MockCallbackHandlerProvider) MetadataURL() string           { return "https://example.com/.well-known/openid-configuration" }
func (m MockCallbackHandlerProvider) RedirectURIValue() string      { return "https://example.com/callback" }
func (m MockCallbackHandlerProvider) UserPoolClientIDValue() string { return "mock-client-id" }

type MockLoginHandlerProvider struct{}
func (m MockLoginHandlerProvider) ClientSecretValue() string     { return "mock-client-secret" }
func (m MockLoginHandlerProvider) MetadataURL() string           { return "https://example.com/.well-known/openid-configuration" }
func (m MockLoginHandlerProvider) RedirectURIValue() string      { return "https://example.com/callback" }
func (m MockLoginHandlerProvider) UserPoolClientIDValue() string { return "mock-client-id" }

type MockLogoutHandlerProvider struct{}
func (m MockLogoutHandlerProvider) MetadataURL() string { return "https://example.com/.well-known/openid-configuration" }

type MockMeHandlerProvider struct{}
func (m MockMeHandlerProvider) Audience() string    { return "mock-audience" }
func (m MockMeHandlerProvider) GetJWKSURI() string  { return "https://example.com/client-id/.well-known/jwks.json" }
func (m MockMeHandlerProvider) Issuer() string      { return "mock-issuer" }
func (m MockMeHandlerProvider) MetadataURL() string { return "https://example.com/.well-known/openid-configuration" }
