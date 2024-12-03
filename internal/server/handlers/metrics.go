package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.Repo.Metrics(r.Context())

	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<title>All metric params</title>
		</head>
		<body>
			%v
		</body>
		</html>
		`, metrics)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(html))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
