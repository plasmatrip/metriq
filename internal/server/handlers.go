package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/plasmatrip/metriq/internal/storage"
)

type Handlers struct {
	Repo storage.Repository
}

func NewHandlers(repo storage.Repository) *Handlers {
	return &Handlers{Repo: repo}
}

func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if _, ok := r.Header[`Content-Type`]; !ok {
		http.Error(w, "Only the 'text/plain' content type is allowed!", http.StatusUnsupportedMediaType)
		return
	}

	uri := strings.Split(r.URL.RequestURI(), "/")

	if len(uri) != updateURILen {
		http.Error(w, "Request not recognized!", http.StatusNotFound)
		return
	}

	//проверяем имя метрики
	metricName := uri[mNamePos]
	if len(metricName) == 0 || !CheckName(metricName) {
		http.Error(w, "The name of the metric is not defined!", http.StatusNotFound)
		return
	}

	//проверяем тип метрики
	metricType := uri[mTypePos]
	if len(metricType) == 0 || !CheckType(metricType) {
		http.Error(w, "The type of the metric is not defined!", http.StatusBadRequest)
		return
	}

	//проверяем значение метрики
	metricValue := uri[mValuePos]
	if err := CheckValue(metricType, metricValue); err != nil {
		http.Error(w, "Unknown value!", http.StatusBadRequest)
		return
	}

	h.Repo.Update(metricName, metricValue)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprint("Successful data update! ", metricType, " ", metricName, " ", metricValue, "\r\n")))
	if err != nil {
		fmt.Println(err.Error())
	}
}
