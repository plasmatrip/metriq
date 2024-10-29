package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	uri := strings.Split(r.URL.RequestURI(), "/")

	if len(uri) != server.ValueURLLen {
		http.Error(w, "Request not recognized!", http.StatusNotFound)
		return
	}

	//проверяем имя метрики
	metricName := uri[server.RequestNamePos]

	//проверяем тип метрики
	metricType := uri[server.RequestTypePos]
	if err := server.CheckMetricType(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, ok := h.Repo.Get(metricName)
	if !ok {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}

	var formatedValue string
	switch metric.MetricType {
	case types.Gauge:
		value, ok := metric.Value.(float64)
		if !ok {
			http.Error(w, "failed to cast the received value to type float64", http.StatusInternalServerError)
			return
		}
		formatedValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case types.Counter:
		value, ok := metric.Value.(int64)
		if !ok {
			http.Error(w, "failed to cast the received value to type int64", http.StatusInternalServerError)
			return
		}
		formatedValue = strconv.FormatInt(int64(value), 10)
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(formatedValue))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
