package config

import (
	"fmt"
	"os"
)

type Config struct {
	Region           string
	ClientSecret     string
	LogoutURI        string
	RedirectURI      string
	Scope            string
	UserPoolClientID string
	UserPoolID       string
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

func (c *Config) MetadataURL() string {
	return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/openid-configuration", c.Region, c.UserPoolID)
}

type MetadataURLProvider interface {
	MetadataURL() string
}
