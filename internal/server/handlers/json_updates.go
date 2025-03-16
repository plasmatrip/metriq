// The JSONUpdates function is a request handler for the HTTP POST method.
// It receives a JSON payload containing an array of metrics, each with a name,
// type, and value. The handler checks the type of each metric and logs an error
// if it is not one of the supported types. It also checks the name of the metric
// and logs an error if it is empty. The handler then iterates over the list of
// metrics and calls the SetMetric method of the repository for each one, storing
// the metric in the repository. The handler logs an error if the repository
// returns an error. The handler returns a JSON response with the list of metrics
// that were successfully written to the repository.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

var mPool = sync.Pool{
	New: func() interface{} {
		return &[]models.Metrics{}
	},
}

func (h *Handlers) JSONUpdates(w http.ResponseWriter, r *http.Request) {
	jMetrics := mPool.Get().(*[]models.Metrics)

	// read body
	// var body []byte
	// _, err := r.Body.Read(body)
	// if err != nil {
	// 	h.lg.Sugar.Infow("error in request handler", "error: ", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// decrypted body
	// decryptedBody, err := cert.DecryptData(body, h.config.CryptoKey)
	// if err != nil {
	// 	h.lg.Sugar.Infow("error in request handler", "error: ", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// unmarshal
	// if err = json.Unmarshal(decryptedBody, &jMetrics); err != nil {
	// 	h.lg.Sugar.Infow("error in request handler", "error: ", err)
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	if err := json.NewDecoder(r.Body).Decode(&jMetrics); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, jMetric := range *jMetrics {
		// проверяем тип метрики
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

	}

	if err := h.Repo.SetMetrics(r.Context(), *jMetrics); err != nil {
		h.lg.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(fmt.Sprintf("%d metrics received", len(*jMetrics)))
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
