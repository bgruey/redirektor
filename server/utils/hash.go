package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func Sha256Base64(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}
