package config

import (
	"testing"

	"cognito-repeater-go/internal/config"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Setenv("AWS_REGION", "ap-northeast-1")
	t.Setenv("AWS_COGNITO_CLIENT_SECRET", "dummy-secret")
	t.Setenv("AWS_COGNITO_LOGOUT_URI", "https://example.com/logout")
	t.Setenv("AWS_COGNITO_REDIRECT_URI", "https://example.com/callback")
	t.Setenv("AWS_COGNITO_SCOPE", "openid")
	t.Setenv("AWS_COGNITO_USER_POOL_CLIENT_ID", "abc123")
	t.Setenv("AWS_COGNITO_USER_POOL_ID", "ap-northeast-1_xyz789")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_xyz789/.well-known/openid-configuration"
	if got := cfg.MetadataURL(); got != want {
		t.Errorf("unexpected metadata URL: got %s, want %s", got, want)
	}
}

func TestLoadConfig_MissingVars(t *testing.T) {
	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error due to missing env vars, got nil")
	}
}
