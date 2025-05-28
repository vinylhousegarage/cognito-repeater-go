package callback

import (
	"encoding/base64"
	"net/url"

	"cognito-repeater-go/internal/config"
)

func BuildTokenRequestBody(code string, cfg config.CognitoMetadataProvider) string {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", cfg.UserPoolClientIDValue())
	form.Set("redirect_uri", cfg.RedirectURIValue())

	return form.Encode()
}

func BuildBasicAuthHeader(cfg config.CognitoMetadataProvider) string {
	raw := cfg.UserPoolClientIDValue() + ":" + cfg.ClientSecretValue()
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}
