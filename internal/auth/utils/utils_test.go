package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestExtractFormValue(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		body := "token=abc.def.ghi"
		req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		token, err := ExtractFormValue(req, zap.NewNop())
		assert.NoError(t, err)
		assert.Equal(t, "abc.def.ghi", token)
	})

	t.Run("failed to parse form", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader("%"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err := ExtractFormValue(req, zap.NewNop())
		assert.ErrorIs(t, err, ErrFailedToParseForm)
	})

	t.Run("missing token", func(t *testing.T) {
		t.Parallel()

		body := ""
		req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err := ExtractFormValue(req, zap.NewNop())
		assert.ErrorIs(t, err, ErrMissingToken)
	})

	t.Run("bad form encoding", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader("{invalid-json}"))
		req.Header.Set("Content-Type", "application/json")

		_, err := ExtractFormValue(req, zap.NewNop())
		assert.ErrorIs(t, err, ErrMissingToken)
	})
}
