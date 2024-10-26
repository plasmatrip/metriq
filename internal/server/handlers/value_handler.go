package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/plasmatrip/metriq/internal/server"
)

func (h *Handlers) ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	uri := strings.Split(r.URL.RequestURI(), "/")

	if len(uri) != server.ValueURILen {
		http.Error(w, "Request not recognized!", http.StatusNotFound)
		return
	}

	//проверяем имя метрики
	metricName := uri[server.RequestNamePos]
	if err := server.CheckMetricName(metricName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//проверяем тип метрики
	metricType := uri[server.RequestTypePos]
	if err := server.CheckMetricType(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var formatedValue string
	switch metricType {
	case server.Gauge:
		value := h.Repo.GetGauge(metricName)
		formatedValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case server.Counter:
		value := h.Repo.GetCounter(metricName)
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
