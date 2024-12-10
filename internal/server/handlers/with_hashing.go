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
		if len(h.config.Key) > 0 {
			reqHash := r.Header.Get("HashSHA256")

			// copyBody, err := r.Body()
			// if err != nil {
			// 	h.lg.Sugar.Infow("error in request handler", "error: ", err)
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

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
				h.lg.Sugar.Infow("error in request handler", "error: ", err)
				http.Error(w, "hashes are not equal", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
