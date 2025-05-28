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
	form.Set("client_id", cfg.UserPoolClientID)
	form.Set("redirect_uri", cfg.RedirectURI)

	return form.Encode()
}

func BuildBasicAuthHeader(cfg config.CognitoMetadataProvider) string {
	raw := cfg.UserPoolClientID + ":" + cfg.ClientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}
