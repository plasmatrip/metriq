package handlers

import (
	"fmt"
	"net/http"

	"github.com/plasmatrip/metriq/internal/types"
)

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	//проверяем тип метрики
	mType := r.PathValue("metricType")
	if err := types.CheckMetricType(mType); err != nil {
		http.Error(w, mType, http.StatusBadRequest)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем имя метрики
	mName := r.PathValue("metricName")
	if len(mName) == 0 {
		http.Error(w, "the name of the metric is empty", http.StatusNotFound)
		return
	}

	//проверяем значение метрики
	metricValue := r.PathValue("metricValue")
	value, err := types.CheckValue(mType, metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Repo.SetMetric(r.Context(), mName, types.Metric{MetricType: mType, Value: value}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprint("Successful data update: ", mType, " ", mName, " ", metricValue, "\r\n")))
	if err != nil {
		fmt.Println(err.Error())
	}
}
