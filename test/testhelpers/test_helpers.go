package testhelpers

import (
	"io"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"

	"go.uber.org/zap"
)

var MockCfg = config.Config{
	Region:           "ap-northeast-1",
	ClientSecret:     "client-secret",
	LogoutURI:        "https://example.com/logout",
	RedirectURI:      "https://localhost/callback",
	Scope:            "openid",
	UserPoolClientID: "test-client",
	UserPoolID:       "test-pool",
}

type MockAllHandlerProviders struct{ *config.Config }

func NewMockRouteDependencies() deps.RouteDependencies {
	return deps.RouteDependencies{
		CallbackProvider: &MockCfg,
		LoginProvider:    &MockCfg,
		LogoutProvider:   &MockCfg,
		MeProvider:       &MockCfg,
		RevokeProvider:   &MockCfg,
		WhoamiProvider:   &MockCfg,
	}
}

type MockMetadataURL struct{ URL string }

func (m *MockMetadataURL) MetadataURL() string { return m.URL }

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

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

type MockRevokeHandlerProvider struct {
	ClientSecretStr     string
	MetadataURLStr      string
	UserPoolClientIDStr string
}

func (m *MockRevokeHandlerProvider) ClientSecretValue() string {
	return m.ClientSecretStr
}

func (m *MockRevokeHandlerProvider) MetadataURL() string {
	return m.MetadataURLStr
}

func (m *MockRevokeHandlerProvider) UserPoolClientIDValue() string {
	return m.UserPoolClientIDStr
}

func NewMockRevokeHandlerProvider(metadataURL, clientSecret, clientID string) deps.RevokeHandlerProvider {
	return &MockRevokeHandlerProvider{
		MetadataURLStr:      metadataURL,
		ClientSecretStr:     clientSecret,
		UserPoolClientIDStr: clientID,
	}
}

func NewMockHTTPClientOKWithAuthEndpoint() httpclient.HTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(
					`{"authorization_endpoint": "https://example.com/oauth2/authorize"}`,
				)),
			}, nil
		},
	}
}

var MockProvider = &MockCfg
var MockClient = NewMockHTTPClientOKWithAuthEndpoint()
var MockLogger = zap.NewNop()

type ReadErrorCloser struct{}

func (r *ReadErrorCloser) Read(p []byte) (int, error) { return 0, errors.New("read error") }
func (r *ReadErrorCloser) Close() error               { return nil }
