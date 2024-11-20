package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandlers(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "Status ok test",
			url:  "/update/gauge/metric/100",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "No name metrics test",
			url:  "/update/gauge//100",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong metrics type",
			url:  "/update/gaaauge/metric/100",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong value",
			url:  "/update/counter/metric/100.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "Wrong value",
			url:  "/update/counter/metric/aa",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	h := NewHandlers(storage.NewStorage(), config.Config{}, logger.Logger{})
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	serv := httptest.NewServer(mux)
	defer serv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodPost, serv.URL+test.url, nil)
			assert.NoError(t, err)

			res, err := serv.Client().Do(request)
			assert.NoError(t, err)
			defer res.Body.Close()
		})
	}
}
