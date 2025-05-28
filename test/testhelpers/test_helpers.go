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

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
