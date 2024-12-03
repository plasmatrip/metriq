package handlers

import (
	"net/http"
)

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := h.Repo.Ping(r.Context())
	if err != nil {
		h.lg.Sugar.Infow(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
