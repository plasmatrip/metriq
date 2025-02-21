// This function is a handler for "/ping" endpoint. It sends a request to database to check if it is alive.
// If the request fails, it returns 500 status code.
// If the request is successfull, it returns 200 status code.
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
