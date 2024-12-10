package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) JSONValue(w http.ResponseWriter, r *http.Request) {
	var jMetric models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&jMetric); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем тип метрики
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

	metric, err := h.Repo.Metric(r.Context(), jMetric.ID)
	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	jMetric = metric.Convert(jMetric.ID)

	resp, err := json.Marshal(jMetric)
	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// если есть ключ, хэшируем ответ
	if len(h.config.Key) > 0 {
		hash, err := h.Sum(resp)
		if err != nil {
			h.lg.Sugar.Infow("error in request handler", "error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HashSHA256", hash)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
