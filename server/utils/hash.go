package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

func Sha256Base64(password string) string {
	h := sha256.Sum256([]byte(password))
	ret := base64.StdEncoding.EncodeToString(h[:])
	ret = strings.ReplaceAll(ret, "/", "-")
	ret = strings.ReplaceAll(ret, "+", "_")
	ret = strings.ReplaceAll(ret, "=", ".")
	return ret
}
