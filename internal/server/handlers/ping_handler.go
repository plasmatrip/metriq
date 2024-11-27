package handlers

import (
	"net/http"
)

func (h *Handlers) PingHandler(w http.ResponseWriter, r *http.Request) {
	err := h.Repo.Ping()
	if err != nil {
		http.Error(w, "Failed to ping datbase", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
}
