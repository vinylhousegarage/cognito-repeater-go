package me

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAccessToken(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{"valid header", "Bearer abc.def.ghi", "abc.def.ghi", false},
		{"missing header", "", "", true},
		{"invalid scheme", "Basic abc", "", true},
		{"no token", "Bearer", "", true},
	}

	for _, c := range cases {
		req := httptest.NewRequest("GET", "/me", nil)
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}

		token, err := extractAccessToken(req)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error but got none", c.name)
		}
		if !c.wantErr && token != c.wantToken {
			t.Errorf("%s: expected token %q but got %q", c.name, c.wantToken, token)
		}
	}
}

func TestGetJWKSURI_Success(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"jwks_uri":"https://example.com/dummy/.well-known/jwks.json"}`))
	}))
	defer ts.Close()

	endpoint, err := GetJWKSURI(http.DefaultClient, ts.URL)

	assert.NoError(t, err, "failed to fetch jwks_uri")
	assert.Equal(t, "https://example.com/dummy/.well-known/jwks.json", endpoint)
}
