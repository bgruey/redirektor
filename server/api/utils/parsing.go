package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func HashFromUrl(prefix string, r *http.Request) (string, error) {
	parts := strings.Split(r.URL.String(), prefix)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid URL")
	}

	target := strings.Split(parts[len(parts)-1], "?")

	return strings.Replace(target[0], "/", "", -1), nil
}
