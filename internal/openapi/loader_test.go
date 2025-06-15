package openapi

import (
	"os"
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "openapi_valid_*.yaml")
	require.NoError(t, err)

	var n int
	n, err = tmpFile.WriteString(content)
	t.Log("byte:", n)
	require.NoError(t, err)

	err = tmpFile.Close()
	require.NoError(t, err)

	return tmpFile.Name()
}

func TestLoadOpenAPISpec_Valid(t *testing.T) {
	t.Parallel()

	content := `
openapi: 3.0.0
info:
  title: Cognito Repeater API (Go)
  version: "1.0.0"
paths: {}
`
	tmpFile := writeTempFile(t, content)

	doc, err := LoadOpenAPISpec(tmpFile, testhelpers.MockLogger)
	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Cognito Repeater API (Go)", doc.Info.Title)
}

func TestLoadOpenAPISpec_FileNotFound(t *testing.T) {
	t.Parallel()

	doc, err := LoadOpenAPISpec("no_such_file.yaml", testhelpers.MockLogger)
	require.Error(t, err)
	assert.Nil(t, doc)
}

func TestLoadOpenAPISpec_InvalidYAML(t *testing.T) {
	t.Parallel()

	doc, err := LoadOpenAPISpec("../../testdata/openapi_invalid.yaml", testhelpers.MockLogger)
	require.Error(t, err)
	assert.Nil(t, doc)
}
