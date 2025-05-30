package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"cognito-repeater-go/internal/auth/deps"
)

func GenerateState() string {
	const stateLength = 16
	b := make([]byte, stateLength)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure random state: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}

func BuildStateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func BuildLoginURL(p deps.LoginHandlerProvider, endpoint, state string) (string, error) {
	loginURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	q := loginURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", p.UserPoolClientIDValue())
	q.Set("redirect_uri", p.RedirectURIValue())
	q.Set("scope", p.ScopeValue())
	q.Set("state", state)
	loginURL.RawQuery = q.Encode()
	return loginURL.String(), nil
}
