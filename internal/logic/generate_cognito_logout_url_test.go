package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCognitoLogoutUrl_ReturnsExpected(t *testing.T) {
	expected := "https://stab.com/logout"
	actual := generateCognitoLogoutUrl()
	assert.Equal(t, expected, actual)
}
