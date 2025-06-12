package utils

import (
	"testing"
)

func TestMask_ShowFullSecrets(t *testing.T) {
	t.Setenv("DEBUG_SHOW_SECRETS", "1")

	got := Mask("secret123")
	want := "secret123"
	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestMask_HiddenFull(t *testing.T) {
	t.Setenv("DEBUG_SHOW_SECRETS", "0")

	tests := []struct {
		input string
		want  string
	}{
		{"", "(empty)"},
		{"abc", "****"},
		{"abcd", "****"},
		{"secret123", "*****t123"},
	}

	for _, tt := range tests {
		got := Mask(tt.input)
		if got != tt.want {
			t.Errorf("Mask(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}
