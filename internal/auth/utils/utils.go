package utils

import (
	"fmt"
	"net/http"
)

func ExtractFormValue(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", fmt.Errorf("failed to parse form: %w", err)
	}

	token := r.FormValue("token")
	if token == "" {
		return "", fmt.Errorf("token is missing from body")
	}

	return token, nil
}
