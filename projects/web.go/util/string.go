package util

import (
	"bytes"
	"encoding/base64"
	"strings"
)

// Base64 returns the contents of b as a base64-encoded string.
func Base64(b []byte) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write(b)
	encoder.Close()
	return buf.String()
}

func ParseUrl(url string) []string {
	vs := strings.Split(url, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item == "" {
			continue
		}
		parts = append(parts, item)
	}
	return parts
}
