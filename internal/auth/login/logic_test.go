package login

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	const state = "abc123"
	c := BuildStateCookie(state)

	assert.Equal(t, "oauth_state", c.Name)
	assert.Equal(t, state, c.Value)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Secure)
	assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
}
