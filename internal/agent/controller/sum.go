package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func (c Controller) Sum(reqBody io.ReadCloser) (string, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(c.cfg.Key))
	h.Write(body)
	dst := h.Sum(nil)
	hash := base64.StdEncoding.EncodeToString(dst)
	return hash, nil
}
