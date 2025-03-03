// The JSONUpdate function is a request handler that processes incoming requests to update a metric.
// It reads the request body, decodes the JSON into a models.Metrics struct and checks that the metric type is valid.
// If the metric type is invalid, it logs an error and returns a 400 status code.
// It then checks that the metric name is not empty. If the name is empty, it logs an error and returns a 404 status code.
// If all checks pass, it calls the SetMetric method of the repository to update the metric.
// The function does not return any data in the response body.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) JSONUpdate(w http.ResponseWriter, r *http.Request) {
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

	//проверяем имя метрики и не пустое ли оно
	if len(jMetric.ID) == 0 {
		h.lg.Sugar.Infow("error in request handler", "error: ", "the name of the metric is empty")
		http.Error(w, "the name of the metric is empty", http.StatusNotFound)
		return
	}
	_, err := h.Repo.Metric(r.Context(), jMetric.ID)
	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", "metric not found")
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	var value any
	switch jMetric.MType {
	case types.Counter:
		value = *jMetric.Delta
	case types.Gauge:
		value = *jMetric.Value
	}

	if err = h.Repo.SetMetric(r.Context(), jMetric.ID, types.Metric{MetricType: jMetric.MType, Value: value}); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Repo.Metric(r.Context(), jMetric.ID)
	if err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", "metric not found")
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
