package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	metrics := h.Repo.Metrics()

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
	_, err := w.Write([]byte(html))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
