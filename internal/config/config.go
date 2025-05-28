package config

import (
	"fmt"
	"os"
)

type CognitoMetadataProvider interface {
	Audience() string
	ClientSecretValue() string
	Issuer() string
	MetadataURL() string
	RedirectURIValue() string
	ScopeValue() string
	UserPoolClientIDValue() string
}

type Config struct {
	Region           string
	ClientSecret     string
	LogoutURI        string
	RedirectURI      string
	Scope            string
	UserPoolClientID string
	UserPoolID       string
}

func (c *Config) Audience() string {
	return c.UserPoolClientID
}

func (c *Config) ClientSecretValue() string {
	return c.ClientSecret
}

func (c *Config) Issuer() string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", c.Region, c.UserPoolID)
}

func (c *Config) MetadataURL() string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/openid-configuration", c.Region, c.UserPoolID)
}

type MetadataURLProvider interface {
	MetadataURL() string
}

func (c *Config) RedirectURIValue() string {
	return c.RedirectURI
}

func (c *Config) ScopeValue() string {
	return c.Scope
}

func (c *Config) UserPoolClientIDValue() string {
	return c.UserPoolClientID
}

func LoadConfig() (*Config, error) {
	required := []string{
		"AWS_REGION",
		"AWS_COGNITO_CLIENT_SECRET",
		"AWS_COGNITO_LOGOUT_URI",
		"AWS_COGNITO_REDIRECT_URI",
		"AWS_COGNITO_SCOPE",
		"AWS_COGNITO_USER_POOL_CLIENT_ID",
		"AWS_COGNITO_USER_POOL_ID",
	}

	missing := []string{}
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return &Config{
		Region:           os.Getenv("AWS_REGION"),
		ClientSecret:     os.Getenv("AWS_COGNITO_CLIENT_SECRET"),
		LogoutURI:        os.Getenv("AWS_COGNITO_LOGOUT_URI"),
		RedirectURI:      os.Getenv("AWS_COGNITO_REDIRECT_URI"),
		Scope:            os.Getenv("AWS_COGNITO_SCOPE"),
		UserPoolClientID: os.Getenv("AWS_COGNITO_USER_POOL_CLIENT_ID"),
		UserPoolID:       os.Getenv("AWS_COGNITO_USER_POOL_ID"),
	}, nil
}
