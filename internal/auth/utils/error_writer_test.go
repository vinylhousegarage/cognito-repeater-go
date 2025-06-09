package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap/zaptest"
)

func TestWritePlainError(t *testing.T) {
	rec := httptest.NewRecorder()
	logger := zaptest.NewLogger(t)

	WritePlainError(rec, http.StatusInternalServerError, errors.New("something went wrong"), logger)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
}
