package callback

import (
	"errors"
	"net/http"
)

func ValidateCallbackRequest(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		return "", errors.New("missing code")
	}

	if state == "" {
		return "", errors.New("missing state")
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return "", errors.New("missing oauth_state cookie")
	}

	if state != cookie.Value {
		return "", errors.New("invalid state")
	}

	return code, nil
}
