package domain

import (
	"crypto/sha1"
	"encoding/base64"
)

func Shorten(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	b := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
