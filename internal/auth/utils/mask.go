package utils

import (
	"os"
	"strings"
)

// showSecretsEnabled checks whether secret masking should be disabled.
// It returns true if DEBUG_SHOW_SECRETS is set to "1", a common convention for boolean true in env vars.
func showSecretsEnabled() bool {
	return os.Getenv("DEBUG_SHOW_SECRETS") == "1"
}

// Mask returns a masked version of the input string for safe logging.
// If DEBUG_SHOW_SECRETS=1, the original string is returned for debugging.
// - If input is empty: returns "(empty)"
// - If length <= 4: returns "****"
// - Otherwise: masks all but the last 4 characters.
func Mask(s string) string {
	if s == "" {
		return "(empty)"
	}
	if showSecretsEnabled() {
		return s
	}
	if len(s) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}
