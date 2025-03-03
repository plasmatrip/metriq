// The Update function handles HTTP requests for updating a metric. It extracts the metric type,
// name, and value from the request path. The function performs several validations:
// 1. Metric Type Validation: It checks if the provided metric type is valid using the CheckMetricType
//    function. If it's not valid, the function responds with a 400 Bad Request status.
// 2. Metric Name Validation: It ensures that the metric name is not empty. If the name is empty, it
//    responds with a 404 Not Found status.
// 3. Metric Value Validation: It validates the metric value based on its type using the CheckValue
//    function. If the value is invalid, it responds with a 400 Bad Request status.
// Once all validations pass, the function proceeds to handle the metric update logic.

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

	if err = h.Repo.SetMetric(r.Context(), mName, types.Metric{MetricType: mType, Value: value}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprint("Successful data update: ", mType, " ", mName, " ", metricValue, "\r\n")))
	if err != nil {
		fmt.Println(err.Error())
	}
}
