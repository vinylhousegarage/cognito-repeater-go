package callback

import (
	"encoding/base64"
	"net/url"

	"cognito-repeater-go/internal/auth/deps"
)

func BuildTokenRequestBody(code string, p deps.CallbackHandlerProvider) string {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", p.UserPoolClientIDValue())
	form.Set("redirect_uri", p.RedirectURIValue())

	return form.Encode()
}

func BuildBasicAuthHeader(p deps.CallbackHandlerProvider) string {
	raw := p.UserPoolClientIDValue() + ":" + p.ClientSecretValue()
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}
