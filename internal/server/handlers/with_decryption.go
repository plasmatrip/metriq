package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
)

func (h Handlers) WithDecryption(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Чтение зашифрованных данных из тела запроса
		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			h.lg.Sugar.Infow("error in request handler", "error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Расшифровка данных
		decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, h.config.CryptoKey, encryptedData)
		if err != nil {
			h.lg.Sugar.Infow("error encryption data", "error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Заменяем тело запроса на расшифрованные данные
		r.Body = io.NopCloser(bytes.NewReader(decryptedData))
		r.ContentLength = int64(len(decryptedData))

		// Передаем запрос следующему обработчику
		next.ServeHTTP(w, r)
	})
}
