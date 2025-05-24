package me

import (
	"net/http/httptest"
	"testing"
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
