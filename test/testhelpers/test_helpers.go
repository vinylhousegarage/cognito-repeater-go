package testhelpers

import (
	"net/http"

	"cognito-repeater-go/internal/config"
)

var MockCfg = &config.Config{
	Region:           "ap-northeast-1",
	ClientSecret:     "client-secret",
	LogoutURI:        "https://example.com/logout",
	RedirectURI:      "https://localhost/callback",
	Scope:            "openid",
	UserPoolClientID: "client-id",
	UserPoolID:       "pool-id",
}

type MockMetadataURL struct{}

func (m *MockMetadataURL) MetadataURL() string {
	return "https://mock.metadata.url"
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
