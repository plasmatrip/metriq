package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func (h Handlers) Sum(resp []byte) (string, error) {
	hm := hmac.New(sha256.New, []byte(h.config.Key))
	hm.Write(resp)
	dst := hm.Sum(nil)
	hash := base64.StdEncoding.EncodeToString(dst)
	return hash, nil
}
