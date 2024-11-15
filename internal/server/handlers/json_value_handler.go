package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) JSONValueHandler(w http.ResponseWriter, r *http.Request) {
	var jMetric models.Metrics

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&jMetric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем тип метрики
	if err := types.CheckMetricType(jMetric.MType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем имя метрики
	if len(jMetric.ID) == 0 {
		http.Error(w, "the name of the metric is empty", http.StatusNotFound)
		return
	}

	metric, ok := h.Repo.Get(jMetric.ID)
	if !ok {
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	jMetric = metric.Convert(jMetric.ID)

	resp, err := json.Marshal(jMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
