package openapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOpenAPISpec_Valid(t *testing.T) {
	doc := LoadOpenAPISpec("../../testdata/openapi_valid.yaml")
	require.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Cognito Repeater API (Go)", doc.Info.Title)
}

func TestLoadOpenAPISpec_FileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for missing file")
		}
	}()
	LoadOpenAPISpec("no_such_file.yaml")
}

func TestLoadOpenAPISpec_InvalidYAML(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid YAML")
		}
	}()
	LoadOpenAPISpec("../../testdata/openapi_invalid.yaml")
}
