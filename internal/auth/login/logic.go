package logic

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure random state: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
}
