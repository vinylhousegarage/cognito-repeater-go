package utils

import (
	"fmt"
	"net/http"
)

func ExtractFormValue(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", fmt.Errorf("failed to parse form: %w", err)
	}

	idToken := r.FormValue("id_token")
	if idToken == "" {
		return "", fmt.Errorf("id_token is missing from body")
	}

	return idToken, nil
}
