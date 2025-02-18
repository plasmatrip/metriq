// The function Sum is used to generate a hash of the response body.
// The generated hash is then written to the "HashSHA256" header of the response.
// The hash is generated using the HMAC-SHA256 algorithm with the secret key
// provided in the configuration. The hash is then encoded to a string using
// base64 encoding.
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
