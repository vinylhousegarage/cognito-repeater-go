package test_helpers

import (
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

type MockEndpointProvider struct{}

func (m *MockEndpointProvider) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/logout", nil
}

type MockMetadataProvider struct{}

func (m *MockMetadataProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}
