package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
)

func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get(`Content-Type`)
	if !strings.Contains(`text/plain`, contentType) {
		http.Error(w, "Only the 'text/plain' content type is allowed!", http.StatusUnsupportedMediaType)
		return
	}

	uri := strings.Split(r.URL.RequestURI(), "/")

	if len(uri) != server.UpdateURILen {
		http.Error(w, "Request not recognized!", http.StatusNotFound)
		return
	}

	//проверяем тип метрики
	metricType := uri[server.RequestTypePos]
	if err := server.CheckMetricType(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем имя метрики
	metricName := uri[server.RequestNamePos]
	if err := server.MetricNameNotEmpty(metricName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	//проверяем значение метрики
	metricValue := uri[server.RequestValuePos]
	if err := server.CheckValue(metricType, metricValue); err != nil {
		http.Error(w, "Unknown value!", http.StatusBadRequest)
		return
	}

	switch metricType {
	case server.Gauge:
		value, _ := strconv.ParseFloat(metricValue, 64)
		h.Repo.UpdateGauge(metricName, storage.Gauge(value))
	case server.Counter:
		value, _ := strconv.ParseInt(metricValue, 10, 64)
		h.Repo.UpdateCounter(value)
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprint("Successful data update! ", metricType, " ", metricName, " ", metricValue, "\r\n")))
	if err != nil {
		fmt.Println(err.Error())
	}
}
