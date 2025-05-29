package callback

import (
	"encoding/base64"
	"net/url"
)

type CallbackHandlerProvider interface {
	ClientSecretValue() string
	MetadataURL() string
	RedirectURIValue() string
	UserPoolClientIDValue() string
}

func BuildTokenRequestBody(code string, p CallbackHandlerProvider) string {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", p.UserPoolClientIDValue())
	form.Set("redirect_uri", p.RedirectURIValue())

	return form.Encode()
}

func BuildBasicAuthHeader(p CallbackHandlerProvider) string {
	raw := p.UserPoolClientIDValue() + ":" + p.ClientSecretValue()
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}
