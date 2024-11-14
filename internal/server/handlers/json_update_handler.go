package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/metriq/internal/model"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) JSONUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var jMetric model.Metrics

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

	var value any
	switch jMetric.MType {
	case types.Counter:
		value = *jMetric.Delta
	case types.Gauge:
		value = *jMetric.Value
	}

	if err := h.Repo.Update(jMetric.ID, types.Metric{MetricType: jMetric.MType, Value: value}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
