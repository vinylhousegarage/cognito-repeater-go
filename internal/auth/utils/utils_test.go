package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFormValueFromBody_Success(t *testing.T) {
	body := "id_token=abc.def.ghi"
	req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	token, err := ExtractFormValue(req)

	assert.NoError(t, err)
	assert.Equal(t, "abc.def.ghi", token)
}

func TestExtractFormValue_MissingToken(t *testing.T) {
	body := "" // id_token がない
	req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := ExtractFormValuey(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id_token is missing")
}

func TestExtractFormValue_BadFormEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader("{invalid-json}"))
	req.Header.Set("Content-Type", "application/json") // 違う形式

	_, err := ExtractFormValue(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse form")
}
