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

		// reqHash := r.Header.Get("HashSHA256")
		// noneHash := r.Header.Get("Hash")
		// if noneHash == "none" {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		reqHash := r.Header.Get("HashSHA256")
		if reqHash == "" {
			next.ServeHTTP(w, r)
			// h.lg.Sugar.Infoln("error in request handler. hash is empty")
			// http.Error(w, "hash is empty", http.StatusBadRequest)
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
