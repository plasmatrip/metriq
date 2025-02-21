// Metrics - GET /metrics - returns a list of all metrics and their current values as HTML page.
// It does not take any parameters. It returns a list of all metrics and their current values as HTML page.
// The list of metrics is retrieved from the repository and then written to the HTTP response.
// If any error occur during the request, it returns an error with the appropriate HTTP status code.
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
