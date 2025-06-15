package openapi

import (
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOpenAPISpec_Valid(t *testing.T) {
	doc, err := LoadOpenAPISpec("../../testdata/openapi_valid.yaml", testhelpers.MockLogger)
	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Cognito Repeater API (Go)", doc.Info.Title)
}

func TestLoadOpenAPISpec_FileNotFound(t *testing.T) {
	doc, err := LoadOpenAPISpec("no_such_file.yaml", testhelpers.MockLogger)
	require.Error(t, err)
	assert.Nil(t, doc)
}

func TestLoadOpenAPISpec_InvalidYAML(t *testing.T) {
	doc, err := LoadOpenAPISpec("../../testdata/openapi_invalid.yaml", testhelpers.MockLogger)
	require.Error(t, err)
	assert.Nil(t, doc)
}
