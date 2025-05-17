package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCognitoLogoutURL_ReturnExpected(t *testing.T) {
	expected := "https://stab.com/logout"
	actual := parseCognitoLogoutURL()
	assert.Equal(t, expected, actual)
}
