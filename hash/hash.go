package hash

import (
	"hash"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type HMAC struct {
	hmac hash.Hash
}

func NewHMAC(key string) *HMAC {
	return &HMAC{
		hmac: hmac.New(sha256.New, []byte(key)),
	}
}

func (h *HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b) 
}