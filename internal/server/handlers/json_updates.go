package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) JSONUpdates(w http.ResponseWriter, r *http.Request) {
	jMetrics := []models.Metrics{}

	if err := json.NewDecoder(r.Body).Decode(&jMetrics); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, jMetric := range jMetrics {
		// проверяем тип метрики
		if err := types.CheckMetricType(jMetric.MType); err != nil {
			h.lg.Sugar.Infow("error in request handler", "error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//проверяем имя метрики
		if len(jMetric.ID) == 0 {
			h.lg.Sugar.Infow("error in request handler", "error: ", "the name of the metric is empty")
			http.Error(w, "the name of the metric is empty", http.StatusNotFound)
			return
		}

	}

	if err := h.Repo.SetMetrics(r.Context(), jMetrics); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(fmt.Sprintf("%d metrics received", len(jMetrics)))
	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
