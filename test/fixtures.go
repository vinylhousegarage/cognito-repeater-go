package test

import "cognito-repeater-go/internal/config"

var mockCfg = &config.Config{
	Region:           "ap-northeast-1",
	ClientSecret:     "client-secret",
	LogoutURI:        "https://example.com/logout",
	RedirectURI:      "https://localhost/callback",
	Scope:            "openid",
	UserPoolClientID: "client-id",
	UserPoolID:       "pool-id",
}
