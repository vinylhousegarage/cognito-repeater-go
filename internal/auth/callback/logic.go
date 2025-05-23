package callback

import (
	"net/url"

	"cognito-repeater-go/internal/config"
)

func BuildTokenRequestBody(code string, cfg *config.Config) string {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", cfg.UserPoolClientID)
	form.Set("redirect_uri", cfg.RedirectURI)

	return form.Encode()
}
