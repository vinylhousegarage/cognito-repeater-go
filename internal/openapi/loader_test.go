package openapi

import (
	"testing"

	"cognito-repeater-go/test/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadEmbeddedOpenAPISpec(t *testing.T) {
	t.Parallel()

	doc, err := LoadEmbeddedOpenAPISpec(testhelpers.MockLogger)
	require.NoError(t, err)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.OpenAPI)
	assert.Equal(t, "Cognito Repeater API (Go)", doc.Info.Title)
}
