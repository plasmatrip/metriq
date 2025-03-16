// This middleware function is used to authenticate incoming requests to the server.
// It computes a sha256 hash of the request body and compares it with the value
// of the "HashSHA256" header from the request. If the values match, the request
// is allowed to proceed to the next handler in the chain. If the values do not
// match, the function returns an error response with a status code of 400.
// The function is used to prevent tampering with the request body and ensure
// that the data is not modified during transport.
package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

func (h Handlers) WithHashing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqHash := r.Header.Get("HashSHA256")
		if reqHash == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.lg.Sugar.Infow("error in request handler", "error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))

		hm := hmac.New(sha256.New, []byte(h.config.Key))
		hm.Write(body)
		dst := hm.Sum(nil)
		sumHash := base64.StdEncoding.EncodeToString(dst)

		if !hmac.Equal([]byte(reqHash), []byte(sumHash)) {
			h.lg.Sugar.Infow("error in request handler", "req: ", reqHash, "sum: ", sumHash)
			http.Error(w, "hashes are not equal", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
