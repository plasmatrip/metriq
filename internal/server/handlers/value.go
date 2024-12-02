package handlers

import (
	"net/http"
	"strconv"

	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) Value(w http.ResponseWriter, r *http.Request) {
	//получаем имя метрики
	mName := r.PathValue("metricName")
	if len(mName) == 0 {
		h.lg.Sugar.Infoln("Metric name is undefined")
		http.Error(w, "Metric name is undefined", http.StatusBadRequest)
		return
	}

	//проверяем тип метрики
	mType := r.PathValue("metricType")
	if err := types.CheckMetricType(mType); err != nil {
		h.lg.Sugar.Infoln("Metric name is undefined")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := h.Repo.Metric(r.Context(), mName)
	if err != nil {
		h.lg.Sugar.Infoln("Metric name is undefined")
		http.Error(w, "Metric not found", http.StatusNotFound)
		return
	}

	var formatedValue string
	switch metric.MetricType {
	case types.Gauge:
		value, ok := metric.Value.(float64)
		if !ok {
			h.lg.Sugar.Infoln("Failed to cast the received value to type float64")
			http.Error(w, "Failed to cast the received value to type float64", http.StatusInternalServerError)
			return
		}
		formatedValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case types.Counter:
		value, ok := metric.Value.(int64)
		if !ok {
			h.lg.Sugar.Infoln("Failed to cast the received value to type int64")
			http.Error(w, "Failed to cast the received value to type int64", http.StatusInternalServerError)
			return
		}
		formatedValue = strconv.FormatInt(int64(value), 10)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(formatedValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
